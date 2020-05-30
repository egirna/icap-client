package icapclient

import (
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {

	t.Run("Request Factory", func(t *testing.T) {
		if _, err := NewRequest("options", "icap://localhost:1344/something", nil, nil); err != nil {
			t.Fatal(err.Error())
		}
		if _, err := NewRequest("respmod", "icap://localhost:1344/something", nil, &http.Response{}); err != nil {
			t.Fatal(err.Error())
		}
		if _, err := NewRequest("reqmod", "icap://localhost:1344/something", &http.Request{}, nil); err != nil {
			t.Fatal(err.Error())
		}
		if _, err := NewRequest("invalid", "icap://localhost:1344/something", nil, nil); err == nil || err.Error() != ErrMethodNotRegistered {
			t.Fatal(err.Error())
		}
		if _, err := NewRequest("options", "http://localhost:1344/something", nil, nil); err == nil || err.Error() != ErrInvalidScheme {
			t.Fatal(err.Error())
		}
		if _, err := NewRequest("options", "icap://", nil, nil); err == nil || err.Error() != ErrInvalidHost {
			t.Fatal(err.Error())
		}
	})

	t.Run("DumpRequest", func(t *testing.T) {

		req, _ := NewRequest(MethodOPTIONS, "icap://localhost:1344/something", nil, nil)

		b, err := DumpRequest(req)

		if err != nil {
			t.Fatal(err.Error())
		}

		wanted := "OPTIONS icap://localhost:1344/something ICAP/1.0\r\n" +
			"Encapsulated:  null-body=0\r\n\r\n"

		got := string(b)

		if wanted != got {
			t.Fail()
		}

	})

}
