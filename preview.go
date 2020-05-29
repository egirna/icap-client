package icapclient

import (
	"bytes"
	"io/ioutil"
	"strconv"
)

// SetPreview sets the preview bytes in the icap header
func (r *Request) SetPreview(maxBytes int) error {

	bodyBytes := []byte{}

	previewBytes := 0

	if r.Method == MethodREQMOD {
		if r.HTTPRequest == nil {
			return nil
		}
		if r.HTTPRequest.Body != nil {
			var err error
			bodyBytes, err = ioutil.ReadAll(r.HTTPRequest.Body)

			if err != nil {
				return err
			}

			defer r.HTTPRequest.Body.Close()
		}
	}

	if r.Method == MethodRESPMOD {
		if r.HTTPResponse == nil {
			return nil
		}

		if r.HTTPResponse.Body != nil {
			var err error
			bodyBytes, err = ioutil.ReadAll(r.HTTPResponse.Body)

			if err != nil {
				return err
			}

			defer r.HTTPResponse.Body.Close()
		}
	}

	previewBytes = len(bodyBytes)

	if previewBytes > 0 {
		r.bodyFittedInPreview = true
	}

	if previewBytes > maxBytes {
		previewBytes = maxBytes
		r.bodyFittedInPreview = false
		r.remainingPreviewBytes = bodyBytes[maxBytes:]
	}

	if r.Method == MethodREQMOD {
		r.HTTPRequest.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	}

	if r.Method == MethodRESPMOD {
		r.HTTPResponse.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	}

	r.Header.Set("Preview", strconv.Itoa(previewBytes))
	r.PreviewBytes = previewBytes
	r.previewSet = true

	return nil

}
