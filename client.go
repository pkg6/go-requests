package requests

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"github.com/pkg6/go-requests/jsons"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
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
)

var (
	hdrUserAgentKey    = http.CanonicalHeaderKey(HttpHeaderUserAgent)
	hostname, _        = os.Hostname()
	defaultClientAgent = fmt.Sprintf(`pkg6/go-request client at  %s`, hostname)
	defaultRetryCount  = 3
	defaultWaitTime    = time.Duration(2000) * time.Millisecond
	defaultRequest     = New()
)

type (
	clientCallback   func(client *Client) error
	requestCallback  func(client *Client, request *http.Request) error
	responseCallback func(client *Client, request *http.Request, response *Response) error
	ErrorHook        func(client *Client, request *http.Request, err error)
	SuccessHook      func(client *Client, response *Response)
)

type Client struct {
	*http.Client
	BaseUrl                string
	Debug                  bool
	Query                  url.Values
	Logger                 *log.Logger
	header                 http.Header
	cookies                map[string]string
	jsonMarshal            func(v any) ([]byte, error)
	jsonUnmarshal          func(data []byte, v any) error
	xmlMarshal             func(v any) ([]byte, error)
	xmlUnmarshal           func(data []byte, v any) error
	middlewares            []MiddlewareFunc
	beforeRequestCallbacks []clientCallback
	afterRequestCallbacks  []requestCallback
	responseCallbacks      []responseCallback
	successHooks           []SuccessHook
	errorHooks             []ErrorHook
	panicHooks             []ErrorHook
	retryCount             int
	retryWaitTime          time.Duration
	lock                   sync.RWMutex
	traceContext           traceContext
	trace                  bool
	attempt                int
}

// defaultHttpClient
//set InsecureSkipVerify = false c.Client.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = false
func defaultHttpClient(localAddr net.Addr) *http.Client {
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

func New() *Client {
	return NewHttpClient(defaultHttpClient(nil))
}

func NewHttpClient(client *http.Client) *Client {
	return new(Client).Clone().WitchHttpClient(client)
}

// Clone
//Parameter initialization
func (c *Client) Clone() *Client {
	c.Debug = false
	c.Query = nil
	c.BaseUrl = ""
	c.header = make(http.Header, 0)
	c.cookies = make(map[string]string, 0)
	c.retryWaitTime = defaultWaitTime
	c.retryCount = defaultRetryCount
	c.beforeRequestCallbacks = make([]clientCallback, 0)
	c.afterRequestCallbacks = make([]requestCallback, 0)
	c.responseCallbacks = make([]responseCallback, 0)
	c.successHooks = make([]SuccessHook, 0)
	c.errorHooks = make([]ErrorHook, 0)
	//refer to https://github.com/gin-gonic/gin/tree/master/internal/json
	c.jsonMarshal = jsons.Marshal
	c.jsonUnmarshal = jsons.Unmarshal
	if c.xmlMarshal == nil {
		c.SetXMLMarshaler(xml.Marshal)
	}
	if c.xmlUnmarshal == nil {
		c.SetXMLUnmarshaler(xml.Unmarshal)
	}
	if c.Logger == nil {
		c.Logger = log.Default()
	}
	if c.header.Get(HttpHeaderUserAgent) == "" {
		c.WithUserAgent(defaultClientAgent)
	}
	c.OnAfterRequest(requestLogger)
	c.OnResponse(responseLogger)
	c.attempt = 1
	return c
}

func (c *Client) WitchHttpClient(client *http.Client) *Client {
	c.Client = client
	return c
}

func (c *Client) SetBaseURL(baseUrl string) *Client {
	c.BaseUrl = baseUrl
	return c
}

func (c *Client) SetQuery(query url.Values) *Client {
	c.Query = query
	return c
}

func (c *Client) SetDebug(debug bool) *Client {
	c.Debug = debug
	return c
}

// SetXMLMarshaler method sets the XML marshaler function to marshal the request body.
// By default, Resty uses `encoding/xml` package to marshal the request body.
func (c *Client) SetXMLMarshaler(marshaler func(v any) ([]byte, error)) *Client {
	c.xmlMarshal = marshaler
	return c
}

// SetXMLUnmarshaler method sets the XML unmarshaler function to unmarshal the response body.
// By default, Resty uses `encoding/xml` package to unmarshal the response body.
func (c *Client) SetXMLUnmarshaler(unmarshaler func(data []byte, v any) error) *Client {
	c.xmlUnmarshal = unmarshaler
	return c
}

// SetRetry is a chaining function,
// which sets retry count and interval when failure for next request.
func (c *Client) SetRetry(retryCount int, retryWaitTime time.Duration) *Client {
	c.retryCount = retryCount
	c.retryWaitTime = retryWaitTime
	return c
}

// Unmarshal content into object from JSON or XML
func (c *Client) Unmarshal(contentType string, b []byte, d any) (err error) {
	if IsJSONType(contentType) {
		err = c.jsonUnmarshal(b, d)
	} else if IsXMLType(contentType) {
		err = c.xmlUnmarshal(b, d)
	}
	return
}
