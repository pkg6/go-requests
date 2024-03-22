package requests

import (
	"context"
	"net/http"
	"net/url"
)

func (c *Client) GetD(ctx context.Context, uri string, data, d any) (response *Response, err error) {
	return c.DoRequestD(ctx, http.MethodGet, uri, data, d)
}
func (c *Client) HeadD(ctx context.Context, uri string, data, d any) (response *Response, err error) {
	return c.DoRequestD(ctx, http.MethodHead, uri, data, d)
}
func (c *Client) PostD(ctx context.Context, uri string, data, d any) (response *Response, err error) {
	return c.DoRequestD(ctx, http.MethodPost, uri, data, d)
}
func (c *Client) PutD(ctx context.Context, uri string, data, d any) (response *Response, err error) {
	return c.DoRequestD(ctx, http.MethodPut, uri, data, d)
}
func (c *Client) PatchD(ctx context.Context, uri string, data, d any) (response *Response, err error) {
	return c.DoRequestD(ctx, http.MethodPatch, uri, data, d)
}
func (c *Client) DeleteD(ctx context.Context, uri string, data, d any) (response *Response, err error) {
	return c.DoRequestD(ctx, http.MethodDelete, uri, data, d)
}
func (c *Client) ConnectD(ctx context.Context, uri string, data, d any) (response *Response, err error) {
	return c.DoRequestD(ctx, http.MethodConnect, uri, data, d)
}
func (c *Client) OptionsD(ctx context.Context, uri string, data, d any) (response *Response, err error) {
	return c.DoRequestD(ctx, http.MethodOptions, uri, data, d)
}
func (c *Client) TraceD(ctx context.Context, uri string, data, d any) (response *Response, err error) {
	return c.DoRequestD(ctx, http.MethodTrace, uri, data, d)
}
func (c *Client) PostJsonD(ctx context.Context, uri string, data, d any) (response *Response, err error) {
	return c.AsJson().DoRequestD(ctx, http.MethodPost, uri, data, d)
}
func (c *Client) PostFormD(ctx context.Context, uri string, data url.Values, d any) (response *Response, err error) {
	w, body, err := buildFormBody(data)
	if err != nil {
		return nil, err
	}
	return c.WithContentType(w.FormDataContentType()).PostD(ctx, uri, body, d)
}
