package icapclient

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// transport represents the transport layer data
type transport struct {
	network      string
	addr         string
	timeout      time.Duration
	readTimeout  time.Duration
	writeTimeout time.Duration
	sckt         net.Conn
	buf          *bufio.ReadWriter
}

func (t *transport) tlsdial() error {
	sckt, err := tls.Dial(t.network, t.addr, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return err
	}

	if err := sckt.SetReadDeadline(time.Now().UTC().Add(t.readTimeout)); err != nil {
		return err
	}

	if err := sckt.SetWriteDeadline(time.Now().UTC().Add(t.writeTimeout)); err != nil {
		return err
	}

	t.sckt = sckt

	return nil
}

// dial fires up a tcp socket
func (t *transport) dial() error {
	sckt, err := net.DialTimeout(t.network, t.addr, t.timeout)

	if err != nil {
		return err
	}

	if err := sckt.SetReadDeadline(time.Now().UTC().Add(t.readTimeout)); err != nil {
		return err
	}

	if err := sckt.SetWriteDeadline(time.Now().UTC().Add(t.writeTimeout)); err != nil {
		return err
	}

	t.sckt = sckt

	return nil
}

// dialWithContext fires up a tcp socket
func (t *transport) dialWithContext(ctx context.Context) error {
	sckt, err := (&net.Dialer{
		Timeout: t.timeout,
	}).DialContext(ctx, t.network, t.addr)

	if err != nil {
		return err
	}

	if err := sckt.SetReadDeadline(time.Now().UTC().Add(t.readTimeout)); err != nil {
		return err
	}

	if err := sckt.SetWriteDeadline(time.Now().UTC().Add(t.writeTimeout)); err != nil {
		return err
	}

	t.sckt = sckt

	return nil
}

// Write writes data to the server
func (t *transport) write(data []byte) (int, error) {
	logDebug("Dumping the message being sent to the server...")
	dumpDebug(string(data))

	return t.sckt.Write(data)
}

func (t *transport) read() *bufio.ReadWriter {

	//var bodyReader io.ReadCloser = emptyReader(0)

	logDebug("Dumping messages received from the server...")
	//	br := bufio.NewReader(t.sckt)

	br := bufio.NewReader(t.sckt)
	bw := bufio.NewWriter(t.sckt)
	c := bufio.NewReadWriter(br, bw)
	t.buf = c

	/*f := rankings.Response.Body
	filepath := "client/re_body.pdf"
	samplefile, _ := os.Create(filepath)

	defer samplefile.Close()

	//	io.Copy(samplefile, resp.ContentResponse.Body)
	samplefile.ReadFrom(f)*/

	return c
}

// Read reads data from server
func (t *transport) readold() (string, error) {
	fmt.Println("begin read")

	data := make([]byte, 0)
	//var bodyReader io.ReadCloser = emptyReader(0)

	logDebug("Dumping messages received from the server...")
	//	br := bufio.NewReader(t.sckt)
	for {
		fmt.Println("read")

		tmp := make([]byte, 1096)

		n, err := t.sckt.Read(tmp)
		//br := bufio.NewReader(tmp)

		if err != nil {
			if err == io.EOF {
				logDebug("End of file detected from EOF error")
				fmt.Println("End of file detected from EOF error")
				break
			}
			return "", err
		}

		if n == 0 {
			logDebug("End of file detected by 0 bytes")
			fmt.Println("End of file detected by 0 bytes")
		}

		data = append(data, tmp[:n]...)
		if string(data) == icap100ContinueMsg { // explicitly breaking because the Read blocks for 100 continue message // TODO: find out why
			logDebug("Stopping because got 100 Continue from the server")
			fmt.Println("Stopping because got 100 Continue from the server")

			break
		}

		if strings.HasSuffix(string(data), "0\r\n\r\n") {
			logDebug("End of the file detected by 0 Double CRLF indicator")
			fmt.Println("End of the file detected by 0 Double CRLF indicator")

			break
		}

		if strings.Contains(string(data), icap204NoModsMsg) {
			logDebug("End of file detected by 204 no modifications and Double CRLF at the end")
			fmt.Println("End of file detected by 204 no modifications and Double CRLF at the end")

			break
		}

		dumpDebug(string(tmp))

	}

	//fmt.Println(string(data))

	return string(data), nil
}

// close closes the tcp connection
func (t *transport) close() error {
	t.buf.Flush()
	t.buf = nil
	return t.sckt.Close()
}
