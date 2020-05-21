package icapclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
		httpReq, err := addHexaRequestBodyByteNotations(*req.HTTPRequest) // getting a copy of the request as don't want to tamper with the original resource
		if err != nil {
			return nil, err
		}
		b, err := httputil.DumpRequestOut(&httpReq, true)

		if err != nil {
			return nil, err
		}

		if httpReq.Body != nil {
			defer httpReq.Body.Close()
		}

		httpReqStr += string(b)
	}

	// Making the HTTP Response message block

	httpRespStr := ""
	if req.HTTPResponse != nil {
		httpResp, err := addHexaResponseBodyByteNotations(*req.HTTPResponse)
		if err != nil {
			return nil, err
		}
		b, err := httputil.DumpResponse(&httpResp, true)

		if err != nil {
			return nil, err
		}

		if httpResp.Body != nil {
			defer httpResp.Body.Close()
		}

		httpRespStr += string(b)
	}

	if httpRespStr != "" && !strings.HasSuffix(httpRespStr, DoubleCRLF) { // if the HTTP Response message block doesn't end with a \r\n\r\n, then going to add one by force for better calculation of byte offsets
		httpRespStr += CRLF
	}

	reqStr = setEncapsulatedHeaderValue(reqStr, httpReqStr, httpRespStr)

	if req.bodyFittedInPreview {
		httpRespStr = addFullBodyInPreviewIndicator(httpRespStr)
	}

	data := []byte(reqStr + httpReqStr + httpRespStr)

	fmt.Println(string(data))

	return data, nil
}

// addHexaResponseBodyByteNotations adds body bytes Hexadecimal notations before and after body chunk
// for example: for a body, "Hello World", this function adds
// b
// Hello World
// 0
func addHexaResponseBodyByteNotations(r http.Response) (http.Response, error) {

	// FIXME: stop modifying the original response

	if r.Body == nil {
		return r, nil
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return r, err
	}

	if len(b) < 1 {
		return r, nil
	}

	bodyStr := string(b)

	bodyStr = fmt.Sprintf("%x\r\n", len(b)) + bodyStr + CRLF + "0" + CRLF // the byte length in Hexadecimal with a body and ending with \r\n0\r\n

	newBodyByte := []byte(bodyStr)

	r.Body = ioutil.NopCloser(bytes.NewBuffer(newBodyByte)) // returning the body bytes back to the body as it is already read once
	r.ContentLength = int64(len(newBodyByte))               // adapting the content length according to the new byte length of the body

	return r, nil

}

// addHexaRequestBodyByteNotations adds body bytes Hexadecimal notations before and after body chunk
// for example: for a body, "Hello World", this function adds
// b
// Hello World
// 0
func addHexaRequestBodyByteNotations(req http.Request) (http.Request, error) {
	if req.Body == nil {
		return req, nil
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return req, err
	}

	if len(b) < 1 {
		return req, nil
	}

	bodyStr := string(b)

	bodyStr = fmt.Sprintf("%x\r\n", len(b)) + bodyStr + CRLF + "0" + CRLF // the byte length in Hexadecimal with a body and ending with \r\n0\r\n

	newBodyByte := []byte(bodyStr)

	req.Body = ioutil.NopCloser(bytes.NewBuffer(newBodyByte)) // returning the body bytes back to the body as it is already read once
	req.ContentLength = int64(len(newBodyByte))               // adapting the content length according to the new byte length of the body

	return req, nil

}
