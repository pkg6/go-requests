package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pkg6/go-requests"
)

func main() {
	data := url.Values{}
	data.Set("k", "v")
	get, _ := requests.New().Get(context.Background(), "http://www.httpbin.org/get", data)
	defer get.Close()
	fmt.Println(get.ReadAllString())
	post, _ := requests.New().
		AsJson().
		Post(context.Background(), "http://www.httpbin.org/post", data)
	defer post.Close()
	fmt.Println(post.ReadAllString())
}
