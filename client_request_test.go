package requests

import (
	"context"
	"net/http"
	"net/url"
	"testing"
)

func TestClient_Get(t *testing.T) {
	data := url.Values{}
	data.Set("k", "v")
	body, _ := DefaultClient.
		Get(context.Background(), "http://www.httpbin.org/get", data)
	defer body.Close()
	if body.StatusCode != http.StatusOK {
		t.Fatalf("get error")
	}
}

func TestClient_Post(t *testing.T) {
	data := url.Values{}
	data.Set("k", "v")
	body, _ := DefaultClient.
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
	body, _ := DefaultClient.
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
	body, _ := DefaultClient.
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
	body, _ := DefaultClient.
		AsJson().
		Put(context.Background(), "http://www.httpbin.org/put", data)
	defer body.Close()
	if body.StatusCode != http.StatusOK {
		t.Fatalf("Put error")
	}
}
