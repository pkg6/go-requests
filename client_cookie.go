package requests

import (
	"net/http"
	"net/http/cookiejar"
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
	c.lock.Lock()
	defer c.lock.Unlock()
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
	c.lock.Lock()
	defer c.lock.Unlock()
	if cookie == "" {
		cookie = MapCookiesToString(c.cookies, cookie)
	}
	if cookie != "" {
		c.WithHeader(HttpHeaderCookie, cookie)
	}
	return c
}
