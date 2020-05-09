package icapclient

import (
	"bufio"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestResponse(t *testing.T) {

	t.Run("ReadResponse", func(t *testing.T) {

		respStr := "ICAP/1.0 200 OK\n" +
			"Date: Mon, 10 Jan 2000  09:55:21 GMT\n" +
			"Server: ICAP-Server-Software/1.0\n" +
			"Connection: close\n" +
			"ISTag: \"W3E4R7U9-L2E4-2\"\n" +
			"Encapsulated: req-hdr=0, null-body=231\n\n" +
			"GET /modified-path HTTP/1.1\r\n" +
			"Host: www.origin-server.com\r\n" +
			"Via: 1.0 icap-server.net (ICAP Example ReqMod Service 1.1)\r\n" +
			"Accept: text/html, text/plain, image/gif\r\n" +
			"Accept-Encoding: gzip, compress\r\n" +
			"If-None-Match: \"xyzzy\", \"r2d2xxxx\"\r\n\r\n"

		resp, err := ReadResponse(bufio.NewReader(strings.NewReader(respStr)))

		if err != nil {
			t.Fatal(err.Error())
		}

		spew.Dump(resp)

	})

}
