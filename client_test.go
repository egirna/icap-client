package icapclient

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestClientDo(t *testing.T) {

	t.Run("Client Request", func(t *testing.T) {
		req, err := NewRequest("options", "icap://127.0.0.1:1344/respmod-icapeg", nil)

		if err != nil {
			t.Fatal(err.Error())
		}

		client := Client{}

		msg, err := client.Do(req)

		if err != nil {
			t.Fatal(err.Error())
		}

		spew.Dump(msg)

	})

}
