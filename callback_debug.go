package requests

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	ctxDebugStartTime CtxKeyString = "_request_debug_logger_start_time"
)

func onAfterRequestByDebug(client *Client, request *http.Request) error {
	if client.Debug {
		now := time.Now()
		var body []byte
		if request.Body != nil {
			body, _ = io.ReadAll(request.Body)
			request.Body = NewReadCloser(body, false)
		}
		client.ctx = context.WithValue(context.Background(), ctxDebugStartTime, now)
		headers, _ := client.JSONMarshal(request.Header)
		reqLog := "\n==============================================================================\n" +
			"~~~ REQUEST ~~~\n" +
			fmt.Sprintf("%s  %s  %s\n", request.Method, request.URL.RequestURI(), request.Proto) +
			fmt.Sprintf("HOST           : %s\n", request.URL.Host) +
			fmt.Sprintf("TIME DURATION  : %v\n", now.Format(time.RFC3339Nano)) +
			fmt.Sprintf("HEADERS        : \n%s\n", string(headers)) +
			fmt.Sprintf("REQUEST BODY           :\n%v\n", string(body)) +
			"------------------------------------------------------------------------------\n"
		client.Logger.Debugf(reqLog)
	}
	return nil
}
func onResponseByDebug(client *Client, request *http.Request, response *Response) error {
	if client.Debug {
		e := time.Now()
		var responseBody []byte
		if response.Body != nil {
			responseBody, _ = io.ReadAll(response.Body)
			response.Body = NewReadCloser(responseBody, false)
		}
		s := client.ctx.Value(ctxDebugStartTime).(time.Time)
		headers, _ := client.JSONMarshal(response.Header)
		debugLog := "~~~ RESPONSE ~~~\n" +
			fmt.Sprintf("CLONE        : %v\n", client.clone) +
			fmt.Sprintf("STATUS       : %s\n", response.Status) +
			fmt.Sprintf("PROTO        : %s\n", response.Proto) +
			fmt.Sprintf("ATTEMPT      : %v\n", client.attempt) +
			fmt.Sprintf("RECEIVED AT  : %v\n", e.Format(time.RFC3339Nano)) +
			fmt.Sprintf("HEADERS      : \n%v\n", string(headers))
		debugLog += fmt.Sprintf("RESPONSE BODY : \n%s\n", string(responseBody))
		debugLog += fmt.Sprintf("TIME CONSUMING : %v\n", e.Sub(s))
		debugLog += "==============================================================================\n"
		client.Logger.Debugf(debugLog)
		client.ctx = context.Background()
	}
	return nil
}

func onResponseByDebugWriter(client *Client, request *http.Request, response *Response) error {
	if client.writer != nil {
		var builder strings.Builder
		builder.WriteString("REQUEST: \n")
		builder.WriteString(fmt.Sprintf("%s %s %s \n", request.Method, request.URL.String(), request.Proto))
		builder.WriteString(fmt.Sprintf("Clone: %d \n", client.clone))
		builder.WriteString(fmt.Sprintf("Attempt: %d \n", client.attempt))
		reqHeader := request.Header
		for s := range reqHeader {
			builder.WriteString(fmt.Sprintf("%s : %s \n", s, reqHeader.Get(s)))
		}
		builder.WriteString("\n")
		builder.WriteString("RESPONSE: \n")
		builder.WriteString(fmt.Sprintf("%s %s \n", response.Proto, response.Status))
		respHeader := response.Header
		for s := range respHeader {
			builder.WriteString(fmt.Sprintf("%s : %s \n", s, respHeader.Get(s)))
		}
		builder.WriteString("\n")
		var responseBody []byte
		if response.Body != nil {
			responseBody, _ = io.ReadAll(response.Body)
			response.Body = NewReadCloser(responseBody, false)
		}
		builder.Write(responseBody)
		_, err := fmt.Fprintln(client.writer, builder.String())
		return err
	}
	return nil
}
