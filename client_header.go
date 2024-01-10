package requests

import (
	"encoding/base64"
)

// WithHeader method sets a single header field and its value in the client instance.
// These headers will be applied to all requests raised from this client instance.
// Also it can be overridden at request level header options.
//		WithHeader("Content-Type", "application/json").
//		WithHeader("Accept", "application/json")
func (c *Client) WithHeader(k, v string) *Client {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Header.Set(k, v)
	return c
}

// WithHeaderMap method sets multiple headers field and its values at one go in the client instance.
// These headers will be applied to all requests raised from this client instance. Also it can be
// overridden at request level headers options.
// For Example: To set `Content-Type` and `Accept` as `application/json`
//		WithHeaders(map[string]string{
//			"Content-Type": "application/json",
//			"Accept": "application/json",
//		})
func (c *Client) WithHeaderMap(headers map[string]string) *Client {
	for h, v := range headers {
		c.WithHeader(h, v)
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
func (c *Client) WithHeaderVerbatim(k, v string) *Client {
	c.Header[k] = []string{v}
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

func (c *Client) WithRandomUserAgent() *Client {
	c.WithUserAgent(RandomUserAgent())
	return c
}
func (c *Client) WithRandomMobileUserAgent() *Client {
	c.WithUserAgent(RandomMobileUserAgent())
	return c
}

// AsForm is a chaining function,
// which sets the HTTP content type as "application/x-www-form-urlencoded" for the next request.
func (c *Client) AsForm() *Client {
	c.WithContentType(HttpHeaderContentTypeForm)
	return c
}

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
	c.WithToken(base64.StdEncoding.EncodeToString([]byte(username+":"+password)), AuthorizationTypeBasic)
	return c
}

// WithToken
//Specify an authorization token for the request.
func (c *Client) WithToken(token string, authorizationType ...string) *Client {
	if len(authorizationType) > 0 {
		token = authorizationType[0] + token
	} else {
		token = AuthorizationTypeBearer + token
	}
	c.Header.Set(HttpHeaderAuthorization, token)
	return c
}
