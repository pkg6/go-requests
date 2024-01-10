package requests

import "net/http"

const (
	clientMiddlewareKey CtxKeyString = "__client_middleware_key"
)

// MiddlewareFunc middleware handler func
type MiddlewareFunc = func(c *Client, r *http.Request) (*Response, error)

// clientMiddleware is the plugin for http client request workflow management.
type clientMiddleware struct {
	client       *Client          // http client.
	handlers     []MiddlewareFunc // mdl handlers.
	handlerIndex int              // current handler index.
	response     *Response        // save resp.
	err          error            // save err.
}

// Use adds one or more middleware handlers to client.
func (c *Client) Use(middlewares ...MiddlewareFunc) *Client {
	c.middlewares = append(c.middlewares, middlewares...)
	return c
}

// Next calls the next middleware.
// This should only be call in HandlerFunc.
func (c *Client) Next(req *http.Request) (*Response, error) {
	if v := req.Context().Value(clientMiddlewareKey); v != nil {
		if m, ok := v.(*clientMiddleware); ok {
			return m.Next(req)
		}
	}
	return c.callRequest(req)
}

// Next calls the next middleware handler.
func (m *clientMiddleware) Next(req *http.Request) (response *Response, err error) {
	if m.err != nil {
		return m.response, m.err
	}
	if m.handlerIndex < len(m.handlers) {
		m.handlerIndex++
		m.response, m.err = m.handlers[m.handlerIndex](m.client, req)
	}
	return m.response, m.err
}
