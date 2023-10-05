package requests

import (
	"context"
	"encoding/json"
	_testdata "github.com/pkg6/go-requests/.testdata"
	"net/http"
	"net/url"
	"testing"
)

var defaultRequest = New()

func TestClientRequestUnmarshal(t *testing.T) {
	var resp _testdata.GitHubUser
	_ = defaultRequest.DoRequestUnmarshal(context.Background(),
		http.MethodGet,
		"https://api.github.com/users/github", nil, &resp)
	if resp.Id != 9919 {
		t.Fatalf("get error")
	}
}
func TestClientRequestBytes(t *testing.T) {
	var resp _testdata.GitHubUser
	b, _ := defaultRequest.DoRequestBytes(context.Background(),
		http.MethodGet,
		"https://api.github.com/users/github", nil)
	_ = json.Unmarshal(b, &resp)
	if resp.Id != 9919 {
		t.Fatalf("get error")
	}
}

func TestClient_Get(t *testing.T) {
	data := url.Values{}
	data.Set("k", "v")
	body, _ := defaultRequest.
		Get(context.Background(), "http://www.httpbin.org/get", data)
	defer body.Close()
	if body.StatusCode != http.StatusOK {
		t.Fatalf("get error")
	}
}

func TestClient_Post(t *testing.T) {
	data := url.Values{}
	data.Set("k", "v")
	body, _ := defaultRequest.
		AsJson().
		Post(context.Background(), "http://www.httpbin.org/post", data)
	defer body.Close()
	if body.StatusCode != http.StatusOK {
		t.Fatalf("Post error")
	}
}

func TestClient_Delete(t *testing.T) {
	data := url.Values{}
	data.Set("k", "v")
	body, _ := defaultRequest.
		AsJson().
		Delete(context.Background(), "http://www.httpbin.org/delete", data)
	defer body.Close()
	if body.StatusCode != http.StatusOK {
		t.Fatalf("Delete error")
	}
}
func TestClient_Patch(t *testing.T) {
	data := url.Values{}
	data.Set("k", "v")
	body, _ := defaultRequest.
		AsJson().
		Patch(context.Background(), "http://www.httpbin.org/patch", data)
	defer body.Close()
	if body.StatusCode != http.StatusOK {
		t.Fatalf("Patch error")
	}
}
func TestClient_Put(t *testing.T) {
	data := url.Values{}
	data.Set("k", "v")
	body, _ := defaultRequest.
		AsJson().
		Put(context.Background(), "http://www.httpbin.org/put", data)
	defer body.Close()
	if body.StatusCode != http.StatusOK {
		t.Fatalf("Put error")
	}
}
