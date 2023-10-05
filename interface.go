package requests

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"time"
)

type (
	ClientCallback   func(client *Client) error
	RequestCallback  func(client *Client, request *http.Request) error
	ResponseCallback func(client *Client, request *http.Request, response *Response) error
	ErrorHook        func(client *Client, request *http.Request, err error)
	SuccessHook      func(client *Client, response *Response)
)

type ClientInterface interface {
	ClientHttpClientClient
	ClientOwnerInterface
}

type ClientHttpClientClient interface {
	SetHttpClient(client *http.Client) ClientInterface
	Timeout(t time.Duration) ClientInterface
	WithProxyUrl(proxyURL string) ClientInterface
	WithTLSKeyCrt(crtFile, keyFile string) ClientInterface
	SetTLSConfig(tlsConfig *tls.Config) ClientInterface
	SetCheckRedirect(fn func(req *http.Request, via []*http.Request) error)
}

type ClientOwnerInterface interface {
	Clone() ClientInterface
	SetDebug(debug bool) ClientInterface
	SetRetry(retryCount int, retryWaitTime time.Duration) ClientInterface
	SetWriter(writer io.Writer) ClientInterface
	SetBaseURL(baseUrl string) ClientInterface
	SetQuery(query url.Values) ClientInterface
	SetCookie(cookie Cookie) ClientInterface
	SetHeader(header http.Header) ClientInterface
	ClientFnInterface
	ClientMiddlewareInterface
	ClientHeaderInterface
	ClientRequestInterface
}

type ClientFnInterface interface {
	SetJSONMarshaler(marshaler func(v interface{}) ([]byte, error)) ClientInterface
	SetJSONUnmarshaler(unmarshaler func(data []byte, v interface{}) error) ClientInterface
	SetXMLMarshaler(marshaler func(v any) ([]byte, error)) ClientInterface
	SetXMLUnmarshaler(unmarshaler func(data []byte, v any) error) ClientInterface
	Unmarshal(contentType string, b []byte, d any) (err error)

	OnBeforeRequest(callback ClientCallback) ClientInterface
	OnAfterRequest(callback RequestCallback) ClientInterface
	OnResponse(callback ResponseCallback) ClientInterface
	OnError(h ErrorHook) ClientInterface
	OnSuccess(h SuccessHook) ClientInterface
	OnPanic(h ErrorHook) ClientInterface
}

type ClientHeaderInterface interface {
	WithClientCookieJar(jar http.CookieJar) ClientInterface
	WithHeader(header, value string) ClientInterface
	WithHeaderMap(headers map[string]string) ClientInterface
	WithHeaderVerbatim(header, value string) ClientInterface
	WithContentType(contentType string) ClientInterface
	WithUserAgent(userAgent string) ClientInterface
	WithRandomUserAgent() ClientInterface
	WithRandomMobileUserAgent() ClientInterface
	AsForm() ClientInterface
	AsJson() ClientInterface
	AsXml() ClientInterface
	WithBasicAuth(username, password string) ClientInterface
	WithToken(token string, tokenType ...string) ClientInterface

	WithCookieString(cookieString string) ClientInterface
	WithCookie(k, v string) ClientInterface
	WithCookieMap(cookies map[string]string) ClientInterface
	WithCookieNextRequest(cache ICache, ttl time.Duration) ClientInterface

	WithRedirectPolicy(policies ...any) ClientInterface
	WithRedirectLimit(redirectLimit int) ClientInterface
}

type ClientMiddlewareInterface interface {
	Use(middlewares ...MiddlewareFunc) ClientInterface
}

type ClientRequestInterface interface {
	ClientDoRequestInterface
	Get(ctx context.Context, uri string, data any) (*Response, error)
	GetUnmarshal(ctx context.Context, uri string, data, d any) error
	GetBytes(ctx context.Context, uri string, data any) ([]byte, error)

	Put(ctx context.Context, uri string, data any) (*Response, error)
	PutUnmarshal(ctx context.Context, uri string, data, d any) error
	PutBytes(ctx context.Context, uri string, data any) ([]byte, error)

	Post(ctx context.Context, uri string, data any) (*Response, error)
	PostUnmarshal(ctx context.Context, uri string, data, d any) error
	PostBytes(ctx context.Context, uri string, data any) ([]byte, error)

	Delete(ctx context.Context, uri string, data any) (*Response, error)
	DeleteUnmarshal(ctx context.Context, uri string, data, d any) error
	DeleteBytes(ctx context.Context, uri string, data any) ([]byte, error)

	Head(ctx context.Context, uri string, data any) (*Response, error)
	HeadUnmarshal(ctx context.Context, uri string, data, d any) error
	HeadBytes(ctx context.Context, uri string, data any) ([]byte, error)

	Patch(ctx context.Context, uri string, data any) (*Response, error)
	PatchUnmarshal(ctx context.Context, uri string, data, d any) error
	PatchBytes(ctx context.Context, uri string, data any) ([]byte, error)

	Connect(ctx context.Context, uri string, data any) (*Response, error)
	ConnectUnmarshal(ctx context.Context, uri string, data, d any) error
	ConnectBytes(ctx context.Context, uri string, data any) ([]byte, error)

	Options(ctx context.Context, uri string, data any) (*Response, error)
	OptionsUnmarshal(ctx context.Context, uri string, data, d any) error
	OptionsBytes(ctx context.Context, uri string, data any) ([]byte, error)

	Trace(ctx context.Context, uri string, data any) (*Response, error)
	TraceUnmarshal(ctx context.Context, uri string, data, d any) error
	TraceBytes(ctx context.Context, uri string, data any) ([]byte, error)

	PostJson(ctx context.Context, uri string, data any) (*Response, error)
	PostJsonUnmarshal(ctx context.Context, uri string, data, d any) error
	PostJsonBytes(ctx context.Context, uri string, data any) ([]byte, error)

	PostForm(ctx context.Context, uri string, data url.Values) (*Response, error)
	PostFormUnmarshal(ctx context.Context, uri string, data url.Values, d any) error
	PostFormBytes(ctx context.Context, uri string, data url.Values) ([]byte, error)

	PostFormWithFiles(ctx context.Context, uri string, data url.Values) (*Response, error)
	PostFormWithFilesUnmarshal(ctx context.Context, uri string, data url.Values, d any) error
	PostFormWithFilesBytes(ctx context.Context, uri string, data url.Values) ([]byte, error)
}

type ClientDoRequestInterface interface {
	DoRequestUnmarshal(ctx context.Context, method string, uri string, data, d any) error
	DoRequestBytes(ctx context.Context, method string, uri string, data any) ([]byte, error)
	DoRequest(ctx context.Context, method, uri string, body any) (response *Response, err error)
}
