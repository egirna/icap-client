package examples

import (
	"fmt"
	"log"
	"net/http"
	"time"

	ic "github.com/egirna/icap-client"
)

func makeRespmodWithCustomDriver() {
	httpReq, err := http.NewRequest(http.MethodGet, "http://localhost:8000/sample.pdf", nil)

	if err != nil {
		log.Fatal(err)
	}

	httpClient := &http.Client{}

	httpResp, err := httpClient.Do(httpReq)

	if err != nil {
		log.Fatal(err)
	}

	optReq, err := ic.NewRequest(ic.MethodOPTIONS, "icap://127.0.0.1:1344/respmod", nil, nil)

	if err != nil {
		log.Fatal(err)
		return
	}

	client := &ic.Client{
		Timeout: 5 * time.Second,
	}

	client.SetDriver(&ic.Driver{
		Host:        "127.0.0.1",
		Port:        1344,
		ReadTimeout: 2 * time.Second,
	})

	optResp, err := client.Do(optReq)

	if err != nil {
		log.Fatal(err)
		return
	}

	req, err := ic.NewRequest(ic.MethodRESPMOD, "icap://127.0.0.1:1344/respmod", httpReq, httpResp)

	if err != nil {
		log.Fatal(err)
	}

	req.SetPreview(optResp.PreviewBytes)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)

}
