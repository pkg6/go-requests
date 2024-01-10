package requests

import (
	"context"
	"crypto/tls"
	"net"
	"net/http/httptrace"
	"time"
)

func (c *Client) withContext(ctx context.Context) context.Context {
	if c.trace {
		if ctx == nil {
			ctx = c.traceContext.createContext(context.Background())
		} else {
			ctx = c.traceContext.createContext(ctx)
		}
	}
	return ctx
}
func (c *Client) EnableTrace() *Client {
	c.trace = true
	return c
}

func (c *Client) getTraceInfo() TraceInfo {
	return c.traceContext.Get(c)
}

type TraceInfo struct {
	// DNSLookup is a duration that transport took to perform
	// DNS lookup.
	DNSLookup time.Duration
	// ConnTime is a duration that took to obtain a successful connection.
	ConnTime time.Duration
	// TCPConnTime is a duration that took to obtain the TCP connection.
	TCPConnTime time.Duration
	// TLSHandshake is a duration that TLS handshake took place.
	TLSHandshake time.Duration
	// ServerTime is a duration that server took to respond first byte.
	ServerTime time.Duration
	// ResponseTime is a duration since first response byte from server to
	// request completion.
	ResponseTime time.Duration
	// TotalTime is a duration that total request took end-to-end.
	TotalTime time.Duration
	// IsConnReused is whether this connection has been previously
	// used for another HTTP request.
	IsConnReused bool
	// IsConnWasIdle is whether this connection was obtained from an
	// idle pool.
	IsConnWasIdle bool
	// ConnIdleTime is a duration how long the connection was previously
	// idle, if IsConnWasIdle is true.
	ConnIdleTime time.Duration
	// RequestAttempt is to represent the request attempt made during a
	// request execution flow, including retry count.
	RequestAttempt int
	// RemoteAddr returns the remote network address.
	RemoteAddr net.Addr
}
type traceContext struct {
	getConn              time.Time
	dnsStart             time.Time
	dnsDone              time.Time
	connectDone          time.Time
	tlsHandshakeStart    time.Time
	tlsHandshakeDone     time.Time
	gotConn              time.Time
	gotFirstResponseByte time.Time
	endTime              time.Time
	gotConnInfo          httptrace.GotConnInfo
}

func (t *traceContext) createContext(ctx context.Context) context.Context {
	return httptrace.WithClientTrace(
		ctx,
		&httptrace.ClientTrace{
			DNSStart: func(_ httptrace.DNSStartInfo) {
				t.dnsStart = time.Now()
			},
			DNSDone: func(_ httptrace.DNSDoneInfo) {
				t.dnsDone = time.Now()
			},
			ConnectStart: func(_, _ string) {
				if t.dnsDone.IsZero() {
					t.dnsDone = time.Now()
				}
				if t.dnsStart.IsZero() {
					t.dnsStart = t.dnsDone
				}
			},
			ConnectDone: func(net, addr string, err error) {
				t.connectDone = time.Now()
			},
			GetConn: func(_ string) {
				t.getConn = time.Now()
			},
			GotConn: func(ci httptrace.GotConnInfo) {
				t.gotConn = time.Now()
				t.gotConnInfo = ci
			},
			GotFirstResponseByte: func() {
				t.gotFirstResponseByte = time.Now()
			},
			TLSHandshakeStart: func() {
				t.tlsHandshakeStart = time.Now()
			},
			TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
				t.tlsHandshakeDone = time.Now()
			},
		},
	)
}

func (t *traceContext) Get(client *Client) TraceInfo {
	ct := t
	ti := TraceInfo{
		DNSLookup:      ct.dnsDone.Sub(ct.dnsStart),
		TLSHandshake:   ct.tlsHandshakeDone.Sub(ct.tlsHandshakeStart),
		ServerTime:     ct.gotFirstResponseByte.Sub(ct.gotConn),
		IsConnReused:   ct.gotConnInfo.Reused,
		IsConnWasIdle:  ct.gotConnInfo.WasIdle,
		ConnIdleTime:   ct.gotConnInfo.IdleTime,
		RequestAttempt: client.attempt,
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
