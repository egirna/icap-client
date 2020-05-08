package icapclient

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

// Driver os the one responsible for driving the transport layer operations
type Driver struct {
	Host string
	Port int
	tcp  *transport
}

// NewDriver is the factory function for Driver
func NewDriver(host string, port int) *Driver {
	return &Driver{
		Host: host,
		Port: port,
	}
}

// Connect fires up a tcp socket connection with the icap server
func (d *Driver) Connect() error {

	d.tcp = &transport{
		network: "tcp",
		addr:    fmt.Sprintf("%s:%d", d.Host, d.Port),
	}

	return d.tcp.dial()
}

// Close closes the socket connection
func (d *Driver) Close() error {
	if d.tcp == nil {

		return errors.New(ErrConnectionNotOpen)
	}

	return d.tcp.close()
}

// Send sends a request to the icap server
func (d *Driver) Send(req *Request) error {

	b, err := DumpRequest(req)

	if err != nil {
		return err
	}

	_, err = d.tcp.write(string(b))

	if err != nil {
		return err
	}

	return nil

}

// Receive returns the respone from the tcp socket connection
func (d *Driver) Receive() (*Response, error) {

	msg, err := d.tcp.read()

	if err != nil {
		return nil, err
	}

	resp, err := ReadResponse(bufio.NewReader(strings.NewReader(msg)))

	if err != nil {
		return nil, err
	}

	return resp, nil
}
