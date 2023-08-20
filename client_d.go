package requests

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
)

func (c *Client) GetD(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestD(ctx, http.MethodGet, uri, data, d)
}
func (c *Client) HeadD(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestD(ctx, http.MethodHead, uri, data, d)
}
func (c *Client) PostD(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestD(ctx, http.MethodPost, uri, data, d)
}
func (c *Client) PutD(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestD(ctx, http.MethodPut, uri, data, d)
}
func (c *Client) PatchD(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestD(ctx, http.MethodPatch, uri, data, d)
}
func (c *Client) DeleteD(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestD(ctx, http.MethodDelete, uri, data, d)
}
func (c *Client) ConnectD(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestD(ctx, http.MethodConnect, uri, data, d)
}
func (c *Client) OptionsD(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestD(ctx, http.MethodOptions, uri, data, d)
}
func (c *Client) TraceD(ctx context.Context, uri string, data, d any) error {
	return c.DoRequestD(ctx, http.MethodTrace, uri, data, d)
}
func (c *Client) PostJsonD(ctx context.Context, uri string, data, d any) error {
	return c.AsJson().DoRequestD(ctx, http.MethodPost, uri, data, d)
}
func (c *Client) PostFormD(ctx context.Context, uri string, data url.Values, d any) error {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k := range data {
		v := data.Get(k)
		if err := w.WriteField(k, v); err != nil {
			return err
		}
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.WithContentType(w.FormDataContentType()).PostD(ctx, uri, body, d)
}
func (c *Client) PostFormWithFilesD(ctx context.Context, uri string, data url.Values, d any) error {
	response, err := c.PostFormWithFiles(ctx, uri, data)
	if err != nil {
		return err
	}
	return response.Unmarshal(d)
}
