package icapclient

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestResponse(t *testing.T) {

	t.Run("ReadResponse REQMOD", func(t *testing.T) { // FIXME: headers and content request aren't being tested properly

		type testSample struct {
			headers      http.Header
			status       string
			statusCode   int
			previewBytes int
			respStr      string
			httpReqStr   string
		}

		sampleTable := []testSample{
			{
				headers: http.Header{
					"Date":         []string{"Mon, 10 Jan 2000  09:55:21 GMT"},
					"Server":       []string{"ICAP-Server-Software/1.0"},
					"ISTag":        []string{"\"W3E4R7U9-L2E4-2\""},
					"Encapsulated": []string{"req-hdr=0, null-body=231"},
				},
				status:       "OK",
				statusCode:   200,
				previewBytes: 0,
				respStr: "ICAP/1.0 200 OK\r\n" +
					"Date: Mon, 10 Jan 2000  09:55:21 GMT\r\n" +
					"Server: ICAP-Server-Software/1.0\r\n" +
					"Connection: close\r\n" +
					"ISTag: \"W3E4R7U9-L2E4-2\"\r\n" +
					"Encapsulated: req-hdr=0, null-body=231\r\n\r\n",
				httpReqStr: "GET /modified-path HTTP/1.1\r\n" +
					"Host: www.origin-server.com\r\n" +
					"Via: 1.0 icap-server.net (ICAP Example ReqMod Service 1.1)\r\n" +
					"Accept: text/html, text/plain, image/gif\r\n" +
					"Accept-Encoding: gzip, compress\r\n" +
					"If-None-Match: \"xyzzy\", \"r2d2xxxx\"\r\n\r\n",
			},
		}

		for _, sample := range sampleTable {
			resp, err := ReadResponse(bufio.NewReader(strings.NewReader(sample.respStr)))
			if err != nil {
				t.Fatal(err.Error())
			}

			if resp.StatusCode != sample.statusCode {
				t.Logf("Wanted ICAP status code: %d , got: %d", sample.statusCode, resp.StatusCode)
				t.Fail()
			}
			if resp.Status != sample.status {
				t.Logf("Wanted ICAP status: %s , got: %s", sample.status, resp.Status)
				t.Fail()
			}
			if resp.PreviewBytes != sample.previewBytes {
				t.Logf("Wanted preview bytes: %d, got: %d", sample.previewBytes, resp.PreviewBytes)
				t.Fail()
			}
			if !reflect.DeepEqual(resp.Header, sample.headers) {
				t.Logf("Wanted ICAP header: %v, got: %v", sample.headers, resp.Header)
				t.Fail()
			}
			if resp.ContentRequest == nil {
				t.Log("ContentRequest is nil")
				t.Fail()
			}

			wantedHTTPReq, err := http.ReadRequest(bufio.NewReader(strings.NewReader(sample.httpReqStr)))
			if err != nil {
				t.Fatal(err.Error())
			}

			if !reflect.DeepEqual(resp.ContentRequest, wantedHTTPReq) {
				t.Logf("Wanted http request: %v, got: %v", wantedHTTPReq, resp.ContentRequest)
				t.Fail()
			}

		}

	})

	t.Run("ReadRequest getting respmod response", func(t *testing.T) {
		respStr := "ICAP/1.0 200 OK\n" +
			"Date: Mon, 10 Jan 2000  09:55:21 GMT\n" +
			"Server: ICAP-Server-Software/1.0\n" +
			"Connection: close\n" +
			"ISTag: \"W3E4R7U9-L2E4-2\"\n" +
			"Encapsulated: res-hdr=0, res-body=222\n\n" +

			"HTTP/1.1 200 OK\r\n" +
			"Date: Mon, 10 Jan 2000  09:55:21 GMT\r\n" +
			"Via: 1.0 icap.example.org (ICAP Example RespMod Service 1.1)\r\n" +
			"Server: Apache/1.3.6 (Unix)\r\n" +
			"ETag: \"63840-1ab7-378d415b\"\r\n" +
			"Content-Type: text/plain\r\n" +
			"Content-Length: 92\r\n\r\n" +

			"This is data that was returned by an origin server, but with value added by an ICAP server."

		resp, err := ReadResponse(bufio.NewReader(strings.NewReader(respStr)))

		if err != nil {
			t.Fatal(err.Error())
		}

		wantedICAPStatusCode := 200
		if resp.StatusCode != 200 {
			t.Errorf("Expected ICAP server response status code to be %d got %d", wantedICAPStatusCode, resp.StatusCode)
		}
		wantedICAPStatus := "OK"
		if resp.Status != "OK" {
			t.Errorf("Expected ICAP server response status to be %s, got %s", wantedICAPStatus, resp.Status)
		}

		if resp.ContentResponse != nil {
			wantedStatusCode := http.StatusOK
			if resp.ContentResponse.StatusCode != wantedStatusCode {
				t.Errorf("Expected http response status code to be %d, got %d", wantedStatusCode, resp.ContentResponse.StatusCode)
			}

			wantedStatus := fmt.Sprintf("%d %s", http.StatusOK, http.StatusText(http.StatusOK))
			if resp.ContentResponse.Status != wantedStatus {
				t.Errorf("Expected http response status to be %s, got %s", wantedStatus, resp.ContentResponse.Status)
			}
			if resp.ContentResponse.ContentLength != 92 {
				t.Errorf("Expected http response content length to be %d, got %d", 92, resp.ContentResponse.ContentLength)
			}

			wantedHeaderValue := "text/plain"
			if resp.ContentResponse.Header.Get("Content-Type") != wantedHeaderValue {
				t.Errorf("Expected http response content-type header value to be %s, got %s",
					wantedHeaderValue, resp.ContentResponse.Header.Get("Content-Type"))
			}

			bdyBytes, err := ioutil.ReadAll(resp.ContentResponse.Body)

			if err != nil {
				t.Fatal(err.Error())
			}

			wantedBodyData := "This is data that was returned by an origin server, but with value added by an ICAP server"
			gotBodyData := strings.TrimSpace(string(bdyBytes))
			if gotBodyData != wantedBodyData {
				t.Errorf("Expected http response body data to be %s, got %s", wantedBodyData, gotBodyData)
			}

		}

		if resp.ContentResponse == nil {
			t.Error("The http response should not be nil")
		}

	})

}
