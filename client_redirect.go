package requests

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

var (
	ErrAutoRedirectDisabled = errors.New("auto redirect is disabled")
)

type (
	// RedirectPolicy to regulate the redirects in the  client.
	// Objects implementing the RedirectPolicy interface can be registered as
	//
	// Apply function should return nil to continue the redirect jounery, otherwise
	// return error to stop the redirect.
	RedirectPolicy interface {
		Apply(req *http.Request, via []*http.Request) error
	}

	// The RedirectPolicyFunc type is an adapter to allow the use of ordinary functions as RedirectPolicy.
	// If f is a function with the appropriate signature, RedirectPolicyFunc(f) is a RedirectPolicy object that calls f.
	RedirectPolicyFunc func(*http.Request, []*http.Request) error
)

// Apply calls f(req, via).
func (f RedirectPolicyFunc) Apply(req *http.Request, via []*http.Request) error {
	return f(req, via)
}

// NoRedirectPolicy is used to disable redirects in the HTTP client
//
//	requests.WithRedirectPolicy(NoRedirectPolicy())
func NoRedirectPolicy() RedirectPolicy {
	return RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		return ErrAutoRedirectDisabled
	})
}

// FlexibleRedirectPolicy is convenient method to create No of redirect policy for HTTP client.
func FlexibleRedirectPolicy(noOfRedirect int) RedirectPolicy {
	return RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		if len(via) >= noOfRedirect {
			return fmt.Errorf("stopped after %d redirects", noOfRedirect)
		}
		checkHostAndAddHeaders(req, via[0])
		return nil
	})
}

// DomainCheckRedirectPolicy is convenient method to define domain name redirect rule in  client.
// Redirect is allowed for only mentioned host in the policy.
//
//	requests.WithRedirectPolicy(DomainCheckRedirectPolicy("host1.com", "host2.org", "host3.net"))
func DomainCheckRedirectPolicy(hostnames ...string) RedirectPolicy {
	hosts := make(map[string]bool)
	for _, h := range hostnames {
		hosts[strings.ToLower(h)] = true
	}
	fn := RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		if ok := hosts[getHostname(req.URL.Host)]; !ok {
			return errors.New("redirect is not allowed as per DomainCheckRedirectPolicy")
		}
		return nil
	})
	return fn
}

//‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Package Unexported methods
//_______________________________________________________________________

func getHostname(host string) (hostname string) {
	if strings.Index(host, ":") > 0 {
		host, _, _ = net.SplitHostPort(host)
	}
	hostname = strings.ToLower(host)
	return
}

// By default Golang will not redirect request headers
// after go throughing various discussion comments from thread
// https://github.com/golang/go/issues/4800
//  will add all the headers during a redirect for the same host
func checkHostAndAddHeaders(request *http.Request, response *http.Request) {
	requestHostname := getHostname(request.URL.Host)
	responseHostname := getHostname(response.URL.Host)
	if strings.EqualFold(requestHostname, responseHostname) {
		for key, val := range response.Header {
			request.Header[key] = val
		}
	} else { // only library User-Agent header is added
		request.Header.Set(hdrUserAgentKey, defaultClientAgent)
	}
}

// WithRedirectLimit limits the number of jumps.
func (c *Client) WithRedirectLimit(redirectLimit int) *Client {
	c.WithRedirectPolicy(func(req *http.Request, via []*http.Request) error {
		if len(via) >= redirectLimit {
			return http.ErrUseLastResponse
		}
		return nil
	})
	return c
}

// WithRedirectPolicy method sets the client redirect poilicy.  provides ready to use
// redirect policies. Wanna create one for yourself refer to `redirect.go`.
//	WithRedirectLimit(20)
//	WithRedirectPolicy(FlexibleRedirectPolicy(20))
//	WithRedirectPolicy(FlexibleRedirectPolicy(20), DomainCheckRedirectPolicy("host1.com", "host2.net"))
func (c *Client) WithRedirectPolicy(policies ...any) *Client {
	if len(policies) == 1 {
		if checkRedirect, ok := policies[0].(func(req *http.Request, via []*http.Request) error); ok {
			c.SetCheckRedirect(checkRedirect)
			return c
		}
	}
	c.SetCheckRedirect(func(req *http.Request, via []*http.Request) error {
		for _, p := range policies {
			if _, ok := p.(RedirectPolicy); ok {
				if err := p.(RedirectPolicy).Apply(req, via); err != nil {
					return err
				}
			} else {
				c.Logger.Warnf("%v does not implement .RedirectPolicy (missing Apply method)", functionName(p))
			}
		}
		return nil
	})
	return c
}
