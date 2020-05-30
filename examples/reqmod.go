package examples

import (
	"fmt"
	"log"
	"net/http"
	"time"

	ic "github.com/egirna/icap-client"
)

func makeReqmodCall() {
	httpReq, err := http.NewRequest(http.MethodGet, "http://localhost:8000/sample.pdf", nil)

	if err != nil {
		log.Fatal(err)
	}

	optReq, err := ic.NewRequest(ic.MethodOPTIONS, "icap://127.0.0.1:1344/reqmod", nil, nil)

	if err != nil {
		log.Fatal(err)
		return
	}

	client := &ic.Client{
		Timeout: 5 * time.Second,
	}

	optResp, err := client.Do(optReq)

	if err != nil {
		log.Fatal(err)
		return
	}

	req, err := ic.NewRequest(ic.MethodREQMOD, "icap://127.0.0.1:1344/reqmod", httpReq, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.SetPreview(optResp.PreviewBytes)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Status)

}
