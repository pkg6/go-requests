package requests

import (
	"context"
	"net/url"
)

// Deprecated: use PostFormD instead. To be removed in v0.3.0 release.
func (c *Client) PostFormWithFilesD(ctx context.Context, uri string, data url.Values, d any) (response *Response, err error) {
	return c.PostFormD(ctx, uri, data, d)
}

// Deprecated: use PostForm instead. To be removed in v0.3.0 release.
func (c *Client) PostFormWithFiles(ctx context.Context, uri string, data url.Values) (*Response, error) {
	return c.PostForm(ctx, uri, data)
}

// Deprecated: use PostFormBytes instead. To be removed in v0.3.0 release.
func (c *Client) PostFormWithFilesBytes(ctx context.Context, uri string, data url.Values) ([]byte, error) {
	return c.PostFormBytes(ctx, uri, data)
}

// Deprecated: use PostFormUnmarshal instead. To be removed in v0.3.0 release.
func (c *Client) PostFormWithFilesUnmarshal(ctx context.Context, uri string, data url.Values, d any) error {
	return c.PostFormUnmarshal(ctx, uri, data, d)
}

// Deprecated: use PostForm instead. To be removed in v0.3.0 release.
func PostFormWithFiles(uri string, data url.Values, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PostFormWithFiles(context.Background(), uri, data)
}

// Deprecated: use PostFormUnmarshal instead. To be removed in v0.3.0 release.
func PostFormWithFilesUnmarshal(uri string, data url.Values, d any, args ...ArgsFunc) error {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PostFormWithFilesUnmarshal(context.Background(), uri, data, d)
}
