package requests

import (
	"context"
	"crypto/tls"
	"net/url"
)

//requests helper function

type ArgsFunc func(client *Client)

func Debug() ArgsFunc {
	return func(client *Client) {
		client.Debug()
	}
}
func WithCookie(k, v string) ArgsFunc {
	return func(client *Client) {
		client.WithCookie(k, v)
	}
}
func WithCookies(cookies map[string]string) ArgsFunc {
	return func(client *Client) {
		client.WithCookies(cookies)
	}
}
func WithCookieString(cookies string) ArgsFunc {
	return func(client *Client) {
		client.WithCookieString(cookies)
	}
}
func WithHeader(header, value string) ArgsFunc {
	return func(client *Client) {
		client.WithHeader(header, value)
	}
}
func WithHeaders(headers map[string]string) ArgsFunc {
	return func(client *Client) {
		client.WithHeaders(headers)
	}
}
func WithContentType(contentType string) ArgsFunc {
	return func(client *Client) {
		client.WithContentType(contentType)
	}
}
func WithUserAgent(userAgent string) ArgsFunc {
	return func(client *Client) {
		client.WithUserAgent(userAgent)
	}
}
func WithBasicAuth(username, password string) ArgsFunc {
	return func(client *Client) {
		client.WithBasicAuth(username, password)
	}
}
func WithToken(token string, Type ...string) ArgsFunc {
	return func(client *Client) {
		client.WithToken(token, Type...)
	}
}
func WithProxyUrl(proxyURL string) ArgsFunc {
	return func(client *Client) {
		client.WithProxyUrl(proxyURL)
	}
}
func WithTLSKeyCrt(crtFile, keyFile string) ArgsFunc {
	return func(client *Client) {
		client.WithTLSKeyCrt(crtFile, keyFile)
	}
}
func SetTLSConfig(tlsConfig *tls.Config) ArgsFunc {
	return func(client *Client) {
		client.SetTLSConfig(tlsConfig)
	}
}

func Get(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Get(context.Background(), uri, data)
}

func Put(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Put(context.Background(), uri, data)
}

func Delete(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Delete(context.Background(), uri, data)
}

func Head(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Head(context.Background(), uri, data)
}

func Patch(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Patch(context.Background(), uri, data)
}

func Connect(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Connect(context.Background(), uri, data)
}

func Options(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Options(context.Background(), uri, data)
}

func Trace(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Trace(context.Background(), uri, data)
}

func Post(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Post(context.Background(), uri, data)
}

func PostJson(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.PostJson(context.Background(), uri, data)
}

func PostForm(uri string, data url.Values, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.PostForm(context.Background(), uri, data)
}

func Request(method, uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.DoRequest(context.Background(), method, uri, data)
}
