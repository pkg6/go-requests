package requests

import (
	"strings"
)

func CookieStringEscape(cookie string) Cookie {
	cookieRaw := make(Cookie)
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

type Cookie map[string][]string

func (v Cookie) Get(key string) string {
	if v == nil {
		return ""
	}
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}
func (v Cookie) Set(key, value string) {
	v[key] = []string{value}
}
func (v Cookie) Add(key, value string) {
	v[key] = append(v[key], value)
}
func (v Cookie) Del(key string) {
	delete(v, key)
}
func (v Cookie) Has(key string) bool {
	_, ok := v[key]
	return ok
}
func (v Cookie) Encode() string {
	cookieStr := ""
	for s := range v {
		if cookieStr != "" {
			cookieStr += ";"
		}
		cookieStr += s + "=" + v.Get(s)
	}
	return cookieStr
}
