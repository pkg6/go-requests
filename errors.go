package requests

import "fmt"

type RequestError struct {
	URI        string
	Method     string
	StatusCode int
	Response   *Response
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("Request %s %s statusCode:%d", r.Method, r.URI, r.StatusCode)
}

type ResponseError struct {
	Response *Response
	Err      error
}

func (e *ResponseError) Error() string {
	return e.Err.Error()
}
