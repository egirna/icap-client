package icapclient

import (
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

	chunkLength := req.ChunkLength

	if chunkLength <= 0 {
		chunkLength = defaultChunkLength
	}

	data := chunkBodyByBytes(req.remainingPreviewBytes, chunkLength)

	if err := c.scktDriver.Send(data); err != nil {
		return nil, err
	}

	resp, err := c.scktDriver.Receive()

	if err != nil {
		return nil, err
	}

	return resp, nil
}
