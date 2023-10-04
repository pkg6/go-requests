package requests

import (
	"net/http"
	"strings"
	"time"
)

func OnAfterRequestWithCookie(cache ICache) RequestCallback {
	return func(client *Client, request *http.Request) error {
		cacheKey := Md5(request.URL.Host)
		if cache.Has(cacheKey) {
			if cookStr, _ := cache.Get(cacheKey); cookStr != "" {
				var cookieRaw CookieRaw
				_ = client.JSONUnmarshal([]byte(cookStr), &cookieRaw)
				request.Header.Set(HttpHeaderCookie, cookieRaw.Encode())
			}
		}
		return nil
	}
}

func OnResponseWithCookie(cache ICache, ttl time.Duration) ResponseCallback {
	return func(client *Client, request *http.Request, response *Response) error {
		cacheKey := Md5(request.URL.Host)
		if !cache.Has(cacheKey) {
			cookieRaw := response.GetCookieMap()
			if len(cookieRaw) > 0 {
				cookieByte, _ := client.JSONMarshal(cookieRaw)
				_ = cache.Set(cacheKey, string(cookieByte), ttl)
			}
		}
		return nil
	}
}

func CookieToRaw(cookie string) CookieRaw {
	cookieRaw := make(CookieRaw)
	if cookie == "" {
		return cookieRaw
	}
	for _, pair := range strings.Split(cookie, ";") {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			cookieRaw.Set(kv[0], kv[1])
		}
	}
	return cookieRaw
}

type CookieRaw map[string][]string

func (v CookieRaw) Get(key string) string {
	if v == nil {
		return ""
	}
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}
func (v CookieRaw) Set(key, value string) {
	v[key] = []string{value}
}
func (v CookieRaw) Add(key, value string) {
	v[key] = append(v[key], value)
}
func (v CookieRaw) Del(key string) {
	delete(v, key)
}
func (v CookieRaw) Has(key string) bool {
	_, ok := v[key]
	return ok
}
func (v CookieRaw) Encode() string {
	cookieStr := ""
	for s := range v {
		if cookieStr != "" {
			cookieStr += ";"
		}
		cookieStr += s + "=" + v.Get(s)
	}
	return cookieStr
}
