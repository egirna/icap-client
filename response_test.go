package icapclient

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestResponse(t *testing.T) {

	t.Run("ReadResponse getting reqmod response", func(t *testing.T) {

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

		wantedICAPStatusCode := 200
		if resp.StatusCode != wantedICAPStatusCode {
			t.Errorf("Expected ICAP server response status code to be %d got %d", wantedICAPStatusCode, resp.StatusCode)
		}
		wantedICAPStatus := "OK"
		if resp.Status != wantedICAPStatus {
			t.Errorf("Expected ICAP server response status to be %s, got %s", wantedICAPStatus, resp.Status)
		}

		if resp.ContentRequest != nil {
			headerValue := "text/html, text/plain, image/gif"
			if resp.ContentRequest.Header.Get("Accept") != "text/html, text/plain, image/gif" {
				t.Errorf("Expected value of http request header Accecpt to be %s , got %s", headerValue,
					resp.ContentRequest.Header.Get("Accept"))
			}

			wantedURI := "/modified-path"
			if resp.ContentRequest.RequestURI != wantedURI {
				t.Errorf("Expected http request requestedURI to be %s, got %s", wantedURI, resp.ContentRequest.RequestURI)
			}
			wantedMethod := http.MethodGet
			if resp.ContentRequest.Method != wantedMethod {
				t.Errorf("Expected http request method  to be %s, got %s", wantedMethod, resp.ContentRequest.Method)
			}

			wantedHost := "www.origin-server.com"
			if resp.ContentRequest.Host != wantedHost {
				t.Errorf("Expected http request host  to be %s, got %s", wantedHost, resp.ContentRequest.Host)
			}
		}

		if resp.ContentRequest == nil {
			t.Error("The http request should not be nil")
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
