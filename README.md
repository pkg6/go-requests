[![Test Status](https://github.com/pkg6/go-requests/actions/workflows/go-test.yml/badge.svg)](https://github.com/pkg6/go-requests/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/pkg6/go-requests)](https://goreportcard.com/report/github.com/pkg6/go-requests)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/pkg6/go-requests?tab=doc)
[![Sourcegraph](https://sourcegraph.com/github.com/pkg6/go-requests/-/badge.svg)](https://sourcegraph.com/github.com/pkg6/go-requests?badge)
[![Release](https://img.shields.io/github/release/pkg6/go-requests.svg?style=flat-square)](https://github.com/pkg6/go-requests/releases)

## 基本介绍

`GoRequests`框架提供了强大便捷易用的HTTP客户端，基于`http.Client`进行扩展开发，对象创建可以通过`requests.New()`包方法，也可以通过`new(requests.Client).Clone()`方法调用，同时你还可以使用`requests.NewHttpClient(requests.DefaultHttpClient(nil))`创建对象。推荐使用`requests.New()`来便捷地创建HTTP客户端对象。

## 安装

~~~
$ go get github.com/pkg6/go-requests
~~~


## 基础使用

~~~
package main

import (
	"context"
	"fmt"
	"github.com/pkg6/go-requests"
	"net/url"
)

func main() {
	data := url.Values{}
	data.Set("k", "v")
	get, _ := requests.New().Get(context.Background(), "http://www.httpbin.org/get", data)
	defer get.Close()
	fmt.Println(get.ReadAllString())
	post, _ := requests.New().AsJson().Post(context.Background(), "http://www.httpbin.org/post", data)
	defer post.Close()
	fmt.Println(post.ReadAllString())
}

~~~

## 文件上传

~~~
u := url.Values{}
u.Set("服务端接受的name名称", "@file:loaclfile.txt")
request := requests.New()
request.PostFormWithFiles(context.Background(), "http://127.0.0.1/upload", u)
~~~

## stream请求与返回处理

[ChatGPT【以微软为例】](https://learn.microsoft.com/zh-cn/azure/ai-services/openai/how-to/function-calling)

~~~
package main

import (
	"fmt"
	"github.com/pkg6/go-requests"
)

func main() {
	payload := `{
  "messages": [
    {
      "role": "system",
      "content": "You are an AI assistant that helps people find information."
    },
    {
      "role": "user",
      "content": "你是谁"
    }
  ],
  "model":"gpt-35-turbo",
  "temperature": 0.7,
  "top_p": 0.95,
  "frequency_penalty": 0,
  "presence_penalty": 0,
  "max_tokens": 800,
  "stop": null,
  "stream":true 
}`

	json, _ := requests.PostJson("end-point url", payload, func(client *requests.Client) {
		client.WithHeader("api-key", "api-key")
	})
	defer json.Close()
	json.ReadStream(func(line []byte, number int64) {
		fmt.Println(string(line))
	})
}

~~~

## 方法列表

[https://pkg.go.dev/github.com/pkg6/go-requests]( https://pkg.go.dev/github.com/pkg6/go-requests)

## 返回对象

~~~
(r *Response) Close() error
(r *Response) TraceInfo() TraceInfo
(r *Response) GetCookie() Cookie
(r *Response) ReadAll() []byte
(r *Response) ReadStream(lineNumberFun func(line []byte, number int64)) int64
(r *Response) ReadAllString() string
(r *Response) ContentType() string
(r *Response) Unmarshal(d any) error
(r *Response) IsSuccess() bool
(r *Response) IsError() bool
~~~

## 加入我们

如果你认可我们的开源项目，有兴趣为 [go-requests](https://github.com/pkg6/go-requests) 的发展做贡献，竭诚欢迎加入我们一起开发完善。无论是[报告错误](https://github.com/pkg6/go-requests/issues)或是 [Pull Request](https://github.com/pkg6/go-requests/pulls) 开发，那怕是修改一个错别字也是对我们莫大的帮助。

## License

go-requests is licensed under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0) License - see the [LICENSE](https://github.com/pkg6/go-requests/blob/main/LICENSE) file for details

