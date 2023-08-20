package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (c *Client) DoRequestD(ctx context.Context, method string, uri string, data, d any) error {
	response, err := c.DoRequest(ctx, method, uri, data)
	if err != nil {
		return nil
	}
	defer response.Close()
	defer func() {
		if err = response.Close(); err != nil {
			c.Logger.Fatalf("%v", err)
		}
	}()
	return response.Unmarshal(d)
}

func (c *Client) DoRequestBytes(ctx context.Context, method string, uri string, data any) []byte {
	response, err := c.DoRequest(ctx, method, uri, data)
	if err != nil {
		return nil
	}
	defer func() {
		if err = response.Close(); err != nil {
			c.Logger.Fatalf("%v", err)
		}
	}()
	return response.ReadAll()
}

func (c *Client) DoRequest(ctx context.Context, method, uri string, body any) (response *Response, err error) {
	if err = c.doBeforeRequestCallbacks(); err != nil {
		return nil, err
	}
	request, err := c.prepareRequest(ctx, method, uri, body)
	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				c.doPanicHooks(request, err)
			} else {
				c.doPanicHooks(request, fmt.Errorf("panic %v", rec))
			}
			panic(rec)
		}
	}()
	if err != nil {
		c.doErrorHooks(request, nil, err)
		return nil, err
	}
	if err = c.doAfterRequestCallbacks(request); err != nil {
		c.doErrorHooks(request, nil, err)
		return nil, err
	}
	// Client middleware.
	if len(c.middlewares) > 0 {
		middlewares := make([]MiddlewareFunc, 0, len(c.middlewares)+1)
		middlewares = append(middlewares, c.middlewares...)
		middlewares = append(middlewares, func(cli *Client, r *http.Request) (*Response, error) {
			return cli.callRequest(r)
		})
		ctx = context.WithValue(request.Context(),
			clientMiddlewareKey,
			&clientMiddleware{
				client:       c,
				handlers:     middlewares,
				handlerIndex: -1,
			})
		request = request.WithContext(ctx)
		response, err = c.Next(request)
	} else {
		response, err = c.callRequest(request)
	}
	if err != nil {
		c.doErrorHooks(request, response, err)
		return nil, err
	}
	if err = c.doResponseCallbacks(request, response); err != nil {
		c.doErrorHooks(request, response, err)
		return nil, err
	}
	c.doErrorHooks(request, response, err)
	return response, err
}
func (c *Client) callRequest(request *http.Request) (response *Response, err error) {
	response = &Response{request: request, client: c}
	reqBodyContent, err := io.ReadAll(request.Body)
	if err != nil {
		return response, fmt.Errorf(`io.ReadAll: %v`, err)
	}
	response.requestBody = reqBodyContent
	request.Body = NewReadCloser(reqBodyContent, false)
	for {
		if response.Response, err = c.Do(request); err != nil {
			c.attempt++
			// The response might not be nil when err != nil.
			if response.Response != nil {
				_ = response.Response.Body.Close()
			}
			if c.retryCount > 0 {
				c.retryCount--
				time.Sleep(c.retryWaitTime)
			} else {
				break
			}
			return response, fmt.Errorf(fmt.Sprintf(`client.Do: %v`, err))
		} else {
			break
		}
	}
	return response, err
}
func (c *Client) prepareBodyDefault(method string, body any) string {
	switch val := body.(type) {
	case []byte:
		return string(val)
	case string:
		return val
	case url.Values:
		return val.Encode()
	case *url.Values:
		return val.Encode()
	case bytes.Buffer:
		return val.String()
	case *bytes.Buffer:
		return val.String()
	case strings.Builder:
		return val.String()
	case map[string]string:
		uv := url.Values{}
		for s, s2 := range val {
			uv.Set(s, s2)
		}
		return uv.Encode()
	case map[string]any:
		return HttpBuildQuery(val)
	}
	if method == http.MethodGet && body != nil {
		if jsonByte, err := c.jsonMarshal(body); err == nil {
			mapAny := make(map[string]any)
			if err = c.jsonUnmarshal(jsonByte, &mapAny); err == nil {
				return HttpBuildQuery(mapAny)
			}
		}
	}
	return ToString(body)
}
func (c *Client) prepareBody(method string, body any) (string, error) {
	var params string
	if body != nil {
		contentType := c.header.Get(HttpHeaderContentType)
		if IsJSONType(contentType) {
			switch body.(type) {
			case string, []byte:
				params = ToString(body)
			default:
				if b, err := c.jsonMarshal(body); err != nil {
					return "", err
				} else {
					params = string(b)
				}
			}
		} else if IsXMLType(contentType) {
			switch body.(type) {
			case string, []byte:
				params = ToString(body)
			default:
				if b, err := c.xmlMarshal(body); err != nil {
					return "", err
				} else {
					params = string(b)
				}
			}
		} else {
			params = c.prepareBodyDefault(method, body)
		}
	}
	return params, nil
}

func (c *Client) prepareRequest(ctx context.Context, method, uri string, body any) (request *http.Request, err error) {
	method = strings.ToUpper(method)
	if len(c.BaseUrl) > 0 {
		uri = c.BaseUrl + strings.Trim(uri, "")
	}
	if !strings.Contains(uri, httpSchemeName) {
		uri = httpSchemeName + "://" + uri
	}
	if c.Query != nil {
		uri = Uri(uri, c.Query).String()
	}
	params, err := c.prepareBody(method, body)
	if err != nil {
		return nil, err
	}
	if method == http.MethodGet {
		var bodyBuffer *bytes.Buffer
		if params != "" {
			contentType := c.header.Get(HttpHeaderContentType)
			if IsJSONType(contentType) || IsXMLType(contentType) {
				bodyBuffer = bytes.NewBuffer([]byte(params))
			} else {
				if strings.Contains(uri, "?") {
					uri = uri + "&" + params
				} else {
					uri = uri + "?" + params
				}
				bodyBuffer = bytes.NewBuffer(nil)
			}
		} else {
			bodyBuffer = bytes.NewBuffer(nil)
		}
		if request, err = http.NewRequest(method, uri, bodyBuffer); err != nil {
			return nil, fmt.Errorf(`http.NewRequest failed with method "%s" and URL "%s"`, method, uri)
		}
	} else {
		paramBytes := []byte(params)
		if request, err = http.NewRequest(method, uri, bytes.NewReader(paramBytes)); err != nil {
			return nil, fmt.Errorf(`http.NewRequest failed for method "%s" and URL "%s"`, method, uri)
		} else {
			if v := c.header.Get(HttpHeaderContentType); v != "" {
				// Custom Content-Type.
				request.Header.Set(HttpHeaderContentType, v)
			} else if len(paramBytes) > 0 {
				if (paramBytes[0] == '[' || paramBytes[0] == '{') && json.Valid(paramBytes) {
					// Auto-detecting and setting the post content format: JSON.
					request.Header.Set(HttpHeaderContentType, HttpHeaderContentTypeJson)
				} else if IsMatchString(httpRegexParamJson, params) {
					// If the parameters passed like "name=value", it then uses form type.
					request.Header.Set(HttpHeaderContentType, HttpHeaderContentTypeForm)
				}
			}
		}
	}
	//重新加载上下文
	if withContext := c.withContext(ctx); withContext != nil {
		request = request.WithContext(withContext)
	}
	//加载cookie
	c.WithCookieString("")
	// Custom header.
	if len(c.header) > 0 {
		request.Header = c.header
	}
	if reqHeaderHost := request.Header.Get(HttpHeaderHost); reqHeaderHost != "" {
		request.Host = reqHeaderHost
	}
	return request, nil
}