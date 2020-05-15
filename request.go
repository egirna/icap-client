package icapclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
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

// SetPreview sets the preview bytes in the icap header
func (r *Request) SetPreview(maxBytes int) error {
	bodyBytes, err := ioutil.ReadAll(r.HTTPResponse.Body)

	if err != nil {
		return err
	}

	previewBytes := len(bodyBytes)

	if len(bodyBytes) > maxBytes {
		previewBytes = maxBytes
	}

	r.Header.Set("Preview", strconv.Itoa(previewBytes))

	r.HTTPResponse.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))

	return nil

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

	reqStr += "Encapsulated: %s\n"
	reqStr += LF

	httpReqStr := ""
	if req.HTTPRequest != nil {
		b, err := httputil.DumpRequestOut(req.HTTPRequest, true)

		if err != nil {
			return nil, err
		}
		httpReqStr += string(b)
	}

	httpRespStr := ""
	if req.HTTPResponse != nil {
		b, err := httputil.DumpResponse(req.HTTPResponse, true)

		if err != nil {
			return nil, err
		}

		httpRespStr += string(b)
	}

	if httpRespStr != "" && !strings.HasSuffix(httpRespStr, DoubleCRLF) {
		httpRespStr += DoubleCRLF
	}

	reqStr = setEncapsulatedHeaderValue(reqStr, httpReqStr, httpRespStr)

	data := []byte(reqStr + httpReqStr + httpRespStr)

	fmt.Println(string(data))

	return data, nil
}

// SetDefaultRequestHeaders assigns some of the headers with its default value if they are not set already
func SetDefaultRequestHeaders(req *Request) {
	if _, exists := req.Header["Allow"]; !exists {
		req.Header.Add("Allow", "204") // assigning 204 by default if Allow not provided
	}
	if _, exists := req.Header["Host"]; !exists {
		hostName, _ := os.Hostname()
		req.Header.Add("Host", hostName)
	}
}
