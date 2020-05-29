package icapclient

import (
	"fmt"
	"io"
	"net"
)

// transport represents the transport layer data
type transport struct {
	network string
	addr    string
	sckt    net.Conn
}

// Dial fires up a tcp socket
func (t *transport) dial() error {
	sckt, err := net.Dial(t.network, t.addr)

	if err != nil {
		return err
	}

	t.sckt = sckt

	return nil
}

// Write writes data to the server
func (t *transport) write(data []byte) (int, error) {
	return t.sckt.Write(data)
}

// Read reads data from server
func (t *transport) read() (string, error) {

	data := make([]byte, 0)

	for {
		tmp := make([]byte, 1096)

		n, err := t.sckt.Read(tmp)

		fmt.Println("The index: ", n)

		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		data = append(data, tmp[:n]...)
		fmt.Println("The data: ", string(data))
		fmt.Println("truth: ", string(data) == "ICAP/1.0 100 Continue\r\n\r\n")
		if string(data) == "ICAP/1.0 100 Continue\r\n\r\n" {
			break
		}
	}

	return string(data), nil
}

func (t *transport) close() error {
	return t.sckt.Close()
}
