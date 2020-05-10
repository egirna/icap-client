package icapclient

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Request represents the icap client request data
type Request struct {
	Method       string
	URL          *url.URL
	Header       http.Header
	HTTPRequest  *http.Request
	HTTPResponse *http.Response
}

// NewRequest is the factory function for Request
func NewRequest(method, urlStr string, httpReq *http.Request, httpResp *http.Response) (*Request, error) {

	method = strings.ToUpper(method)

	if valid, err := validMethod(method); !valid {
		return nil, err
	}

	u, err := url.Parse(urlStr)

	if err != nil {
		return nil, err
	}

	if valid, err := validURL(u); !valid {
		return nil, err
	}

	req := &Request{
		Method:       method,
		URL:          u,
		Header:       make(map[string][]string),
		HTTPRequest:  httpReq,
		HTTPResponse: httpResp,
	}

	return req, nil
}

// DumpRequest returns the given request in its ICAP/1.x wire
// representation.
func DumpRequest(req *Request) ([]byte, error) {

	reqStr := fmt.Sprintf("%s %s %s\n", req.Method, req.URL.RequestURI(), ICAPVersion)

	for headerName, vals := range req.Header {
		for _, val := range vals {
			reqStr += fmt.Sprintf("%s: %s\n", headerName, val)
		}
	}

	reqStr += LF

	if req.HTTPRequest != nil {
		b, err := httputil.DumpRequestOut(req.HTTPRequest, true)

		if err != nil {
			return nil, err
		}
		reqStr += string(b)
	}

	if req.HTTPResponse != nil {
		b, err := httputil.DumpResponse(req.HTTPResponse, true)

		if err != nil {
			return nil, err
		}

		reqStr += string(b)
	}

	return []byte(reqStr), nil
}
