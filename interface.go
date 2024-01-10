package requests

import (
	"net/http"
	"time"
)

type (
	ClientCallback   func(client *Client) error
	RequestCallback  func(client *Client, request *http.Request) error
	ResponseCallback func(client *Client, request *http.Request, response *Response) error
	ErrorHook        func(client *Client, request *http.Request, err error)
	SuccessHook      func(client *Client, response *Response)

	CtxKeyString string
)

type CacheInterface interface {
	Set(key, value string, ttl time.Duration) error
	Get(key string) (string, error)
	Has(key string) bool
	Delete(key string) error
	CleanExpired() error
}
type LoggerInterface interface {
	Errorf(format string, v ...any)
	Warnf(format string, v ...any)
	Debugf(format string, v ...any)
}
