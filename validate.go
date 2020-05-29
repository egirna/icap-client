package icapclient

import (
	"errors"
	"net/http"
	"net/url"
)

func validMethod(method string) (bool, error) {
	if _, registered := registeredMethods[method]; !registered {
		return false, errors.New(ErrMethodNotRegistered)
	}

	return true, nil
}

func validURL(url *url.URL) (bool, error) {

	if url.Scheme != SchemeICAP {
		return false, errors.New(ErrInvalidScheme)
	}

	if url.Host == "" {
		return false, errors.New(ErrInvalidHost)
	}

	return true, nil
}

func validMethodWithHTTP(httpReq *http.Request, httpResp *http.Response, method string) (bool, error) {
	if method == MethodREQMOD && httpReq == nil {
		return false, errors.New(ErrREQMODWithNoReq)
	}
	if method == MethodREQMOD && httpResp != nil {
		return false, errors.New(ErrREQMODWithResp)
	}
	if method == MethodRESPMOD && httpResp == nil {
		return false, errors.New(ErrRESPMODWithNoResp)
	}

	return true, nil
}

// Validate validates the ICAP request
func (r *Request) Validate() error {

	if valid, err := validMethod(r.Method); !valid {
		return err
	}

	if valid, err := validURL(r.URL); !valid {
		return err
	}

	if valid, err := validMethodWithHTTP(r.HTTPRequest, r.HTTPResponse, r.Method); !valid {
		return err
	}

	return nil
}
