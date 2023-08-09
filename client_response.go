package requests

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
)

type Response struct {
	*http.Response               // Response is the underlying http.Response object of certain request.
	request        *http.Request // Request is the underlying http.Request object of certain request.
	requestBody    []byte        // The body bytes of certain request, only available in Dump feature.
	client         *Client
	cookies        map[string]string
}

func (r *Response) TraceInfo() TraceInfo {
	ct := r.client.traceContext
	ti := TraceInfo{
		DNSLookup:      ct.dnsDone.Sub(ct.dnsStart),
		TLSHandshake:   ct.tlsHandshakeDone.Sub(ct.tlsHandshakeStart),
		ServerTime:     ct.gotFirstResponseByte.Sub(ct.gotConn),
		IsConnReused:   ct.gotConnInfo.Reused,
		IsConnWasIdle:  ct.gotConnInfo.WasIdle,
		ConnIdleTime:   ct.gotConnInfo.IdleTime,
		RequestAttempt: r.client.attempt,
	}
	// Calculate the total time accordingly,
	// when connection is reused
	if ct.gotConnInfo.Reused {
		ti.TotalTime = ct.endTime.Sub(ct.getConn)
	} else {
		ti.TotalTime = ct.endTime.Sub(ct.dnsStart)
	}
	// Only calculate on successful connections
	if !ct.connectDone.IsZero() {
		ti.TCPConnTime = ct.connectDone.Sub(ct.dnsDone)
	}
	// Only calculate on successful connections
	if !ct.gotConn.IsZero() {
		ti.ConnTime = ct.gotConn.Sub(ct.getConn)
	}
	// Only calculate on successful connections
	if !ct.gotFirstResponseByte.IsZero() {
		ti.ResponseTime = ct.endTime.Sub(ct.gotFirstResponseByte)
	}
	// Capture remote address info when connection is non-nil
	if ct.gotConnInfo.Conn != nil {
		ti.RemoteAddr = ct.gotConnInfo.Conn.RemoteAddr()
	}
	return ti
}

func (r *Response) GetCookieMap() map[string]string {
	if r.cookies == nil {
		r.cookies = make(map[string]string)
		for _, c := range r.Cookies() {
			r.cookies[c.Name] = c.Value
		}
	}
	return r.cookies
}

func (r *Response) GetCookie(key string) string {
	r.GetCookieMap()
	if cookie, ok := r.cookies[key]; ok {
		return cookie
	}
	return ""
}

func (r *Response) GetCookieString() string {
	return MapCookiesToString(r.GetCookieMap(), "")
}

// ReadAll retrieves and returns the response content as []byte.
func (r *Response) ReadAll() []byte {
	// Response might be nil.
	if r == nil || r.Response == nil {
		return []byte{}
	}
	body, err := io.ReadAll(r.Response.Body)
	if err != nil {
		return nil
	}
	return body
}

// ReadStreamArgs
//openai chatgpt: requests.ReadStreamArgs{
//		DataPrefix:  []byte("data: "),
//		ErrorPrefix: []byte(`data: {"error":`),
//		EndPrefix:   []byte("[DONE]"),
//		Callback: func(line []byte) {
//			fmt.Println(string(line))
//		},
//	}
type ReadStreamArgs struct {
	DataPrefix  []byte
	ErrorPrefix []byte
	EndPrefix   []byte
	Callback    func(line []byte)
}

func (r *Response) ReadStream(readStreamArgs ReadStreamArgs) uint {
	var (
		emptyMessagesCount uint
		hasErrorPrefix     bool
		emptyArgs          bool
	)
	if readStreamArgs.DataPrefix == nil ||
		readStreamArgs.ErrorPrefix == nil ||
		readStreamArgs.EndPrefix == nil ||
		readStreamArgs.Callback == nil {
		emptyArgs = true
	}
	if emptyArgs || r.IsError() {
		return emptyMessagesCount
	}
	for {
		rawLine, err := bufio.NewReader(r.Response.Body).ReadBytes('\n')
		if err != nil {
			break
		}
		spaceLine := bytes.TrimSpace(rawLine)
		if bytes.HasPrefix(spaceLine, readStreamArgs.ErrorPrefix) {
			hasErrorPrefix = true
		}
		if !bytes.HasPrefix(spaceLine, readStreamArgs.DataPrefix) || hasErrorPrefix {
			if hasErrorPrefix {
				spaceLine = bytes.TrimPrefix(spaceLine, readStreamArgs.DataPrefix)
			}
			emptyMessagesCount++
			continue
		}
		noPrefixLine := bytes.TrimPrefix(spaceLine, readStreamArgs.DataPrefix)
		if string(noPrefixLine) == string(readStreamArgs.EndPrefix) {
			break
		}
		readStreamArgs.Callback(noPrefixLine)
	}
	return emptyMessagesCount
}

// ReadAllString retrieves and returns the response content as string.
func (r *Response) ReadAllString() string {
	return string(r.ReadAll())
}

// ContentType response header Content-Type
func (r *Response) ContentType() string {
	return r.Response.Header.Get(HttpHeaderContentType)
}

// Unmarshal content into object from JSON or XML
func (r *Response) Unmarshal(d any) error {
	return r.client.Unmarshal(r.ContentType(), r.ReadAll(), d)
}

// IsSuccess method returns true if HTTP status `code >= 200 and <= 299` otherwise false.
func (r *Response) IsSuccess() bool {
	return r.StatusCode > 199 && r.StatusCode < 300
}

// IsError method returns true if HTTP status `code >= 400` otherwise false.
func (r *Response) IsError() bool {
	return r.StatusCode > 399
}

// Close closes the response when it will never be used.
func (r *Response) Close() error {
	if r == nil || r.Response == nil {
		return nil
	}
	return r.Response.Body.Close()
}
