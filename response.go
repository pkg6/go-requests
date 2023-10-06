package requests

import (
	"bufio"
	"io"
	"net/http"
)

type Response struct {
	*http.Response               // Response is the underlying http.Response object of certain request.
	request        *http.Request // Request is the underlying http.Request object of certain request.
	client         *Client
}

// Close closes the response when it will never be used.
func (r *Response) Close() error {
	if r == nil || r.Response == nil {
		return nil
	}
	r.client = nil
	r.request = nil
	return r.Response.Body.Close()
}

func (r *Response) TraceInfo() TraceInfo {
	return r.client.getTraceInfo()
}

func (r *Response) GetCookie() Cookie {
	raw := make(Cookie)
	for _, c := range r.Cookies() {
		raw.Set(c.Name, c.Value)
	}
	return raw
}

// ReadAll retrieves and returns the response content as []byte.
func (r *Response) ReadAll() []byte {
	// Response might be nil.
	if r == nil || r.Response == nil {
		return []byte{}
	}
	body, err := io.ReadAll(r.Response.Body)
	if err != nil {
		return nil
	}
	return body
}

// ReadStream
//Microsoft ChatGPT Data Structure start
//data: {"id":"","object":"","created":0,"model":"","prompt_filter_results":[{"prompt_index":0,"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"choices":[],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"role":"assistant"},"content_filter_results":{}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"content":"我"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"content":"是"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"content":"一个"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"content":"AI"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"content":"助"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"content":"手"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"content":"，"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"content":"可以"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"content":"帮"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"content":"助"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"content":"您"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"content":"查"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"content":"找"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"content":"信息"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":null,"delta":{"content":"。"},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}}}],"usage":null}
//data: {"id":"chatcmpl-801dpnRDRT0u4Nx5tec615WTJOXJt","object":"chat.completion.chunk","created":1695017001,"model":"gpt-35-turbo","choices":[{"index":0,"finish_reason":"stop","delta":{},"content_filter_results":{}}],"usage":null}
//data: [DONE]
//Microsoft ChatGPT Data Structure end
func (r *Response) ReadStream(lineNumberFun func(line []byte, number int64)) int64 {
	var number int64
	if r.IsError() {
		return number
	}
	if !IsStreamType(r.ContentType()) {
		lineNumberFun(r.ReadAll(), number)
		return number
	}
	b := bufio.NewReader(r.Response.Body)
	for {
		rawLine, _, err := b.ReadLine()
		if err != nil {
			break
		}
		lineNumberFun(rawLine, number)
		number += 1
	}
	return number
}

// ReadAllString retrieves and returns the response content as string.
func (r *Response) ReadAllString() string {
	return string(r.ReadAll())
}

// ContentType response header Content-Type
func (r *Response) ContentType() string {
	return r.Response.Header.Get(HttpHeaderContentType)
}

// Unmarshal content into object from JSON or XML
func (r *Response) Unmarshal(d any) error {
	return r.client.Unmarshal(r.ContentType(), r.ReadAll(), d)
}

// IsSuccess method returns true if HTTP status `code >= 200 and <= 299` otherwise false.
func (r *Response) IsSuccess() bool {
	return r.StatusCode > 199 && r.StatusCode < 300
}

// IsError method returns true if HTTP status `code >= 400` otherwise false.
func (r *Response) IsError() bool {
	return r.StatusCode > 399
}
