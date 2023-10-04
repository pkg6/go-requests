package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// AnyString converts `any` to string.
// It's most commonly used converting function.
func AnyString(any any) string {
	if any == nil {
		return ""
	}
	switch value := any.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	case time.Time:
		if value.IsZero() {
			return ""
		}
		return value.String()
	case *time.Time:
		if value == nil {
			return ""
		}
		return value.String()
	case time.Duration:
		return value.String()
	default:
		// Empty checks.
		if value == nil {
			return ""
		}
		if f, ok := value.(error); ok {
			// If the variable implements the Error() interface,
			// then use that interface to perform the conversion
			return f.Error()
		}
		// Reflect checks.
		var (
			rv   = reflect.ValueOf(value)
			kind = rv.Kind()
		)
		switch kind {
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return ""
			}
		case reflect.String:
			return rv.String()
		}
		if kind == reflect.Ptr {
			return AnyString(rv.Elem().Interface())
		}
		// Finally, we use json.Marshal to convert.
		if jsonContent, err := json.Marshal(value); err != nil {
			return fmt.Sprint(value)
		} else {
			return string(jsonContent)
		}
	}
}

func IsMatchString(expr string, s string) bool {
	return regexp.MustCompile(expr).MatchString(s)
}

// IsJSONType method is to check JSON content type or not
func IsJSONType(s string) bool {
	return IsMatchString(`(?i:(application|text)/(json|.*\+json|json\-.*)(;|$))`, s)
}

// IsXMLType method is to check XML content type or not
func IsXMLType(s string) bool {
	return IsMatchString(`(?i:(application|text)/(xml|.*\+xml)(;|$))`, s)
}

func IsStreamType(s string) bool {
	return strings.Contains(s, "event-stream")
}

func functionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func URIQuery(uri string, query ...url.Values) *url.URL {
	u, _ := url.Parse(uri)
	q := u.Query()
	for _, value := range query {
		for k := range value {
			q.Set(k, value.Get(k))
		}
	}
	u.RawQuery = q.Encode()
	return u
}

func UrlValues(uvs ...url.Values) url.Values {
	query := url.Values{}
	for _, value := range uvs {
		for k := range value {
			query.Set(k, value.Get(k))
		}
	}
	return query
}

// HttpBuildQuery Generate get request parameters
func HttpBuildQuery(data map[string]any) string {
	var buf bytes.Buffer
	for k, v := range data {
		buildQuery(&buf, k, v)
	}
	return buf.String()
}

func buildQuery(buf *bytes.Buffer, key string, val any) {
	if buf.Len() > 0 {
		buf.WriteByte('&')
	}
	buf.WriteString(url.QueryEscape(key))
	buf.WriteByte('=')
	switch v := val.(type) {
	case map[string]any:
		for sk, sv := range v {
			buildQuery(buf, fmt.Sprintf("%s[%s]", key, sk), sv)
		}
	case []any:
		for _, item := range v {
			buildQuery(buf, key, item)
		}
	default:
		buf.WriteString(url.QueryEscape(fmt.Sprint(v)))
	}
}
