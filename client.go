package requests

import (
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	httpSchemeName                 = `http`
	httpParamFileHolder            = `@file:`
	httpRegexParamJson             = `^[\w\[\]]+=.+`
	httpHeaderHost                 = `Host`
	httpHeaderCookie               = `Cookie`
	HttpHeaderUserAgent            = `User-Agent`
	HttpHeaderContentType          = `Content-Type`
	HttpHeaderContentTypePlainText = `text/plain; charset=utf-8`
	HttpHeaderContentTypeJson      = `application/json`
	HttpHeaderContentTypeXml       = `application/xml`
	HttpHeaderContentTypeForm      = `application/x-www-form-urlencoded`
)

var (
	hdrUserAgentKey = http.CanonicalHeaderKey(HttpHeaderUserAgent)

	hostname, _        = os.Hostname()
	defaultClientAgent = fmt.Sprintf(`pkg6/go-request client at  %s`, hostname)
	defaultRetryCount  = 3
	defaultWaitTime    = time.Duration(2000) * time.Millisecond
	DefaultClient      = new(Client).Clone()
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
	debug                  bool
	header                 http.Header
	cookies                map[string]string
	prefix                 string
	middlewares            []MiddlewareFunc
	beforeRequestCallbacks []clientCallback
	afterRequestCallbacks  []requestCallback
	responseCallbacks      []responseCallback
	successHooks           []SuccessHook
	errorHooks             []ErrorHook
	panicHooks             []ErrorHook
	retryCount             int
	retryWaitTime          time.Duration
	log                    *log.Logger
	jsonMarshal            func(v any) ([]byte, error)
	jsonUnmarshal          func(data []byte, v any) error
	xmlMarshal             func(v any) ([]byte, error)
	xmlUnmarshal           func(data []byte, v any) error
	lock                   sync.RWMutex
	traceContext           traceContext
	trace                  bool
	attempt                int
}

// HttpClient
//set InsecureSkipVerify = false c.Client.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = false
func HttpClient(localAddr net.Addr) *http.Client {
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
	return new(Client).Clone()
}

func NewHttpClient(client *http.Client) *Client {
	c := &Client{}
	c.Client = client
	return c
}

// Clone
//Parameter initialization
func (c *Client) Clone() *Client {
	if c.Client == nil {
		c.Client = HttpClient(nil)
	}
	c.debug = false
	c.header = make(http.Header, 0)
	c.cookies = make(map[string]string, 0)
	c.retryWaitTime = defaultWaitTime
	c.retryCount = defaultRetryCount
	c.beforeRequestCallbacks = make([]clientCallback, 0)
	c.afterRequestCallbacks = make([]requestCallback, 0)
	c.responseCallbacks = make([]responseCallback, 0)
	c.successHooks = make([]SuccessHook, 0)
	c.errorHooks = make([]ErrorHook, 0)
	if c.jsonMarshal == nil {
		c.SetJSONMarshaler(json.Marshal)
	}
	if c.jsonUnmarshal == nil {
		c.SetJSONUnmarshaler(json.Unmarshal)
	}
	if c.xmlMarshal == nil {
		c.SetXMLMarshaler(xml.Marshal)
	}
	if c.xmlUnmarshal == nil {
		c.SetXMLUnmarshaler(xml.Unmarshal)
	}
	if c.log == nil {
		c.log = log.Default()
	}
	if c.header.Get(HttpHeaderUserAgent) == "" {
		c.WithUserAgent(defaultClientAgent)
	}
	c.OnAfterRequest(requestLogger)
	c.OnResponse(responseLogger)
	c.attempt = 1
	return c
}
func (c *Client) Debug() *Client {
	c.debug = true
	return c
}

// SetPrefix sets the request server URL prefix.
func (c *Client) SetPrefix(prefix string) *Client {
	c.prefix = prefix
	return c
}

// SetJSONMarshaler method sets the JSON marshaler function to marshal the request body.
// By default, Resty uses `encoding/json` package to marshal the request body.
func (c *Client) SetJSONMarshaler(marshaler func(v any) ([]byte, error)) *Client {
	c.jsonMarshal = marshaler
	return c
}

// SetJSONUnmarshaler method sets the JSON unmarshaler function to unmarshal the response body.
// By default, Resty uses `encoding/json` package to unmarshal the response body.
func (c *Client) SetJSONUnmarshaler(unmarshaler func(data []byte, v any) error) *Client {
	c.jsonUnmarshal = unmarshaler
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
