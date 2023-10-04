package requests

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

// BrowserMode enables browser mode of the client.
// When browser mode is enabled, it automatically saves and sends cookie content
// from and to server.
func (c *Client) BrowserMode() *Client {
	jar, _ := cookiejar.New(nil)
	return c.WithCookieJar(jar)
}
func (c *Client) WithCookieJar(jar http.CookieJar) *Client {
	c.Client.Jar = jar
	return c
}
func (c *Client) WithCookieString(cookieString string) *Client {
	c.WithCookies(CookieStringEscape(cookieString))
	return c
}

func (c *Client) WithCookie(k, v string) *Client {
	c.Cookie.Set(k, v)
	return c
}

func (c *Client) WithCookieMap(cookies map[string]string) *Client {
	for k, v := range cookies {
		c.WithCookie(k, v)
	}
	return c
}

func (c *Client) WithCookies(cookie Cookie) *Client {
	c.Cookie = cookie
	return c
}

// WithCookieNextRequest
//  The first access to the root domain name will cache cookie data, and the second access will carry the cookie data from the cache until the cache expires and is regenerated
//	cache := requests.NewFileCache("you path/cache")
//	WithCookieNextRequest(cache, time.Hour)
func (c *Client) WithCookieNextRequest(cache ICache, ttl time.Duration) *Client {
	//set cookie
	c.OnResponse(OnResponseWithCookie(cache, ttl))
	// get cookie
	c.OnAfterRequest(OnAfterRequestWithCookie(cache))
	return c
}
