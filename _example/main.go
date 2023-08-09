package main

import (
	"fmt"
	"github.com/pkg6/go-requests"
	"os"
)

func main() {
	//data := url.Values{}
	//data.Set("k", "v")
	//get, _ := requests.New().Get(context.Background(), "http://www.httpbin.org/get", data)
	//defer get.Close()
	//fmt.Println(get.ReadAllString())
	//post, _ := requests.New().
	//	AsJson().
	//	Post(context.Background(), "http://www.httpbin.org/post", data)
	//defer post.Close()
	//fmt.Println(post.ReadAllString())
	f, _ := os.Open("./_example/github.png")
	fmt.Println(requests.Md5Reader(f))
}
