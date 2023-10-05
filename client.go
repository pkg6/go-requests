package requests

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	httpSchemeName      = `http`
	httpParamFileHolder = `@file:`
	httpRegexParamJson  = `^[\w\[\]]+=.+`

	HttpHeaderHost   = `Host`
	HttpHeaderCookie = `Cookie`

	HttpHeaderUserAgent     = `User-Agent`
	HttpHeaderAuthorization = "Authorization"
	HttpHeaderAccept        = "Accept"
	HttpMIMEEventStream     = "text/event-stream"
	HttpHeaderCacheControl  = "Cache-Control"
	HttpHeaderConnection    = "Connection"
	HttpHeaderContentType   = `Content-Type`

	charsetUTF8                          = "charset=UTF-8"
	HttpHeaderContentTypeJson            = `application/json`
	HttpHeaderContentTypeJsonCharsetUTF8 = HttpHeaderContentTypeJson + "; " + charsetUTF8
	HttpHeaderContentTypeXml             = `application/xml`
	HttpHeaderContentTypeXmlCharsetUTF8  = HttpHeaderContentTypeXml + "; " + charsetUTF8
	HttpHeaderContentTypeForm            = `application/x-www-form-urlencoded`

	AuthorizationTypeBearer = "Bearer "
	AuthorizationTypeBasic  = "Basic "
)

var (
	hdrUserAgentKey    = http.CanonicalHeaderKey(HttpHeaderUserAgent)
	hostname, _        = os.Hostname()
	defaultClientAgent = fmt.Sprintf(`github.com/pkg6/go-requests at  %s`, hostname)
	defaultRetryCount  = 3
	defaultWaitTime    = time.Duration(2000) * time.Millisecond
)

type Client struct {
	*http.Client
	Debug bool

	BaseUrl       string
	Query         url.Values
	Header        http.Header
	Cookie        Cookie
	Logger        LoggerInterface
	JSONMarshal   func(v any) ([]byte, error)
	JSONUnmarshal func(data []byte, v any) error
	XMLMarshal    func(v any) ([]byte, error)
	XMLUnmarshal  func(data []byte, v any) error

	//os.Stderr
	writer io.Writer

	middlewares            []MiddlewareFunc
	beforeRequestCallbacks []ClientCallback
	afterRequestCallbacks  []RequestCallback
	responseCallbacks      []ResponseCallback
	successHooks           []SuccessHook
	errorHooks             []ErrorHook
	panicHooks             []ErrorHook

	retryCount    int
	retryWaitTime time.Duration
	attempt       int

	traceContext traceContext
	trace        bool

	clone int
	lock  sync.RWMutex
	ctx   context.Context
}

// DefaultHttpClient
//set InsecureSkipVerify = false c.Client.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = false
func DefaultHttpClient(localAddr net.Addr) *http.Client {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	if localAddr != nil {
		dialer.LocalAddr = localAddr
	}
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableKeepAlives: true,
	}
	return &http.Client{Transport: transport}
}

func New() ClientInterface {
	return NewClient()
}

func NewWithHttpClient(client *http.Client) ClientInterface {
	return new(Client).SetHttpClient(client).Clone()
}

func NewClient() *Client {
	c := new(Client)
	c.Clone()
	return c
}

// Clone
//Parameter initialization
func (c *Client) Clone() ClientInterface {
	c.Debug = false
	if c.Client == nil {
		c.Client = DefaultHttpClient(nil)
	}
	c.BaseUrl = ""
	c.Query = make(url.Values, 0)
	c.Header = make(http.Header, 0)
	c.Cookie = make(Cookie, 0)

	if c.Logger == nil {
		c.SetLogger(DefaultLogger())
	}
	if c.JSONMarshal == nil {
		c.SetJSONMarshaler(json.Marshal)
	}
	if c.JSONUnmarshal == nil {
		c.SetJSONUnmarshaler(json.Unmarshal)
	}
	if c.XMLMarshal == nil {
		c.SetXMLMarshaler(xml.Marshal)
	}
	if c.XMLUnmarshal == nil {
		c.SetXMLUnmarshaler(xml.Unmarshal)
	}

	c.middlewares = make([]MiddlewareFunc, 0)
	c.beforeRequestCallbacks = make([]ClientCallback, 0)
	c.afterRequestCallbacks = make([]RequestCallback, 0)
	c.responseCallbacks = make([]ResponseCallback, 0)
	c.successHooks = make([]SuccessHook, 0)
	c.errorHooks = make([]ErrorHook, 0)
	c.panicHooks = make([]ErrorHook, 0)

	c.retryCount = defaultRetryCount
	c.retryWaitTime = defaultWaitTime

	c.traceContext = traceContext{}
	c.trace = false

	if c.ctx == nil {
		c.ctx = context.Background()
	}
	if c.Header.Get(HttpHeaderUserAgent) == "" {
		c.WithUserAgent(defaultClientAgent)
	}
	c.writer = nil
	c.OnAfterRequest(requestLogger)
	c.OnResponse(responseLogger)
	c.OnResponse(writerRequestResponseLog)
	c.attempt = 1
	c.clone += 1
	return c
}

func (c *Client) SetHttpClient(client *http.Client) ClientInterface {
	c.Client = client
	return c
}
func (c *Client) SetDebug(debug bool) ClientInterface {
	c.Debug = debug
	return c
}
func (c *Client) SetLogger(logger LoggerInterface) ClientInterface {
	c.Logger = logger
	return c
}
func (c *Client) SetWriter(writer io.Writer) ClientInterface {
	c.writer = writer
	return c
}

func (c *Client) SetBaseURL(baseUrl string) ClientInterface {
	c.BaseUrl = baseUrl
	return c
}

func (c *Client) SetQuery(query url.Values) ClientInterface {
	c.Query = query
	return c
}

func (c *Client) SetCookie(cookie Cookie) ClientInterface {
	c.Cookie = cookie
	return c
}
func (c *Client) SetHeader(header http.Header) ClientInterface {
	c.Header = header
	return c
}

// SetJSONMarshaler method sets the JSON marshaler function to marshal the request body.
// By default,  uses `encoding/json` package to marshal the request body.
func (c *Client) SetJSONMarshaler(marshaler func(v interface{}) ([]byte, error)) ClientInterface {
	c.JSONMarshal = marshaler
	return c
}

// SetJSONUnmarshaler method sets the JSON unmarshaler function to unmarshal the response body.
func (c *Client) SetJSONUnmarshaler(unmarshaler func(data []byte, v interface{}) error) ClientInterface {
	c.JSONUnmarshal = unmarshaler
	return c
}

// SetXMLMarshaler method sets the XML marshaler function to marshal the request body.
func (c *Client) SetXMLMarshaler(marshaler func(v any) ([]byte, error)) ClientInterface {
	c.XMLMarshal = marshaler
	return c
}

// SetXMLUnmarshaler method sets the XML unmarshaler function to unmarshal the response body.
// By default,  uses `encoding/xml` package to unmarshal the response body.
func (c *Client) SetXMLUnmarshaler(unmarshaler func(data []byte, v any) error) ClientInterface {
	c.XMLUnmarshal = unmarshaler
	return c
}

// SetRetry is a chaining function,
// which sets retry count and interval when failure for next request.
func (c *Client) SetRetry(retryCount int, retryWaitTime time.Duration) ClientInterface {
	c.retryCount = retryCount
	c.retryWaitTime = retryWaitTime
	return c
}

// SetTimeout sets the request timeout for the client.
func (c *Client) SetTimeout(t time.Duration) ClientInterface {
	c.Client.Timeout = t
	return c
}

// WithProxyUrl set proxy for the client.
// This func will do nothing when the parameter `proxyURL` is empty or in wrong pattern.
// The correct pattern is like `http://USER:PASSWORD@IP:PORT` or `socks5://USER:PASSWORD@IP:PORT`.
// Only `http` and `socks5` proxies are supported currently.
func (c *Client) WithProxyUrl(proxyURL string) ClientInterface {
	if strings.TrimSpace(proxyURL) == "" {
		return c
	}
	_proxy, err := url.Parse(proxyURL)
	if err != nil {
		c.Logger.Errorf(`%+v`, err)
		return c
	}
	if _proxy.Scheme == httpSchemeName {
		if v, ok := c.Transport.(*http.Transport); ok {
			v.Proxy = http.ProxyURL(_proxy)
		}
	} else {
		auth := &proxy.Auth{}
		user := _proxy.User.Username()
		if user != "" {
			auth.User = user
			password, hasPassword := _proxy.User.Password()
			if hasPassword && password != "" {
				auth.Password = password
			}
		} else {
			auth = nil
		}
		// refer to the source code, error is always nil
		dialer, err := proxy.SOCKS5("tcp", _proxy.Host, auth, &net.Dialer{
			Timeout:   c.Client.Timeout,
			KeepAlive: c.Client.Timeout,
		})
		if err != nil {
			c.Logger.Errorf(`%+v`, err)
			return c
		}
		if v, ok := c.Transport.(*http.Transport); ok {
			v.DialContext = func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return dialer.Dial(network, addr)
			}
		}
	}
	return c
}

// WithTLSKeyCrt sets the certificate and key file for TLS configuration of client.
func (c *Client) WithTLSKeyCrt(crtFile, keyFile string) ClientInterface {
	crt, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		c.Logger.Errorf("LoadKeyCrt failed")
		return c
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	tlsConfig.Time = time.Now
	tlsConfig.Rand = rand.Reader
	tlsConfig.InsecureSkipVerify = true
	c.SetTLSConfig(tlsConfig)
	return c
}

// SetTLSConfig sets the TLS configuration of client.
func (c *Client) SetTLSConfig(tlsConfig *tls.Config) ClientInterface {
	v, ok := c.Transport.(*http.Transport)
	if !ok {
		c.Logger.Errorf(`cannot set TLSClientConfig for custom Transport of the client`)
		return c
	}
	v.TLSClientConfig = tlsConfig
	return c
}

func (c *Client) SetCheckRedirect(fn func(req *http.Request, via []*http.Request) error) {
	c.CheckRedirect = fn
}

// Unmarshal content into object from JSON or XML
func (c *Client) Unmarshal(contentType string, b []byte, d any) (err error) {
	if IsJSONType(contentType) {
		err = c.JSONUnmarshal(b, d)
	} else if IsXMLType(contentType) {
		err = c.XMLUnmarshal(b, d)
	}
	return
}
