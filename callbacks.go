package requests

import (
	"context"
	"fmt"
	"github.com/pkg6/go-requests/jsons"
	"io"
	"net/http"
	"time"
)

const requestLogCtxKey = "request_log_star_time"

func requestLogger(client *Client, request *http.Request) error {
	if client.Debug {
		now := time.Now()
		var body []byte
		if request.Body != nil {
			body, _ = io.ReadAll(request.Body)
			request.Body = NewReadCloser(body, false)
		}
		client.Context = context.WithValue(context.Background(), requestLogCtxKey, now)
		headers, _ := jsons.Marshal(request.Header)
		reqLog := "\n==============================================================================\n" +
			"~~~ REQUEST ~~~\n" +
			fmt.Sprintf("%s  %s  %s\n", request.Method, request.URL.RequestURI(), request.Proto) +
			fmt.Sprintf("HOST           : %s\n", request.URL.Host) +
			fmt.Sprintf("TIME DURATION  : %v\n", now.Format(time.RFC3339Nano)) +
			fmt.Sprintf("HEADERS        : \n%s\n", string(headers)) +
			fmt.Sprintf("BODY           :\n%v\n", string(body)) +
			"------------------------------------------------------------------------------\n"
		client.Logger.Println(reqLog)
	}
	return nil
}
func responseLogger(client *Client, request *http.Request, response *Response) error {
	if client.Debug {
		e := time.Now()
		var reqBodyContent []byte
		if response.Body != nil {
			reqBodyContent, _ = io.ReadAll(response.Body)
			response.requestBody = reqBodyContent
			response.Body = NewReadCloser(reqBodyContent, false)
		}
		s := client.Context.Value(requestLogCtxKey).(time.Time)
		headers, _ := jsons.Marshal(response.Header)
		debugLog := "~~~ RESPONSE ~~~\n" +
			fmt.Sprintf("CLONE        : %v\n", client.clone) +
			fmt.Sprintf("STATUS       : %s\n", response.Status) +
			fmt.Sprintf("PROTO        : %s\n", response.Proto) +
			fmt.Sprintf("ATTEMPT      : %v\n", client.attempt) +
			fmt.Sprintf("RECEIVED AT  : %v\n", e.Format(time.RFC3339Nano)) +
			fmt.Sprintf("HEADERS      : \n%v\n", string(headers))
		debugLog += fmt.Sprintf("BODY : \n%s\n", string(reqBodyContent))
		debugLog += fmt.Sprintf("TIME CONSUMING : %v\n", e.Sub(s))
		debugLog += "==============================================================================\n"
		client.Logger.Println(debugLog)
	}
	return nil
}
