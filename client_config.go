package requests

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// BrowserMode enables browser mode of the client.
// When browser mode is enabled, it automatically saves and sends cookie content
// from and to server.
func (c *Client) BrowserMode() *Client {
	c.Jar, _ = cookiejar.New(nil)
	return c
}
func (c *Client) WithCookieJar(jar http.CookieJar) *Client {
	c.Client.Jar = jar
	return c
}
func (c *Client) WithCookie(k, v string) *Client {
	c.cookies[k] = v
	return c
}
func (c *Client) WithCookies(cookies map[string]string) *Client {
	for k, v := range cookies {
		c.WithCookie(k, v)
	}
	return c
}
func (c *Client) WithCookieString(cookie string) *Client {
	if cookie == "" {
		cookie = MapCookiesToString(c.cookies, cookie)
	}
	if cookie != "" {
		c.WithHeader(httpHeaderCookie, cookie)
	}
	return c
}

// WithHeader method sets a single header field and its value in the client instance.
// These headers will be applied to all requests raised from this client instance.
// Also it can be overridden at request level header options.
//		WithHeader("Content-Type", "application/json").
//		WithHeader("Accept", "application/json")
func (c *Client) WithHeader(header, value string) *Client {
	c.header.Set(header, value)
	return c
}

// WithHeaders method sets multiple headers field and its values at one go in the client instance.
// These headers will be applied to all requests raised from this client instance. Also it can be
// overridden at request level headers options.
// For Example: To set `Content-Type` and `Accept` as `application/json`
//		WithHeaders(map[string]string{
//			"Content-Type": "application/json",
//			"Accept": "application/json",
//		})
func (c *Client) WithHeaders(headers map[string]string) *Client {
	for h, v := range headers {
		c.header.Set(h, v)
	}
	return c
}

// WithHeaderVerbatim method is to set a single header field and its value verbatim in the current request.
//
// For Example: To set `all_lowercase` and `UPPERCASE` as `available`.
//		WithHeaderVerbatim("all_lowercase", "available").
//		WithHeaderVerbatim("UPPERCASE", "available")
//
// Also you can override header value, which was set at client instance level.
func (c *Client) WithHeaderVerbatim(header, value string) *Client {
	c.header[header] = []string{value}
	return c
}

// WithContentType is a chaining function,
// which sets HTTP content type for the next request.
func (c *Client) WithContentType(contentType string) *Client {
	c.WithHeader(HttpHeaderContentType, contentType)
	return c
}
func (c *Client) WithUserAgent(userAgent string) *Client {
	c.WithHeader(HttpHeaderUserAgent, userAgent)
	return c
}

// AsForm is a chaining function,
// which sets the HTTP content type as "application/x-www-form-urlencoded" for the next request.
func (c *Client) AsForm() *Client {
	c.WithContentType(HttpHeaderContentTypeForm)
	return c
}

// AsStream is a Stream
func (c *Client) AsStream() *Client {
	c.WithHeaders(map[string]string{
		"Accept":        "text/event-stream",
		"Cache-Control": "no-cache",
		"Connection":    "keep-alive",
	})
	return c
}

// AsJson is a chaining function,
// which sets the HTTP content type as "application/json" for the next request.
//
// Note that it also checks and encodes the parameter to JSON format automatically.
func (c *Client) AsJson() *Client {
	c.WithContentType(HttpHeaderContentTypeJson)
	return c
}

// AsXml is a chaining function,
// which sets the HTTP content type as "application/xml" for the next request.
//
// Note that it also checks and encodes the parameter to XML format automatically.
func (c *Client) AsXml() *Client {
	c.WithContentType(HttpHeaderContentTypeXml)
	return c
}

// WithBasicAuth
//Specify the basic authentication username and password for the request.
func (c *Client) WithBasicAuth(username, password string) *Client {
	c.WithToken(base64.StdEncoding.EncodeToString([]byte(username+":"+password)), "Basic ")
	return c
}

// WithToken
//Specify an authorization token for the request.
func (c *Client) WithToken(token string, Type ...string) *Client {
	if len(Type) > 0 {
		token = Type[0] + token
	} else {
		token = "Bearer " + token
	}
	c.header.Set("Authorization", token)
	return c
}

// Timeout sets the request timeout for the client.
func (c *Client) Timeout(t time.Duration) *Client {
	c.Client.Timeout = t
	return c
}

// WithProxyUrl set proxy for the client.
// This func will do nothing when the parameter `proxyURL` is empty or in wrong pattern.
// The correct pattern is like `http://USER:PASSWORD@IP:PORT` or `socks5://USER:PASSWORD@IP:PORT`.
// Only `http` and `socks5` proxies are supported currently.
func (c *Client) WithProxyUrl(proxyURL string) *Client {
	if strings.TrimSpace(proxyURL) == "" {
		return c
	}
	_proxy, err := url.Parse(proxyURL)
	if err != nil {
		c.log.Fatalf(`%+v`, err)
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
		dialer, err := proxy.SOCKS5(
			"tcp",
			_proxy.Host,
			auth,
			&net.Dialer{
				Timeout:   c.Client.Timeout,
				KeepAlive: c.Client.Timeout,
			},
		)
		if err != nil {
			c.log.Fatalf(`%+v`, err)
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
func (c *Client) WithTLSKeyCrt(crtFile, keyFile string) *Client {
	crt, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		c.log.Fatalf("LoadKeyCrt failed")
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
func (c *Client) SetTLSConfig(tlsConfig *tls.Config) *Client {
	v, ok := c.Transport.(*http.Transport)
	if !ok {
		c.log.Fatalf(`cannot set TLSClientConfig for custom Transport of the client`)
	}
	v.TLSClientConfig = tlsConfig
	return c
}

// WithRedirectLimit limits the number of jumps.
func (c *Client) WithRedirectLimit(redirectLimit int) *Client {
	c.WithRedirectPolicy(func(req *http.Request, via []*http.Request) error {
		if len(via) >= redirectLimit {
			return http.ErrUseLastResponse
		}
		return nil
	})
	return c
}

// WithRedirectPolicy method sets the client redirect poilicy. Resty provides ready to use
// redirect policies. Wanna create one for yourself refer to `redirect.go`.
//	WithRedirectLimit(20)
//	WithRedirectPolicy(FlexibleRedirectPolicy(20))
//	WithRedirectPolicy(FlexibleRedirectPolicy(20), DomainCheckRedirectPolicy("host1.com", "host2.net"))
func (c *Client) WithRedirectPolicy(policies ...any) *Client {
	if len(policies) == 1 {
		if checkRedirect, ok := policies[0].(func(req *http.Request, via []*http.Request) error); ok {
			c.WithCheckRedirect(checkRedirect)
			return c
		}
	}
	c.WithCheckRedirect(func(req *http.Request, via []*http.Request) error {
		for _, p := range policies {
			if _, ok := p.(RedirectPolicy); ok {
				if err := p.(RedirectPolicy).Apply(req, via); err != nil {
					return err
				}
			} else {
				c.log.Fatalf("%v does not implement resty.RedirectPolicy (missing Apply method)", functionName(p))
			}
		}
		// looks good, go ahead
		return nil
	})
	return c
}

func (c *Client) WithCheckRedirect(fn func(req *http.Request, via []*http.Request) error) {
	c.CheckRedirect = fn
}
