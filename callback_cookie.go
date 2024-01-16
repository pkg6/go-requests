package requests

import (
	"net/http"
	"time"
)

// onAfterRequestWithCookieGet get cookie on request
func onAfterRequestWithCookieGet(cache CacheInterface) RequestCallback {
	return func(client *Client, request *http.Request) error {
		cacheKey := Md5(request.URL.Host)
		if cache.Has(cacheKey) {
			if cookStr, _ := cache.Get(cacheKey); cookStr != "" {
				var cookieRaw Cookie
				_ = client.JSONUnmarshal([]byte(cookStr), &cookieRaw)
				request.Header.Set(HttpHeaderCookie, cookieRaw.Encode())
			}
		}
		return nil
	}
}

// onResponseNextRequestWithCookieSet set url cookie on cache
func onResponseNextRequestWithCookieSet(cache CacheInterface, ttl time.Duration) ResponseCallback {
	return func(client *Client, request *http.Request, response *Response) error {
		cacheKey := Md5(request.URL.Host)
		if !cache.Has(cacheKey) {
			cookieRaw := response.GetCookie()
			if len(cookieRaw) > 0 {
				cookieByte, _ := client.JSONMarshal(cookieRaw)
				_ = cache.Set(cacheKey, string(cookieByte), ttl)
			}
		}
		return nil
	}
}
