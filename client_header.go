package requests

import "encoding/base64"

// WithHeader method sets a single header field and its value in the client instance.
// These headers will be applied to all requests raised from this client instance.
// Also it can be overridden at request level header options.
//		WithHeader("Content-Type", "application/json").
//		WithHeader("Accept", "application/json")
func (c *Client) WithHeader(header, value string) *Client {
	c.header.Set(header, value)
	return c
}

// WithHeaders method sets multiple headers field and its values at one go in the client instance.
// These headers will be applied to all requests raised from this client instance. Also it can be
// overridden at request level headers options.
// For Example: To set `Content-Type` and `Accept` as `application/json`
//		WithHeaders(map[string]string{
//			"Content-Type": "application/json",
//			"Accept": "application/json",
//		})
func (c *Client) WithHeaders(headers map[string]string) *Client {
	for h, v := range headers {
		c.header.Set(h, v)
	}
	return c
}

// WithHeaderVerbatim method is to set a single header field and its value verbatim in the current request.
//
// For Example: To set `all_lowercase` and `UPPERCASE` as `available`.
//		WithHeaderVerbatim("all_lowercase", "available").
//		WithHeaderVerbatim("UPPERCASE", "available")
//
// Also you can override header value, which was set at client instance level.
func (c *Client) WithHeaderVerbatim(header, value string) *Client {
	c.header[header] = []string{value}
	return c
}

// WithContentType is a chaining function,
// which sets HTTP content type for the next request.
func (c *Client) WithContentType(contentType string) *Client {
	c.WithHeader(HttpHeaderContentType, contentType)
	return c
}
func (c *Client) WithUserAgent(userAgent string) *Client {
	c.WithHeader(HttpHeaderUserAgent, userAgent)
	return c
}

// AsForm is a chaining function,
// which sets the HTTP content type as "application/x-www-form-urlencoded" for the next request.
func (c *Client) AsForm() *Client {
	c.WithContentType(HttpHeaderContentTypeForm)
	return c
}

// AsStream is a Stream
//func (c *Client) AsStream() *Client {
//	c.WithHeaders(map[string]string{
//		HttpHeaderAccept:       HttpMIMEEventStream,
//		HttpHeaderCacheControl: "no-cache",
//		HttpHeaderConnection:   "keep-alive",
//	})
//	return c
//}

// AsJson is a chaining function,
// which sets the HTTP content type as "application/json" for the next request.
//
// Note that it also checks and encodes the parameter to JSON format automatically.
func (c *Client) AsJson() *Client {
	c.WithContentType(HttpHeaderContentTypeJson)
	return c
}

// AsXml is a chaining function,
// which sets the HTTP content type as "application/xml" for the next request.
//
// Note that it also checks and encodes the parameter to XML format automatically.
func (c *Client) AsXml() *Client {
	c.WithContentType(HttpHeaderContentTypeXml)
	return c
}

// WithBasicAuth
//Specify the basic authentication username and password for the request.
func (c *Client) WithBasicAuth(username, password string) *Client {
	c.WithToken(base64.StdEncoding.EncodeToString([]byte(username+":"+password)), "Basic ")
	return c
}

// WithToken
//Specify an authorization token for the request.
func (c *Client) WithToken(token string, Type ...string) *Client {
	if len(Type) > 0 {
		token = Type[0] + token
	} else {
		token = "Bearer " + token
	}
	c.header.Set(HttpHeaderAuthorization, token)
	return c
}
