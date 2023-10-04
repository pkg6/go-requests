package requests

import (
	"context"
	"crypto/tls"
	"io"
	"net/url"
	"os"
)

//requests helper function

type ArgsFunc func(client *Client)

func Debug() ArgsFunc {
	return func(client *Client) {
		client.Debug = true
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
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.Get(context.Background(), uri, data)
}

func GetUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.GetD(context.Background(), uri, data, d)
}

func Put(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.Put(context.Background(), uri, data)
}
func PutUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.PutD(context.Background(), uri, data, d)
}
func Delete(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.Delete(context.Background(), uri, data)
}
func DeleteUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.DeleteD(context.Background(), uri, data, d)
}
func Head(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.Head(context.Background(), uri, data)
}
func HeadUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.HeadD(context.Background(), uri, data, d)
}
func Patch(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.Patch(context.Background(), uri, data)
}
func PatchUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.PatchD(context.Background(), uri, data, d)
}
func Connect(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.Connect(context.Background(), uri, data)
}
func ConnectUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.ConnectD(context.Background(), uri, data, d)
}
func Options(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.Options(context.Background(), uri, data)
}
func OptionsUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.OptionsD(context.Background(), uri, data, d)
}
func Trace(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.Trace(context.Background(), uri, data)
}
func TraceUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.TraceD(context.Background(), uri, data, d)
}
func Post(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.Post(context.Background(), uri, data)
}
func PostUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.PostD(context.Background(), uri, data, d)
}
func PostJson(uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.PostJson(context.Background(), uri, data)
}
func PostJsonUnmarshal(uri string, data, d any, args ...ArgsFunc) error {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.PostJsonD(context.Background(), uri, data, d)
}
func PostForm(uri string, data url.Values, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.PostForm(context.Background(), uri, data)
}
func PostFormUnmarshal(uri string, data url.Values, d any, args ...ArgsFunc) error {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.PostFormD(context.Background(), uri, data, d)
}
func PostFormWithFiles(uri string, data url.Values, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.PostFormWithFiles(context.Background(), uri, data)
}
func PostFormWithFilesUnmarshal(uri string, data url.Values, d any, args ...ArgsFunc) error {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.PostFormWithFilesD(context.Background(), uri, data, d)
}
func Request(method, uri string, data any, args ...ArgsFunc) (*Response, error) {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.DoRequest(context.Background(), method, uri, data)
}

func RequestUnmarshal(method, uri string, data, d any, args ...ArgsFunc) error {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	return client.DoRequestD(context.Background(), method, uri, data, d)
}

// RequestUnmarshalSaveFile
//	var resp _testdata.GitHubUser
//	err := requests.RequestUnmarshalSaveFile(http.MethodGet,
//		"https://api.github.com/users/github",
//		nil, &resp, "github_user.log")
//	fmt.Println(err)
//	fmt.Println(resp)
func RequestUnmarshalSaveFile(method, uri string, data, d any, fileName string, args ...ArgsFunc) error {
	client := new(Client).Clone()
	for _, arg := range args {
		arg(client)
	}
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	client.SetWriter(file)
	resp, err := client.DoRequest(context.Background(), method, uri, data)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, _ := io.ReadAll(resp.Body)
	if err = client.Unmarshal(resp.ContentType(), body, d); err != nil {
		return err
	}
	return nil
}
