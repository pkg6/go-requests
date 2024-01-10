package requests

import (
	"net/url"
	"reflect"
	"testing"
)

func TestIsJSONType(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect bool
	}{
		{"application/json", true},
		{"application/xml+json", true},
		{"application/vnd.foo+json", true},

		{"application/json; charset=utf-8", true},
		{"application/vnd.foo+json; charset=utf-8", true},

		{"text/json", true},
		{"text/xml+json", true},
		{"text/vnd.foo+json", true},

		{"application/foo-json", false},
		{"application/foo.json", false},
		{"application/vnd.foo-json", false},
		{"application/vnd.foo.json", false},
		{"application/json+xml", false},

		{"text/foo-json", false},
		{"text/foo.json", false},
		{"text/vnd.foo-json", false},
		{"text/vnd.foo.json", false},
		{"text/json+xml", false},
	} {
		result := IsJSONType(test.input)

		if result != test.expect {
			t.Errorf("failed on %q: want %v, got %v", test.input, test.expect, result)
		}
	}
}

func TestIsXMLType(t *testing.T) {
	for _, test := range []struct {
		input  string
		expect bool
	}{
		{"application/xml", true},
		{"application/json+xml", true},
		{"application/vnd.foo+xml", true},

		{"application/xml; charset=utf-8", true},
		{"application/vnd.foo+xml; charset=utf-8", true},

		{"text/xml", true},
		{"text/json+xml", true},
		{"text/vnd.foo+xml", true},

		{"application/foo-xml", false},
		{"application/foo.xml", false},
		{"application/vnd.foo-xml", false},
		{"application/vnd.foo.xml", false},
		{"application/xml+json", false},

		{"text/foo-xml", false},
		{"text/foo.xml", false},
		{"text/vnd.foo-xml", false},
		{"text/vnd.foo.xml", false},
		{"text/xml+json", false},
	} {
		result := IsXMLType(test.input)
		if result != test.expect {
			t.Errorf("failed on %q: want %v, got %v", test.input, test.expect, result)
		}
	}
}

func TestUri(t *testing.T) {
	type args struct {
		uri   string
		query []url.Values
	}
	q := url.Values{}
	q.Set("t", "1")
	parse, _ := url.Parse("http://127.0.0.1?t=1")
	var (
		tests = []struct {
			name string
			args args
			want *url.URL
		}{
			{
				name: "url",
				args: args{
					uri:   "http://127.0.0.1",
					query: []url.Values{q},
				},
				want: parse,
			},
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := URIQuery(tt.args.uri, tt.args.query...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Uri() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyString(t *testing.T) {
	type args struct {
		any any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "map",
			args: args{any: map[string]string{"t": "1"}},
			want: `{"t":"1"}`,
		},
		{
			name: "byte",
			args: args{any: []byte("go-request")},
			want: `go-request`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyString(tt.args.any); got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHttpBuildQuery(t *testing.T) {
	for _, test := range []struct {
		input  map[string]any
		expect string
	}{
		{map[string]any{"a": "1"}, "a=1"},
		//{map[string]any{"a": "1", "b": "2"}, "a=1&b=2"},
	} {
		result := HttpBuildQuery(test.input)
		if result != test.expect {
			t.Errorf("failed on %q: want %v, got %v", test.input, test.expect, result)
		}
	}
}
