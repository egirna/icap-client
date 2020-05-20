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
		httpReq, err := addHexaRequestBodyByteNotations(*req.HTTPRequest)
		if err != nil {
			return nil, err
		}
		b, err := httputil.DumpRequestOut(&httpReq, true)

		if err != nil {
			return nil, err
		}
		httpReqStr += string(b)
	}

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

		httpRespStr += string(b)
	}

	if httpRespStr != "" && !strings.HasSuffix(httpRespStr, DoubleCRLF) {
		httpRespStr += DoubleCRLF
	}

	reqStr = setEncapsulatedHeaderValue(reqStr, httpReqStr, httpRespStr)

	data := []byte(reqStr + httpReqStr + httpRespStr)

	return data, nil
}

// addHexaResponseBodyByteNotations adds body bytes Hexadecimal notations before and after body chunk
// for example: for a body, "Hello World", this function adds
// b
// Hello World
// 0
func addHexaResponseBodyByteNotations(r http.Response) (http.Response, error) {

	if r.Body == nil {
		return r, nil
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return r, err
	}

	defer r.Body.Close()

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

	defer req.Body.Close()

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
