package icapclient

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

// Request represents the icap client request data
type Request struct {
	Method                string
	URL                   *url.URL
	Header                http.Header
	HTTPRequest           *http.Request
	HTTPResponse          *http.Response
	ChunkLength           int
	PreviewBytes          int
	previewSet            bool
	bodyFittedInPreview   bool
	remainingPreviewBytes []byte
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

// SetDefaultRequestHeaders assigns some of the headers with its default value if they are not set already
func (r *Request) SetDefaultRequestHeaders() {
	if _, exists := r.Header["Allow"]; !exists {
		r.Header.Add("Allow", "204") // assigning 204 by default if Allow not provided
	}
	if _, exists := r.Header["Host"]; !exists {
		hostName, _ := os.Hostname()
		r.Header.Add("Host", hostName)
	}
}

// DumpRequest returns the given request in its ICAP/1.x wire
// representation.
func DumpRequest(req *Request) ([]byte, error) {

	// Making the ICAP message block

	reqStr := fmt.Sprintf("%s %s %s\n", req.Method, req.URL.RequestURI(), ICAPVersion)

	for headerName, vals := range req.Header {
		for _, val := range vals {
			reqStr += fmt.Sprintf("%s: %s\n", headerName, val)
		}
	}

	reqStr += "Encapsulated: %s\n" // will populate the Encapsulated header value after making the http Request & Response messages
	reqStr += LF

	// Making the HTTP Request message block

	httpReqStr := ""
	if req.HTTPRequest != nil {
		b, err := httputil.DumpRequestOut(req.HTTPRequest, true)

		if err != nil {
			return nil, err
		}

		httpReqStr += string(b)

		if req.previewSet {
			keepPreviewBodyBytes(&httpReqStr, req.PreviewBytes)
		}

		if !bodyAlreadyChunked(httpReqStr) {
			addHexaBodyByteNotations(&httpReqStr)
		}
	}

	// Making the HTTP Response message block

	httpRespStr := ""
	if req.HTTPResponse != nil {
		b, err := httputil.DumpResponse(req.HTTPResponse, true)

		if err != nil {
			return nil, err
		}

		httpRespStr += string(b)

		if req.previewSet {
			keepPreviewBodyBytes(&httpRespStr, req.PreviewBytes)
		}

		if !bodyAlreadyChunked(httpRespStr) {
			addHexaBodyByteNotations(&httpRespStr)
		}
	}

	if httpRespStr != "" && !strings.HasSuffix(httpRespStr, DoubleCRLF) { // if the HTTP Response message block doesn't end with a \r\n\r\n, then going to add one by force for better calculation of byte offsets
		httpRespStr += CRLF
	}

	reqStr = setEncapsulatedHeaderValue(reqStr, httpReqStr, httpRespStr)

	if req.previewSet && req.bodyFittedInPreview {
		httpRespStr = addFullBodyInPreviewIndicator(httpRespStr)
	}

	data := []byte(reqStr + httpReqStr + httpRespStr)

	return data, nil
}
