package requests

import (
	"context"
	"net/url"
)

//requests helper function

func Get(uri string, data any, args ...func(client *Client)) (*Response, error) {
	requests := new(Client).Clone()
	for _, arg := range args {
		arg(requests)
	}
	return requests.Get(context.Background(), uri, data)
}

func Put(uri string, data any, args ...func(client *Client)) (*Response, error) {
	requests := new(Client).Clone()
	for _, arg := range args {
		arg(requests)
	}
	return requests.Put(context.Background(), uri, data)
}

func Delete(uri string, data any, args ...func(client *Client)) (*Response, error) {
	requests := new(Client).Clone()
	for _, arg := range args {
		arg(requests)
	}
	return requests.Delete(context.Background(), uri, data)
}

func Head(uri string, data any, args ...func(client *Client)) (*Response, error) {
	requests := new(Client).Clone()
	for _, arg := range args {
		arg(requests)
	}
	return requests.Head(context.Background(), uri, data)
}

func Patch(uri string, data any, args ...func(client *Client)) (*Response, error) {
	requests := new(Client).Clone()
	for _, arg := range args {
		arg(requests)
	}
	return requests.Patch(context.Background(), uri, data)
}

func Connect(uri string, data any, args ...func(client *Client)) (*Response, error) {
	requests := new(Client).Clone()
	for _, arg := range args {
		arg(requests)
	}
	return requests.Connect(context.Background(), uri, data)
}

func Options(uri string, data any, args ...func(client *Client)) (*Response, error) {
	requests := new(Client).Clone()
	for _, arg := range args {
		arg(requests)
	}
	return requests.Options(context.Background(), uri, data)
}

func Trace(uri string, data any, args ...func(client *Client)) (*Response, error) {
	requests := new(Client).Clone()
	for _, arg := range args {
		arg(requests)
	}
	return requests.Trace(context.Background(), uri, data)
}

func Post(uri string, data any, args ...func(client *Client)) (*Response, error) {
	requests := new(Client).Clone()
	for _, arg := range args {
		arg(requests)
	}
	return requests.Post(context.Background(), uri, data)
}

func PostJson(uri string, data any, args ...func(client *Client)) (*Response, error) {
	requests := new(Client).Clone()
	for _, arg := range args {
		arg(requests)
	}
	return requests.PostJson(context.Background(), uri, data)
}

func PostForm(uri string, data url.Values, args ...func(client *Client)) (*Response, error) {
	requests := new(Client).Clone()
	for _, arg := range args {
		arg(requests)
	}
	return requests.PostForm(context.Background(), uri, data)
}

func Request(method, uri string, data any, args ...func(client *Client)) (*Response, error) {
	requests := new(Client).Clone()
	for _, arg := range args {
		arg(requests)
	}
	return requests.DoRequest(context.Background(), method, uri, data)
}
