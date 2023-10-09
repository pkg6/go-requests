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
	return c.WithContentType(w.FormDataContentType()).Post(ctx, uri, body)
}
