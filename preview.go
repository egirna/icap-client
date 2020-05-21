package icapclient

import (
	"bytes"
	"io/ioutil"
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

// SetPreview sets the preview bytes in the icap header
func (r *Request) SetPreview(maxBytes int) error {

	// FIXME: Preview bytes not being set properly as the addHexaResponseBodyByteNotations is not copying but modifiy the original response

	if r.HTTPResponse == nil {
		return nil
	}

	bodyBytes, err := ioutil.ReadAll(r.HTTPResponse.Body)

	if err != nil {
		return err
	}

	defer r.HTTPResponse.Body.Close()

	respWithNotation, err := addHexaResponseBodyByteNotations(*r.HTTPResponse)

	if err != nil {
		return err
	}

	bdyBytesWithNotation, err := ioutil.ReadAll(respWithNotation.Body)

	if err != nil {
		return err
	}

	defer respWithNotation.Body.Close()

	spew.Dump(string(bodyBytes))

	previewBytes := len(bdyBytesWithNotation)
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
