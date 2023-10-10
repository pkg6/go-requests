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
	return c.WithContentType(w.FormDataContentType()).PostUnmarshal(ctx, uri, body, d)
}
func (c *Client) PostFormWithFilesUnmarshal(ctx context.Context, uri string, data url.Values, d any) error {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k := range data {
		v := data.Get(k)
		if strings.Contains(v, HttpParamFileHolder) {
			localPathFile := strings.ReplaceAll(strings.ReplaceAll(v, HttpParamFileHolder, ""), " ", "")
			osFile, err := os.Open(localPathFile)
			if err != nil {
				return err
			}
			ioWriter, err := w.CreateFormFile(k, k)
			if err != nil {
				return err
			}
			if _, err = io.Copy(ioWriter, osFile); err != nil {
				return err
			}
		} else {
			if err := w.WriteField(k, v); err != nil {
				return err
			}
		}
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.WithContentType(w.FormDataContentType()).PostUnmarshal(ctx, uri, body, d)
}
