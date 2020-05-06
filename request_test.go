package icapclient

import (
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {

	t.Run("Request Factory", func(t *testing.T) {
		if _, err := NewRequest("options", "icap://localhost:1344/something", nil); err != nil {
			t.Fatal(err.Error())
		}
		if _, err := NewRequest("respmod", "icap://localhost:1344/something", &http.Request{}); err != nil {
			t.Fatal(err.Error())
		}
		if _, err := NewRequest("reqmod", "icap://localhost:1344/something", &http.Request{}); err != nil {
			t.Fatal(err.Error())
		}
		if _, err := NewRequest("invalid", "icap://localhost:1344/something", nil); err == nil || err.Error() != ErrMethodNotRegistered {
			t.Fatal(err.Error())
		}
		if _, err := NewRequest("options", "http://localhost:1344/something", nil); err == nil || err.Error() != ErrInvalidScheme {
			t.Fatal(err.Error())
		}
		if _, err := NewRequest("options", "icap://", nil); err == nil || err.Error() != ErrInvalidHost {
			t.Fatal(err.Error())
		}
	})

}
