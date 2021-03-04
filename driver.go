package icapclient

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Driver os the one responsible for driving the transport layer operations
type Driver struct {
	Host          string
	Port          int
	DialerTimeout time.Duration
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	tcp           *transport
}

// NewDriver is the factory function for Driver
func NewDriver(host string, port int) *Driver {
	return &Driver{
		Host: host,
		Port: port,
	}
}

//TlSDial Connecting with a custom root-certificate set.
func (d *Driver) TlSDial() error {
	d.tcp = &transport{
		network:      "tcp",
		addr:         fmt.Sprintf("%s:%d", d.Host, d.Port),
		timeout:      d.DialerTimeout,
		readTimeout:  d.ReadTimeout,
		writeTimeout: d.WriteTimeout,
	}

	return d.tcp.tlsdial() //return conn

}

// Connect fires up a tcp socket connection with the icap server
func (d *Driver) Connect() error {
	fmt.Println("Connect")

	d.tcp = &transport{
		network:      "tcp",
		addr:         fmt.Sprintf("%s:%d", d.Host, d.Port),
		timeout:      d.DialerTimeout,
		readTimeout:  d.ReadTimeout,
		writeTimeout: d.WriteTimeout,
	}

	return d.tcp.dial()
}

// ConnectWithContext connects to the server satisfying the context
func (d *Driver) ConnectWithContext(ctx context.Context) error {
	d.tcp = &transport{
		network:      "tcp",
		addr:         fmt.Sprintf("%s:%d", d.Host, d.Port),
		timeout:      d.DialerTimeout,
		readTimeout:  d.ReadTimeout,
		writeTimeout: d.WriteTimeout,
	}

	return d.tcp.dialWithContext(ctx)
}

// Close closes the socket connection
func (d *Driver) Close() error {
	fmt.Println("close")

	if d.tcp == nil {

		return errors.New(ErrConnectionNotOpen)
	}

	return d.tcp.close()
}

// Send sends a request to the icap server
func (d *Driver) Send(data []byte) error {
	fmt.Println("send")

	_, err := d.tcp.write(data)

	if err != nil {
		return err
	}

	return nil

}

// Receive returns the respone from the tcp socket connection
func (d *Driver) Receive() (*Response, error) {
	fmt.Println("Receive")

	mgs := d.tcp.read()

	resp, err := ReadRespons(mgs)

	if err != nil {
		return nil, err
	}

	logDebug("The final *ic.Response from tcp messages...")
	//	dumpDebug(resp)

	dumpDebug(resp)

	//	icapresp := new(IcapRequest)

	return resp, nil
}
