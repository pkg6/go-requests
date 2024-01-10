package requests

import (
	"net/http"
)

// OnBeforeRequest method appends request callback into the before request chain.
//	client.OnBeforeRequest(func(c *requests.Client) error {
//			// Now you have access to Client and Request instance
//			// manipulate it as per your need
//			return nil 	// if its success otherwise return error
//		})
func (c *Client) OnBeforeRequest(callback ClientCallback) *Client {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.beforeRequestCallbacks = append(c.beforeRequestCallbacks, callback)
	return c
}
func (c *Client) doBeforeRequestCallbacks() error {
	for _, fn := range c.beforeRequestCallbacks {
		if err := fn(c); err != nil {
			return err
		}
	}
	return nil
}

// OnAfterRequest method appends request callback into the before request chain.
//	client.OnAfterRequest(func(client *requests.Client, request *http.Request) error{
//			// Now you have access to Client and Request instance
//			// manipulate it as per your need
//
//			return nil 	// if its success otherwise return error
//		})
func (c *Client) OnAfterRequest(callback RequestCallback) *Client {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.afterRequestCallbacks = append(c.afterRequestCallbacks, callback)
	return c
}
func (c *Client) doAfterRequestCallbacks(request *http.Request) error {
	for _, fn := range c.afterRequestCallbacks {
		if err := fn(c, request); err != nil {
			return err
		}
	}
	return nil
}

// OnResponse method appends response callback into the after response chain.
//	client.OnResponse(func(request *http.Request, response *requests.Response) error {
//			// Now you have access to Client and Response instance
//			// manipulate it as per your need
//
//			return nil 	// if its success otherwise return error
//		})
func (c *Client) OnResponse(callback ResponseCallback) *Client {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.responseCallbacks = append(c.responseCallbacks, callback)
	return c
}
func (c *Client) doResponseCallbacks(request *http.Request, response *Response) error {
	for _, fn := range c.responseCallbacks {
		if err := fn(c, request, response); err != nil {
			return err
		}
	}
	return nil
}

// OnSuccess method adds a callback that will be run whenever a request execution
// succeeds.  This is called after all retries have been attempted (if any).
//
// Out of the OnSuccess, OnError, OnInvalid, OnPanic callbacks, exactly one
// set will be invoked for each call to Request.Execute() that comletes.
func (c *Client) OnSuccess(h SuccessHook) *Client {
	c.successHooks = append(c.successHooks, h)
	return c
}
func (c *Client) doSuccessHooks(resp *Response) {
	for _, h := range c.successHooks {
		h(c, resp)
	}
}

// OnError method adds a callback that will be run whenever a request execution fails.
// This is called after all retries have been attempted (if any).
// If there was a response from the server, the error will be wrapped in *ResponseError
// which has the last response received from the server.
//
//	client.OnError(func(*http.Request, error){
//		if v, ok := err.(*requests.ResponseError); ok {
//			// Do something with v.Response
//		}
//		// Log the error, increment a metric, etc...
//	})
//
// Out of the OnSuccess, OnError, OnInvalid, OnPanic callbacks, exactly one
// set will be invoked for each call to Request.Execute() that comletes.
func (c *Client) OnError(h ErrorHook) *Client {
	c.errorHooks = append(c.errorHooks, h)
	return c
}

// Helper to run errorHooks hooks.
// It wraps the error in a ResponseError if the response is not nil
// so hooks can access it.
func (c *Client) doErrorHooks(request *http.Request, response *Response, err error) {
	if err != nil {
		if response == nil {
			err = &ResponseError{Response: response, Err: err}
		}
		for _, h := range c.errorHooks {
			h(c, request, err)
		}
	} else {
		for _, h := range c.successHooks {
			h(c, response)
		}
	}
}

// OnPanic method adds a callback that will be run whever a request execution
// panics.
//
// Out of the OnSuccess, OnError, OnInvalid, OnPanic callbacks, exactly one
// set will be invoked for each call to Request.Execute() that completes.
// If an OnSuccess, OnError, or OnInvalid callback panics, then the exactly
// one rule can be violated.
func (c *Client) OnPanic(h ErrorHook) *Client {
	c.panicHooks = append(c.panicHooks, h)
	return c
}

// Helper to run panicHooks hooks.
func (c *Client) doPanicHooks(request *http.Request, err error) {
	for _, h := range c.panicHooks {
		h(c, request, err)
	}
}
