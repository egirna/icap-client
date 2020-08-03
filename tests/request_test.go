package tests

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	ic "github.com/egirna/icap-client"
)

func TestRequest(t *testing.T) {

	t.Run("Request Factory", func(t *testing.T) {
		if _, err := ic.NewRequest(ic.MethodOPTIONS, "icap://localhost:1344/something", nil, nil); err != nil {
			t.Log(err.Error())
			t.Fail()
		}
		if _, err := ic.NewRequest(ic.MethodRESPMOD, "icap://localhost:1344/something", nil, &http.Response{}); err != nil {
			t.Log(err.Error())
			t.Fail()

		}
		if _, err := ic.NewRequest(ic.MethodREQMOD, "icap://localhost:1344/something", &http.Request{}, nil); err != nil {
			t.Log(err.Error())
			t.Fail()
		}
		if _, err := ic.NewRequest("invalid", "icap://localhost:1344/something", nil, nil); err == nil ||
			err.Error() != ic.ErrMethodNotRegistered {
			t.Log(err.Error())
			t.Fail()
		}
		if _, err := ic.NewRequest(ic.MethodOPTIONS, "http://localhost:1344/something", nil, nil); err == nil ||
			err.Error() != ic.ErrInvalidScheme {
			t.Log(err.Error())
			t.Fail()
		}
		if _, err := ic.NewRequest(ic.MethodOPTIONS, "icap://", nil, nil); err == nil || err.Error() != ic.ErrInvalidHost {
			t.Log(err.Error())
			t.Fail()
		}
		if _, err := ic.NewRequest(ic.MethodREQMOD, "icap://localhost:1344/something", nil, nil); err == nil ||
			err.Error() != ic.ErrREQMODWithNoReq {
			t.Log(err.Error())
			t.Fail()
		}
		if _, err := ic.NewRequest(ic.MethodREQMOD, "icap://localhost:1344/something", &http.Request{}, &http.Response{}); err == nil ||
			err.Error() != ic.ErrREQMODWithResp {
			t.Log(err.Error())
			t.Fail()
		}
		if _, err := ic.NewRequest(ic.MethodRESPMOD, "icap://localhost:1344/something", &http.Request{}, nil); err == nil ||
			err.Error() != ic.ErrRESPMODWithNoResp {
			t.Log(err.Error())
			t.Fail()
		}
	})

	t.Run("DumpRequest OPTIONS", func(t *testing.T) {

		req, _ := ic.NewRequest(ic.MethodOPTIONS, "icap://localhost:1344/something", nil, nil)

		b, err := ic.DumpRequest(req)

		if err != nil {
			t.Fatal(err.Error())
		}

		wanted := "OPTIONS icap://localhost:1344/something ICAP/1.0\r\n" +
			"Encapsulated:  null-body=0\r\n\r\n"

		got := string(b)

		if wanted != got {
			t.Logf("wanted: %s, got: %s\n", wanted, got)
			t.Fail()
		}

	})

	t.Run("DumpRequest REQMOD", func(t *testing.T) { // FIXME: add proper wanted string and complete this unit test
		httpReq, _ := http.NewRequest(http.MethodGet, "http://someurl.com", nil)

		req, _ := ic.NewRequest(ic.MethodREQMOD, "icap://localhost:1344/something", httpReq, nil)

		b, err := ic.DumpRequest(req)
		if err != nil {
			t.Fatal(err.Error())
		}

		wanted := "REQMOD icap://localhost:1344/something ICAP/1.0\r\n" +
			"Encapsulated:  req-hdr=0, null-body=109\r\n\r\n" +
			"GET http://someurl.com HTTP/1.1\r\n" +
			"Host: someurl.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Accept-Encoding: gzip\r\n\r\n"

		got := string(b)

		if wanted != got {
			t.Logf("wanted: \n%s\ngot: \n%s\n", wanted, got)
			t.Fail()
		}

		httpReq, _ = http.NewRequest(http.MethodPost, "http://someurl.com", bytes.NewBufferString("Hello World"))

		req, _ = ic.NewRequest(ic.MethodREQMOD, "icap://localhost:1344/something", httpReq, nil)

		b, err = ic.DumpRequest(req)
		if err != nil {
			t.Fatal(err.Error())
		}

		wanted = "REQMOD icap://localhost:1344/something ICAP/1.0\r\n" +
			"Encapsulated:  req-hdr=0, req-body=130\r\n\r\n" +
			"POST http://someurl.com HTTP/1.1\r\n" +
			"Host: someurl.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Content-Length: 11\r\n" +
			"Accept-Encoding: gzip\r\n\r\n" +
			"b\r\n" +
			"Hello World\r\n" +
			"0\r\n\r\n"

		got = string(b)

		if wanted != got {
			t.Logf("wanted: \n%s\ngot: \n%s\n", wanted, got)
			t.Fail()
		}

	})

	t.Run("DumpRequest RESPMOD", func(t *testing.T) {
		httpReq, _ := http.NewRequest(http.MethodPost, "http://someurl.com", bytes.NewBufferString("Hello World"))
		httpResp := &http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Proto:      "HTTP/1.0",
			ProtoMajor: 1,
			ProtoMinor: 0,
			Header: http.Header{
				"Content-Type":   []string{"plain/text"},
				"Content-Length": []string{"11"},
			},
			ContentLength: 11,
			Body:          ioutil.NopCloser(strings.NewReader("Hello World")),
		}

		req, _ := ic.NewRequest(ic.MethodRESPMOD, "icap://localhost:1344/something", httpReq, httpResp)

		b, err := ic.DumpRequest(req)
		if err != nil {
			t.Fatal(err.Error())
		}

		wanted := "RESPMOD icap://localhost:1344/something ICAP/1.0\r\n" +
			"Encapsulated:  req-hdr=0, req-body=130, res-hdr=145, res-body=210\r\n\r\n" +
			"POST http://someurl.com HTTP/1.1\r\n" +
			"Host: someurl.com\r\n" +
			"User-Agent: Go-http-client/1.1\r\n" +
			"Content-Length: 11\r\n" +
			"Accept-Encoding: gzip\r\n\r\n" +
			"Hello World\r\n\r\n" +
			"HTTP/1.0 200 OK\r\n" +
			"Content-Length: 11\r\n" +
			"Content-Type: plain/text\r\n\r\n" +
			"b\r\n" +
			"Hello World\r\n" +
			"0\r\n\r\n"

		got := string(b)

		if wanted != got {
			t.Logf("wanted: \n%s\ngot: \n%s\n", wanted, got)
			t.Fail()
		}

	})

	t.Run("SetDefaultRequestHeaders", func(t *testing.T) {
		req, _ := ic.NewRequest(ic.MethodOPTIONS, "icap://localhost:1344/something", nil, nil)
		req.SetDefaultRequestHeaders()

		if val, exists := req.Header["Allow"]; !exists || len(val) < 1 || val[0] != "204" {
			t.Log("Must have Allow header with 204 as value")
			t.Fail()
		}

		hname, _ := os.Hostname()
		if val, exists := req.Header["Host"]; !exists || len(val) < 1 || val[0] != hname {
			t.Logf("Must have Host header with %s as value", hname)
			t.Fail()
		}

		req, _ = ic.NewRequest(ic.MethodOPTIONS, "icap://localhost:1344/something", nil, nil)
		req.Header.Set("Host", "somehost")
		req.SetDefaultRequestHeaders()

		if val, exists := req.Header["Host"]; !exists || len(val) < 1 || val[0] != "somehost" {
			t.Logf("Must have Host header with %s as value", "somehost")
			t.Fail()
		}

	})

	t.Run("ExtendHeader", func(t *testing.T) {
		hdr := http.Header{
			"Name":    []string{"some_name"},
			"Address": []string{"some_address1", "some_address2"},
			"Allow":   []string{"205"},
		}

		req, _ := ic.NewRequest(ic.MethodOPTIONS, "icap://localhost:1344/something", nil, nil)
		req.SetDefaultRequestHeaders()
		if err := req.ExtendHeader(hdr); err != nil {
			t.Fatal(err.Error())
		}

		if val, exists := req.Header["Allow"]; !exists || len(val) < 2 || !reflect.DeepEqual(val, []string{"204", "205"}) {
			t.Log("Must have Allow header with {204,205} as value")
			t.Fail()
		}

		hname, _ := os.Hostname()
		if val, exists := req.Header["Host"]; !exists || len(val) < 1 || val[0] != hname {
			t.Logf("Must have Host header with %s as value", hname)
			t.Fail()
		}

		if val, exists := req.Header["Name"]; !exists || len(val) < 1 || val[0] != "some_name" {
			t.Log("Must have Name header with some_name as value")
			t.Fail()
		}

		if val, exists := req.Header["Address"]; !exists || len(val) < 2 || !reflect.DeepEqual(val,
			[]string{"some_address1", "some_address2"}) {
			t.Log("Must have Address header with {some_address1, some_address2} as value")
			t.Fail()
		}

	})

	t.Run("SetPreview REQMOD", func(t *testing.T) {
		bodyStr := "Hello World! Bye Bye World!"
		bodyData := bytes.NewBufferString(bodyStr)
		httpReq, _ := http.NewRequest(http.MethodPost, "http://someurl.com", bodyData)
		req, _ := ic.NewRequest(ic.MethodREQMOD, "icap://localhost:1344/something", httpReq, nil)

		if err := req.SetPreview(11); err != nil {
			t.Fatal(err.Error())
		}

		if req.PreviewBytes != 11 {
			t.Logf("Wanted preview bytes:%d, got:%d", 11, req.PreviewBytes)
			t.Fail()
		}

		bdyBytes, _ := ioutil.ReadAll(req.HTTPRequest.Body)

		if string(bdyBytes) != bodyStr {
			t.Logf("Wanted body string:%s, got:%s", bodyStr, string(bdyBytes))
			t.Fail()
		}

	})

	// TODO: add the rest of the request unit tests starting with SetPreview RESPMOD

}
