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

func GetUnmarshal(uri string, data, d any, args ...ArgsFunc) (err error) {
	resp, err := Get(uri, data, args...)
	if err != nil {
		return
	}
	return resp.Unmarshal(d)
}

func Put(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Put(context.Background(), uri, data)
}
func PutUnmarshal(uri string, data, d any, args ...ArgsFunc) (err error) {
	resp, err := Put(uri, data, args...)
	if err != nil {
		return
	}
	return resp.Unmarshal(d)
}
func Delete(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Delete(context.Background(), uri, data)
}
func DeleteUnmarshal(uri string, data, d any, args ...ArgsFunc) (err error) {
	resp, err := Delete(uri, data, args...)
	if err != nil {
		return
	}
	return resp.Unmarshal(d)
}
func Head(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Head(context.Background(), uri, data)
}
func HeadUnmarshal(uri string, data, d any, args ...ArgsFunc) (err error) {
	resp, err := Head(uri, data, args...)
	if err != nil {
		return
	}
	return resp.Unmarshal(d)
}
func Patch(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Patch(context.Background(), uri, data)
}
func PatchUnmarshal(uri string, data, d any, args ...ArgsFunc) (err error) {
	resp, err := Patch(uri, data, args...)
	if err != nil {
		return
	}
	return resp.Unmarshal(d)
}
func Connect(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Connect(context.Background(), uri, data)
}
func ConnectUnmarshal(uri string, data, d any, args ...ArgsFunc) (err error) {
	resp, err := Connect(uri, data, args...)
	if err != nil {
		return
	}
	return resp.Unmarshal(d)
}
func Options(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Options(context.Background(), uri, data)
}
func OptionsUnmarshal(uri string, data, d any, args ...ArgsFunc) (err error) {
	resp, err := Options(uri, data, args...)
	if err != nil {
		return
	}
	return resp.Unmarshal(d)
}
func Trace(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Trace(context.Background(), uri, data)
}
func TraceUnmarshal(uri string, data, d any, args ...ArgsFunc) (err error) {
	resp, err := Trace(uri, data, args...)
	if err != nil {
		return
	}
	return resp.Unmarshal(d)
}
func Post(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Post(context.Background(), uri, data)
}
func PostUnmarshal(uri string, data, d any, args ...ArgsFunc) (err error) {
	resp, err := Post(uri, data, args...)
	if err != nil {
		return
	}
	return resp.Unmarshal(d)
}
func PostJson(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.PostJson(context.Background(), uri, data)
}
func PostJsonUnmarshal(uri string, data, d any, args ...ArgsFunc) (err error) {
	resp, err := PostJson(uri, data, args...)
	if err != nil {
		return
	}
	return resp.Unmarshal(d)
}
func PostForm(uri string, data url.Values, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.PostForm(context.Background(), uri, data)
}
func PostFormUnmarshal(uri string, data url.Values, d any, args ...ArgsFunc) (err error) {
	resp, err := PostForm(uri, data, args...)
	if err != nil {
		return
	}
	return resp.Unmarshal(d)
}
func PostFormWithFiles(uri string, data url.Values, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.PostFormWithFiles(context.Background(), uri, data)
}
func PostFormWithFilesUnmarshal(uri string, data url.Values, d any, args ...ArgsFunc) (err error) {
	resp, err := PostFormWithFiles(uri, data, args...)
	if err != nil {
		return
	}
	return resp.Unmarshal(d)
}
func Request(method, uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.DoRequest(context.Background(), method, uri, data)
}

func RequestUnmarshal(method, uri string, data, d any, args ...ArgsFunc) (err error) {
	resp, err := Request(method, uri, data, args...)
	if err != nil {
		return
	}
	return resp.Unmarshal(d)
}
