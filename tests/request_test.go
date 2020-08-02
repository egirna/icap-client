package tests

import (
	"net/http"
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
			"Encapsulated: req-hdr=0, null-body=170"

		got := string(b)

		if wanted != got {
			t.Logf("wanted: \n%s\ngot: \n%s\n", wanted, got)
			t.Fail()
		}

	})

	t.Run("DumpRequest RESPMOD", func(t *testing.T) {
		// TODO: add respmod DumpReqest function tests here
	})

}
