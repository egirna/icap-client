package icapclient

import (
	"bytes"
	"io/ioutil"
	"strconv"
)

// SetPreview sets the preview bytes in the icap header
func (r *Request) SetPreview(maxBytes int) error {

	if r.HTTPResponse == nil {
		return nil
	}

	bodyBytes, err := ioutil.ReadAll(r.HTTPResponse.Body)

	if err != nil {
		return err
	}

	defer r.HTTPResponse.Body.Close()

	previewBytes := len(bodyBytes)
	r.bodyFittedInPreview = true

	if len(bodyBytes) > maxBytes {
		previewBytes = maxBytes
		r.bodyFittedInPreview = false
		r.remainingPreviewBytes = bodyBytes[maxBytes:]
	}

	r.Header.Set("Preview", strconv.Itoa(previewBytes))
	r.PreviewBytes = previewBytes
	r.previewSet = true

	r.HTTPResponse.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))

	return nil

}
