package requests

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func (c *Client) GetBytes(ctx context.Context, uri string, data any) ([]byte, error) {
	return c.DoRequestBytes(ctx, http.MethodGet, uri, data)
}
func (c *Client) HeadBytes(ctx context.Context, uri string, data any) ([]byte, error) {
	return c.DoRequestBytes(ctx, http.MethodHead, uri, data)
}
func (c *Client) PostBytes(ctx context.Context, uri string, data any) ([]byte, error) {
	return c.DoRequestBytes(ctx, http.MethodPost, uri, data)
}
func (c *Client) PutBytes(ctx context.Context, uri string, data any) ([]byte, error) {
	return c.DoRequestBytes(ctx, http.MethodPut, uri, data)
}
func (c *Client) PatchBytes(ctx context.Context, uri string, data any) ([]byte, error) {
	return c.DoRequestBytes(ctx, http.MethodPatch, uri, data)
}
func (c *Client) DeleteBytes(ctx context.Context, uri string, data any) ([]byte, error) {
	return c.DoRequestBytes(ctx, http.MethodDelete, uri, data)
}
func (c *Client) ConnectBytes(ctx context.Context, uri string, data any) ([]byte, error) {
	return c.DoRequestBytes(ctx, http.MethodConnect, uri, data)
}
func (c *Client) OptionsBytes(ctx context.Context, uri string, data any) ([]byte, error) {
	return c.DoRequestBytes(ctx, http.MethodOptions, uri, data)
}
func (c *Client) TraceBytes(ctx context.Context, uri string, data any) ([]byte, error) {
	return c.DoRequestBytes(ctx, http.MethodTrace, uri, data)
}
func (c *Client) PostJsonBytes(ctx context.Context, uri string, data any) ([]byte, error) {
	return c.AsJson().DoRequestBytes(ctx, http.MethodPost, uri, data)
}
func (c *Client) PostFormBytes(ctx context.Context, uri string, data url.Values) ([]byte, error) {
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
	return c.WithContentType(w.FormDataContentType()).PostBytes(ctx, uri, body)
}
func (c *Client) PostFormWithFilesBytes(ctx context.Context, uri string, data url.Values) ([]byte, error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k := range data {
		v := data.Get(k)
		if strings.Contains(v, HttpParamFileHolder) {
			localPathFile := strings.ReplaceAll(strings.ReplaceAll(v, HttpParamFileHolder, ""), " ", "")
			osFile, err := os.Open(localPathFile)
			if err != nil {
				return nil, err
			}
			ioWriter, err := w.CreateFormFile(k, k)
			if err != nil {
				return nil, err
			}
			if _, err = io.Copy(ioWriter, osFile); err != nil {
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
	return c.WithContentType(w.FormDataContentType()).PostBytes(ctx, uri, body)
}
