package requests

import (
	"context"
	"net/url"
)

//requests helper function

func Get(uri string, data any, args ...func(client *Client)) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Get(context.Background(), uri, data)
}

func Put(uri string, data any, args ...func(client *Client)) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Put(context.Background(), uri, data)
}

func Delete(uri string, data any, args ...func(client *Client)) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Delete(context.Background(), uri, data)
}

func Head(uri string, data any, args ...func(client *Client)) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Head(context.Background(), uri, data)
}

func Patch(uri string, data any, args ...func(client *Client)) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Patch(context.Background(), uri, data)
}

func Connect(uri string, data any, args ...func(client *Client)) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Connect(context.Background(), uri, data)
}

func Options(uri string, data any, args ...func(client *Client)) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Options(context.Background(), uri, data)
}

func Trace(uri string, data any, args ...func(client *Client)) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Trace(context.Background(), uri, data)
}

func Post(uri string, data any, args ...func(client *Client)) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.Post(context.Background(), uri, data)
}

func PostJson(uri string, data any, args ...func(client *Client)) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.PostJson(context.Background(), uri, data)
}

func PostForm(uri string, data url.Values, args ...func(client *Client)) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.PostForm(context.Background(), uri, data)
}

func Request(method, uri string, data any, args ...func(client *Client)) (*Response, error) {
	client := new(Client).Clone().WitchHttpClient(defaultHttpClient(nil))
	for _, arg := range args {
		arg(client)
	}
	return client.DoRequest(context.Background(), method, uri, data)
}
