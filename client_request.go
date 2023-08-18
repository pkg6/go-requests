package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Get send GET request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Get(ctx context.Context, uri string, data any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodGet, uri, data)
}

// Put send PUT request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Put(ctx context.Context, uri string, data any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodPut, uri, data)
}

// Post sends request using HTTP method POST and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Post(ctx context.Context, uri string, data any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodPost, uri, data)
}

// Delete send DELETE request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Delete(ctx context.Context, uri string, data any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodDelete, uri, data)
}

// Head send HEAD request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Head(ctx context.Context, uri string, data any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodHead, uri, data)
}

// Patch send PATCH request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Patch(ctx context.Context, uri string, data any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodPatch, uri, data)
}

// Connect send CONNECT request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Connect(ctx context.Context, uri string, data any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodConnect, uri, data)
}

// Options send OPTIONS request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Options(ctx context.Context, uri string, data any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodOptions, uri, data)
}

// Trace send TRACE request and returns the response object.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) Trace(ctx context.Context, uri string, data any) (*Response, error) {
	return c.DoRequest(ctx, http.MethodTrace, uri, data)
}

func (c *Client) PostJson(ctx context.Context, uri string, data any) (*Response, error) {
	return c.AsJson().Post(ctx, uri, data)
}

// PostForm is different from net/http.PostForm.
// It's a wrapper of Post method, which sets the Content-Type as "multipart/form-data;".
// and It will automatically set boundary characters for the request body and Content-Type.
//
// It's Seem like the following case:
//
// Content-Type: multipart/form-data; boundary=----Boundarye4Ghaog6giyQ9ncN
//
// And form data is like:
// ------Boundarye4Ghaog6giyQ9ncN
// Content-Disposition: form-data; name="checkType"
//
// none
//
// It's used for sending form data.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) PostForm(ctx context.Context, uri string, data url.Values) (*Response, error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k := range data {
		v := data.Get(k)
		if err := w.WriteField(k, v); err != nil {
			return nil, err
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return c.WithContentType(w.FormDataContentType()).Post(ctx, uri, body)
}

// PostFormWithFiles is different from net/http.PostForm.
// It's a wrapper of Post method, which sets the Content-Type as "multipart/form-data;".
// and It will automatically set boundary characters for the request body and Content-Type.
//
// It's Seem like the following case:
//
// Content-Type: multipart/form-data; boundary=----Boundarye4Ghaog6giyQ9ncN
//
// And form data is like:
// ------Boundarye4Ghaog6giyQ9ncN
// Content-Disposition: form-data; name="checkType"
//
// none
//
// It's used for sending form data.
// Note that the response object MUST be closed if it'll never be used.
func (c *Client) PostFormWithFiles(ctx context.Context, uri string, data url.Values) (*Response, error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k := range data {
		v := data.Get(k)
		if strings.Contains(v, httpParamFileHolder) {
			localPathFile := strings.ReplaceAll(strings.ReplaceAll(v, httpParamFileHolder, ""), " ", "")
			osfile, err := os.Open(localPathFile)
			if err != nil {
				return nil, err
			}
			ioWriter, err := w.CreateFormFile(k, k)
			if err != nil {
				return nil, err
			}
			if _, err = io.Copy(ioWriter, osfile); err != nil {
				return nil, err
			}
		} else {
			if err := w.WriteField(k, v); err != nil {
				return nil, err
			}
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return c.WithContentType(w.FormDataContentType()).Post(ctx, uri, body)
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
