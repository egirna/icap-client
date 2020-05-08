package icapclient

import (
	"bufio"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Response represents the icap server response data
type Response struct {
	StatusCode      int
	Status          string
	Header          http.Header
	ContentRequest  *http.Request
	ContentResponse *http.Response
}

// ReadResponse converts a Reader to a icapclient Response
func ReadResponse(b *bufio.Reader) (*Response, error) {

	resp := &Response{
		Header: make(map[string][]string),
	}

	icapRespStr, err := b.ReadString('\n') //  reading the first line which is mandatory for any ICAP server to provide

	if err != nil {
		return nil, err
	}

	ss := strings.Split(icapRespStr, " ") // splitting the first line to analyze each component

	if len(ss) < 3 { // should have at least three, for example: ICAP/1.0 200 OK
		return nil, errors.New(ErrInvalidICAPResponse)
	}

	if ss[0] != ICAPVersion {
		return nil, errors.New(ErrInvalidICAPResponse)
	}

	statusCode, err := strconv.Atoi(ss[1])

	if err != nil {
		return nil, err
	}

	resp.StatusCode = statusCode
	resp.Status = strings.TrimSpace(ss[2])

	for err == nil && (icapRespStr != "\n" && err != io.EOF) { // reading the rest of the headers for icap until there is \n or EOF which indicates the end of ICAP portion of the message

		icapRespStr, err = b.ReadString('\n')

		headerVal := strings.Split(icapRespStr, ":")

		if len(headerVal) > 2 {
			resp.Header.Add(strings.TrimSpace(headerVal[0]),
				strings.TrimSpace(strings.TrimSpace(strings.Join(headerVal[1:], ""))))
		}

	}

	if err == io.EOF { // there is nothing but the ICAP portion in the message
		return resp, nil
	}

	httpStr, err := b.ReadString('\n') // reading the first line of the HTTP portion of the message

	if err != nil {
		return nil, err
	}

	httpStr = strings.TrimSpace(httpStr)

	ss = strings.Split(httpStr, " ")

	if len(ss) < 3 {
		return nil, errors.New(ErrInvalidHTTPResponse)
	}

	httpScheme := ss[0]

	httpMsg := httpStr + "\n" // the "\n" needs to be added here explicitly because trimming the string removes the "\n"s

	for err == nil && (httpStr != "\n" || err != io.EOF) {
		httpStr, err = b.ReadString('\n')

		httpMsg += strings.TrimSpace(httpStr) + "\n"
	}

	if httpScheme != HTTPVersion { // this expects the first line and first word of the message to be any of the HTTP Methods, which in turn indicates its an HTTP Request
		resp.ContentRequest, err = http.ReadRequest(bufio.NewReader(strings.NewReader(httpMsg)))

		if err != nil {
			return nil, err
		}
	}

	if httpScheme == HTTPVersion { // this indicates its a HTTP Response
		resp.ContentResponse, err = http.ReadResponse(bufio.NewReader(strings.NewReader(httpMsg)), resp.ContentRequest)

		if err != nil {
			return nil, err
		}
	}

	return resp, nil

}
