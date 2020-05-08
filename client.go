package icapclient

import (
	"strconv"
)

// Client represents the icap client who makes the icap server calls
type Client struct {
	scktDriver *Driver
}

// Do makes the call
func (c *Client) Do(req *Request) (*Response, error) {
	return c.do(req)
}

func (c *Client) do(req *Request) (*Response, error) {

	port, err := strconv.Atoi(req.URL.Port())

	if err != nil {
		return nil, err
	}

	c.scktDriver = NewDriver(req.URL.Hostname(), port)

	if err := c.scktDriver.Connect(); err != nil {
		return nil, err
	}

	defer c.scktDriver.Close()

	if err := c.scktDriver.Send(req); err != nil {
		return nil, err
	}

	resp, err := c.scktDriver.Receive()

	if err != nil {
		return nil, err
	}

	return resp, nil

}
