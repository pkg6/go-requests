package requests

import (
	"context"
	"crypto/tls"
)

type IRoute interface {
	GetMethod() string
	GetUri() string
	GetHeader() map[string]string
	GetTlsConfig() *tls.Config
	GetCookies() map[string]string
	GetBody() any
	GetD() any
}

type Route struct {
	Method, Uri     string
	Header, Cookies map[string]string
	TlsConfig       *tls.Config
	Body, D         any
}

func (r *Route) GetMethod() string {
	return r.Method
}

func (r *Route) GetUri() string {
	return r.Uri
}

func (r *Route) GetHeader() map[string]string {
	return r.Header
}

func (r *Route) GetTlsConfig() *tls.Config {
	return r.TlsConfig
}

func (r *Route) GetCookies() map[string]string {
	return r.Cookies
}

func (r *Route) GetBody() any {
	return r.Body
}

func (r *Route) GetD() any {
	return r.D
}

// RequestRoute
// route := &requests.Route{
//		Uri:    "https://api.github.com/users/github",
//		Method: http.MethodGet,
//		D:      &.tests.GitHubUser{},
//	}
//	_ = requests.RequestRoute(route)
//	fmt.Println(route.GetD().(*.tests.GitHubUser))
func RequestRoute(route IRoute) error {
	client := new(Client).Clone().SetHttpClient(DefaultHttpClient(nil))
	client.WithHeaderMap(route.GetHeader())
	client.SetTLSConfig(route.GetTlsConfig())
	client.WithCookieMap(route.GetCookies())
	request, err := client.DoRequest(context.Background(), route.GetMethod(), route.GetUri(), route.GetBody())
	if err != nil {
		return err
	}
	return request.Unmarshal(route.GetD())
}
