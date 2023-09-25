## openssl版本

~~~
openssl version
OpenSSL 1.1.1u  30 May 2023
~~~

## 自认证根证书

~~~
创建根证书密钥文件(自己做CA)root.key(需要输入两次私钥密码)
openssl genrsa -des3 -out root.key 2048
创建根证书的申请文件root.csr(输入root.key的密码)
openssl req -new -key root.key -out root.csr
创建根证书root.crt：
openssl x509 -req -days 3650 -sha256 -extensions v3_ca -signkey root.key -in root.csr -out root.crt
~~~

## SSL单向认证

~~~
创建服务器证书秘钥(需要输入两次私钥密码)
openssl genrsa –des3 -out server.key 2048
去除key口令(需要输入私钥密码)
openssl rsa -in server.key -out server.key
创建服务器证书申请文件server.csr
openssl req -new -key server.key -out server.csr
创建服务器证书server.crt
openssl x509 -req -days 365 -sha256 -extensions v3_req -CA root.crt -CAkey root.key -CAcreateserial -in server.csr -out server.crt
~~~

## 使用go启动一个单向https服务

~~~
package main

import (
	"io"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "golang https server")
}

func main() {
	http.HandleFunc("/", handler)
	if err := http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
~~~

## SSL双向认证

~~~
创建客户端证书私钥
openssl genrsa -des3 -out client.key 2048
去除key口令
openssl rsa -in client.key -out client.key
创建客户端证书申请文件client.csr
openssl req -new -key client.key -out client.csr
创建客户端证书文件client.crt
openssl x509 -req -days 365 -sha256 -extensions v3_req -CA root.crt -CAkey root.key -CAcreateserial -in client.csr -out client.crt
将客户端证书文件client.crt和客户端证书密钥文件client.key合并成客户端证书安装包client.pfx
openssl pkcs12 -export -in client.crt -inkey client.key -out client.pfx
~~~

## 使用go启动一个双向https服务

~~~
package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type httpsHandler struct {
}

func (*httpsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "golang https server!!!")
}

func main() {
	pool := x509.NewCertPool()
	caCertPath := "root.crt"
	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		log.Fatal("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)
	s := &http.Server{
		Addr:    ":8080",
		Handler: &httpsHandler{},
		TLSConfig: &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}
	if err = s.ListenAndServeTLS("server.crt", "server.key"); err != nil {
		log.Fatal("ListenAndServeTLS err:", err)
	}
}
~~~

## 使用go-requests访问https服务器

~~~
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/pkg6/go-requests"
	"log"
	"os"
)

func main() {
	client := requests.DefaultClient
	pool := x509.NewCertPool()
	caCertPath := "root.crt"
	caCrt, err := os.ReadFile(caCertPath)
	if err != nil {
		log.Fatal("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)
	cliCrt, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		log.Fatal("LoadX509KeyPair err:", err)
		return
	}
	client.SetTLSConfig(&tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{cliCrt},
	})
	client.Get(context.Background(), "https://localhost:8080", nil)
}
~~~

