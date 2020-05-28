package icapclient

import (
	"fmt"
	"net/http"
	"strconv"
)

// Client represents the icap client who makes the icap server calls
type Client struct {
	scktDriver *Driver
}

// Do makes the call
func (c *Client) Do(req *Request) (*Response, error) {

	port, err := strconv.Atoi(req.URL.Port())

	if err != nil {
		return nil, err
	}

	c.scktDriver = NewDriver(req.URL.Hostname(), port)

	if err := c.scktDriver.Connect(); err != nil {
		return nil, err
	}

	defer c.scktDriver.Close()

	req.SetDefaultRequestHeaders()

	d := []byte{}
	if req.Method == MethodOPTIONS {
		var err error
		d, err = DumpRequest(req)

		if err != nil {
			return nil, err
		}
	} else {
		d = []byte("RESPMOD /respmod-icapeg ICAP/1.0\nPreview: 64\nAllow: 204\nHost: Anondos-MacBook-Pro.local\nEncapsulated:  req-hdr=0, res-hdr=112, res-body=339\n\nGET /download/eicar.com HTTP/1.1\r\nHost: www.eicar.org\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\nHTTP/1.1 200 OK\r\nContent-Length: 68\r\nCache-Control: private\r\nContent-Disposition: attachment; filename=\"eicar.com\"\r\nContent-Type: application/octet-stream\r\nDate: Thu, 28 May 2020 13:01:55 GMT\r\nServer: Apache/2.4.10 (Debian)\r\n\r\n44\r\nX5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*\r\n0; ieof\r\n\r\n")
	}

	fmt.Println(string(d))

	if err := c.scktDriver.Send(d); err != nil {
		return nil, err
	}

	resp, err := c.scktDriver.Receive()

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusContinue && !req.bodyFittedInPreview && req.previewSet {
		return c.DoRemaining(req)
	}

	return resp, nil
}

// DoRemaining requests an ICAP server with the remaining body bytes which did not fit in the preview in the original request
func (c *Client) DoRemaining(req *Request) (*Response, error) {

	data := chunkBodyByBytes(req.remainingPreviewBytes, req.ChunkLength)

	if err := c.scktDriver.Send(data); err != nil {
		return nil, err
	}

	resp, err := c.scktDriver.Receive()

	if err != nil {
		return nil, err
	}

	return resp, nil
}
