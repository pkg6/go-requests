package requests

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func requestLogger(client *Client, request *http.Request) error {
	if client.Debug {
		reqLog := "\n==============================================================================\n" +
			"~~~ REQUEST ~~~\n" +
			fmt.Sprintf("%s  %s  %s\n", request.Method, request.URL.RequestURI(), request.Proto) +
			fmt.Sprintf("HOST   : %s\n", request.URL.Host) +
			fmt.Sprintf("HEADERS:\n%s\n", request.Header) +
			fmt.Sprintf("TIME DURATION: %v\n", time.Now()) +
			fmt.Sprintf("BODY   :\n%v\n", request.Body) +
			"------------------------------------------------------------------------------\n"
		client.log.Println(reqLog)
	}
	return nil
}
func responseLogger(client *Client, request *http.Request, response *Response) error {
	if client.Debug {
		var reqBodyContent []byte
		if response.Body != nil {
			reqBodyContent, _ = io.ReadAll(response.Body)
			response.requestBody = reqBodyContent
			response.Body = NewReadCloser(reqBodyContent, false)
		}
		debugLog := "~~~ RESPONSE ~~~\n" +
			fmt.Sprintf("STATUS       : %s\n", response.Status) +
			fmt.Sprintf("PROTO        : %s\n", response.Proto) +
			fmt.Sprintf("ATTEMPT      : %v\n", client.attempt) +
			fmt.Sprintf("RECEIVED AT  : %v\n", time.Now().Format(time.RFC3339Nano)) +
			fmt.Sprintf("TIME DURATION: %v\n", time.Now()) +
			fmt.Sprintf("HEADERS: %v\n", response.Header)
		debugLog += fmt.Sprintf("BODY: %s\n", string(reqBodyContent))
		debugLog += "==============================================================================\n"
		client.log.Println(debugLog)
	}
	return nil
}
