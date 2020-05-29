package icapclient

import (
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

		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		data = append(data, tmp[:n]...)
		if string(data) == icap100ContinueMsg { // explicitly breaking because the Read blocks for 100 continue message // TODO: find out why
			break
		}

	}

	return string(data), nil
}

func (t *transport) close() error {
	return t.sckt.Close()
}
