package icapclient

import (
	"net/http"
	"net/url"
)

// Request represents the icap client request data
type Request struct {
	Method     string
	URL        *url.URL
	HTTRequest *http.Request
}

// NewRequest is the factory function for Request
func NewRequest(method, urlStr string, httpReq *http.Request) (*Request, error) {

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
		Method:     method,
		URL:        u,
		HTTRequest: httpReq,
	}

	return req, nil
}
