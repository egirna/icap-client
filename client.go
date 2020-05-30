package icapclient

import (
	"net/http"
	"strconv"
	"time"
)

// Client represents the icap client who makes the icap server calls
type Client struct {
	scktDriver *Driver
	Timeout    time.Duration
}

// Do makes the call
func (c *Client) Do(req *Request) (*Response, error) {

	port, err := strconv.Atoi(req.URL.Port())

	if err != nil {
		return nil, err
	}

	if c.scktDriver == nil {
		c.scktDriver = NewDriver(req.URL.Hostname(), port)
	}

	c.setDefaultTimeouts()

	if req.ctx != nil {
		if err := c.scktDriver.ConnectWithContext(*req.ctx); err != nil {
			return nil, err
		}
	} else {
		if err := c.scktDriver.Connect(); err != nil {
			return nil, err
		}
	}

	defer c.scktDriver.Close()

	req.SetDefaultRequestHeaders()

	d, err := DumpRequest(req)

	if err != nil {
		return nil, err
	}

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

	data := req.remainingPreviewBytes

	if !bodyAlreadyChunked(string(data)) {
		ds := string(data)
		addHexaBodyByteNotations(&ds)
		data = []byte(ds)
	}

	if err := c.scktDriver.Send(data); err != nil {
		return nil, err
	}

	resp, err := c.scktDriver.Receive()

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SetDriver sets a new socket driver with the client
func (c *Client) SetDriver(d *Driver) {
	c.scktDriver = d
}

func (c *Client) setDefaultTimeouts() {
	if c.Timeout == 0 {
		c.Timeout = defaultTimeout
	}

	if c.scktDriver.DialerTimeout == 0 {
		c.scktDriver.DialerTimeout = c.Timeout
	}

	if c.scktDriver.ReadTimeout == 0 {
		c.scktDriver.ReadTimeout = c.Timeout
	}

	if c.scktDriver.WriteTimeout == 0 {
		c.scktDriver.WriteTimeout = c.Timeout
	}
}
