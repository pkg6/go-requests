package requests

import (
	"context"
	"net/http"
	"net/url"
)

func (c *Client) GetUnmarshal(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestUnmarshal(ctx, http.MethodGet, uri, data, d)
}
func (c *Client) HeadUnmarshal(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestUnmarshal(ctx, http.MethodHead, uri, data, d)
}
func (c *Client) PostUnmarshal(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestUnmarshal(ctx, http.MethodPost, uri, data, d)
}
func (c *Client) PutUnmarshal(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestUnmarshal(ctx, http.MethodPut, uri, data, d)
}
func (c *Client) PatchUnmarshal(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestUnmarshal(ctx, http.MethodPatch, uri, data, d)
}
func (c *Client) DeleteUnmarshal(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestUnmarshal(ctx, http.MethodDelete, uri, data, d)
}
func (c *Client) ConnectUnmarshal(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestUnmarshal(ctx, http.MethodConnect, uri, data, d)
}
func (c *Client) OptionsUnmarshal(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestUnmarshal(ctx, http.MethodOptions, uri, data, d)
}
func (c *Client) TraceUnmarshal(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestUnmarshal(ctx, http.MethodTrace, uri, data, d)
}
func (c *Client) PostJsonUnmarshal(ctx context.Context, uri string, data, d any) error {
	return c.AsJson().DoRequestUnmarshal(ctx, http.MethodPost, uri, data, d)
}
func (c *Client) PostFormUnmarshal(ctx context.Context, uri string, data url.Values, d any) error {
	w, body, err := buildFormBody(data)
	if err != nil {
		return err
	}
	return c.WithContentType(w.FormDataContentType()).PostUnmarshal(ctx, uri, body, d)
}
