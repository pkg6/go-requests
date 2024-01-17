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
		client.SetDebug(true)
	}
}
func WithCookie(k, v string) ArgsFunc {
	return func(client *Client) {
		client.WithCookie(k, v)
	}
}
func WithCookies(cookies map[string]string) ArgsFunc {
	return func(client *Client) {
		client.WithCookieMap(cookies)
	}
}

func WithHeader(header, value string) ArgsFunc {
	return func(client *Client) {
		client.WithHeader(header, value)
	}
}
func WithHeaders(headers map[string]string) ArgsFunc {
	return func(client *Client) {
		client.WithHeaderMap(headers)
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
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.Get(context.Background(), uri, data)
}
func GetD(uri string, data, d any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.GetD(context.Background(), uri, data, d)
}
func GetUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.GetUnmarshal(context.Background(), uri, data, d)
}

func Put(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.Put(context.Background(), uri, data)
}
func PutD(uri string, data, d any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PutD(context.Background(), uri, data, d)
}
func PutUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PutUnmarshal(context.Background(), uri, data, d)
}
func Delete(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.Delete(context.Background(), uri, data)
}
func DeleteD(uri string, data, d any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.DeleteD(context.Background(), uri, data, d)
}
func DeleteUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.DeleteUnmarshal(context.Background(), uri, data, d)
}
func Head(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.Head(context.Background(), uri, data)
}
func HeadD(uri string, data, d any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.HeadD(context.Background(), uri, data, d)
}
func HeadUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.HeadUnmarshal(context.Background(), uri, data, d)
}
func Patch(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.Patch(context.Background(), uri, data)
}
func PatchD(uri string, data, d any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PatchD(context.Background(), uri, data, d)
}
func PatchUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PatchUnmarshal(context.Background(), uri, data, d)
}
func Connect(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.Connect(context.Background(), uri, data)
}
func ConnectD(uri string, data, d any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.ConnectD(context.Background(), uri, data, d)
}
func ConnectUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.ConnectUnmarshal(context.Background(), uri, data, d)
}
func Options(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.Options(context.Background(), uri, data)
}
func OptionsD(uri string, data, d any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.OptionsD(context.Background(), uri, data, d)
}
func OptionsUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.OptionsUnmarshal(context.Background(), uri, data, d)
}
func Trace(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.Trace(context.Background(), uri, data)
}
func TraceD(uri string, data, d any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.TraceD(context.Background(), uri, data, d)
}
func TraceUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.TraceUnmarshal(context.Background(), uri, data, d)
}
func Post(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.Post(context.Background(), uri, data)
}
func PostD(uri string, data, d any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PostD(context.Background(), uri, data, d)
}
func PostUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PostUnmarshal(context.Background(), uri, data, d)
}
func PostJson(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PostJson(context.Background(), uri, data)
}
func PostJsonD(uri string, data, d any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PostJsonD(context.Background(), uri, data, d)
}
func PostJsonUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PostJsonUnmarshal(context.Background(), uri, data, d)
}
func PostForm(uri string, data url.Values, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PostForm(context.Background(), uri, data)
}
func PostFormD(uri string, data url.Values, d any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PostFormD(context.Background(), uri, data, d)
}
func PostFormUnmarshal(uri string, data url.Values, d any, args ...ArgsFunc) error {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PostFormUnmarshal(context.Background(), uri, data, d)
}
func PostFormWithFiles(uri string, data url.Values, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PostFormWithFiles(context.Background(), uri, data)
}
func PostFormWithFilesUnmarshal(uri string, data url.Values, d any, args ...ArgsFunc) error {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.PostFormWithFilesUnmarshal(context.Background(), uri, data, d)
}
func Request(method, uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.DoRequest(context.Background(), method, uri, data)
}
func RequestD(method, uri string, data, d any, args ...ArgsFunc) (*Response, error) {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.DoRequestD(context.Background(), method, uri, data, d)
}

func RequestUnmarshal(method, uri string, data, d any, args ...ArgsFunc) error {
	client := NewClient()
	for _, arg := range args {
		arg(client)
	}
	return client.DoRequestUnmarshal(context.Background(), method, uri, data, d)
}

func Download(uri, fileName string) error {
	return NewClient().Download(context.Background(), uri, fileName)
}
