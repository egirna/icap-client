package icapclient

import (
	"fmt"
	"strings"
)

type Client struct {
	Host string
	tcp  *transport
}

func (c *Client) Do(req *Request) (string, error) {
	return c.do(req)
}

func (c *Client) do(req *Request) (string, error) {

	if c.tcp == nil {
		c.tcp = &transport{
			network: "tcp",
			addr:    req.URL.Host,
		}

		if err := c.tcp.dial(); err != nil {
			return "", err
		}
	}

	method := strings.ToUpper(req.Method)
	uri := req.URL.RequestURI()

	_, err := c.tcp.write(fmt.Sprintf("%s %s %s\n\n", method, uri, ICAPVersion))

	if err != nil {
		return "", err
	}

	msg, err := c.tcp.read()

	if err != nil {
		return "", err
	}

	if err := c.tcp.close(); err != nil {
		return "", err
	}

	return msg, nil
}
