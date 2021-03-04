package icapclient

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"net/url"
	"strconv"
	"strings"
)

// Response represents the icap server response data
type Responseold struct {
	StatusCode      int
	Status          string
	PreviewBytes    int
	Header          textproto.MIMEHeader
	ContentRequest  *http.Request
	ContentResponse *http.Response
}

var (
	optionValues = map[string]bool{
		PreviewHeader:          true,
		MethodsHeader:          true,
		AllowHeader:            true,
		TransferPreviewHeader:  true,
		ServiceHeader:          true,
		ISTagHeader:            true,
		OptBodyTypeHeader:      true,
		MaxConnectionsHeader:   true,
		OptionsTTLHeader:       true,
		ServiceIDHeader:        true,
		TransferIgnoreHeader:   true,
		TransferCompleteHeader: true,
	}
)

// An emptyReader is an io.ReadCloser that always returns os.EOF.
type emptyReader byte

func (emptyReader) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (emptyReader) Close() error {
	return nil
}

// A continueReader sends a "100 Continue" message the first time Read
// is called, creates a ChunkedReader, and reads from that.
type continueReader struct {
	buf *bufio.ReadWriter // the underlying connection
	cr  io.Reader         // the ChunkedReader
}

func (c *continueReader) Read(p []byte) (n int, err error) {
	if c.cr == nil {
		_, err := c.buf.WriteString("ICAP/1.0 100 Continue\r\n\r\n")
		if err != nil {
			return 0, err
		}
		err = c.buf.Flush()
		if err != nil {
			return 0, err
		}
		c.cr = newChunkedReader(c.buf)
	}

	return c.cr.Read(p)
}

type badStringError struct {
	what string
	str  string
}

func (e *badStringError) Error() string { return fmt.Sprintf("%s %q", e.what, e.str) }

//Response represents a parsed ICAP request.
type Response struct {
	Method       string               // REQMOD, RESPMOD, OPTIONS, etc.
	RawURL       string               // The URL given in the request.
	URL          *url.URL             // Parsed URL.
	Proto        string               // The protocol version.
	Header       textproto.MIMEHeader // The ICAP header
	RemoteAddr   string               // the address of the computer sending the request
	Preview      []byte               // the body data for an ICAP preview
	StatusCode   int
	Status       string
	PreviewBytes int
	// The HTTP messages.
	Body            []byte
	ContentRequest  *http.Request
	ContentResponse *http.Response
}

//ReadRespons chunked
func ReadRespons(b *bufio.ReadWriter) (resp *Response, err error) {
	tp := textproto.NewReader(b.Reader)

	resp = new(Response)

	// Read first line.
	var s string
	s, err = tp.ReadLine()
	if err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}

	//samplefile.ReadFrom(s)
	f := strings.SplitN(string(s), " ", 3)
	if len(f) < 3 {
		return nil, err
	}
	getStatusWithCode(f[1], f[2])
	resp.StatusCode, resp.Status, err = getStatusWithCode(f[1], f[2])

	resp.Method = f[0]

	resp.Header, err = tp.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(req.Header))

	s = resp.Header.Get("Encapsulated")
	if s == "" {
		return resp, nil // No HTTP headers or body.
	}

	eList := strings.Split(s, ", ")
	var initialOffset, reqHdrLen, respHdrLen int
	var hasBody bool
	var prevKey string
	var prevValue int
	for _, item := range eList {
		eq := strings.Index(item, "=")
		if eq == -1 {
			return nil, &badStringError{"malformed Encapsulated: header", s}
		}
		key := item[:eq]
		value, err := strconv.Atoi(item[eq+1:])
		if err != nil {
			return nil, &badStringError{"malformed Encapsulated: header", s}
		}

		// Calculate the length of the previous section.
		switch prevKey {
		case "":
			initialOffset = value
		case "req-hdr":
			reqHdrLen = value - prevValue
		case "res-hdr":
			respHdrLen = value - prevValue
		case "req-body", "opt-body", "res-body", "null-body":
			return nil, fmt.Errorf("%s must be the last section", prevKey)
		}

		switch key {
		case "req-hdr", "res-hdr", "null-body":
		case "req-body", "res-body", "opt-body":
			hasBody = true
		default:
			return nil, &badStringError{"invalid key for Encapsulated: header", key}
		}

		prevValue = value
		prevKey = key

	}

	// Read the HTTP headers.

	var rawReqHdr, rawRespHdr []byte
	if initialOffset > 0 {
		junk := make([]byte, initialOffset)
		_, err = io.ReadFull(b, junk)
		if err != nil {
			return nil, err
		}
	}
	if reqHdrLen > 0 {
		rawReqHdr = make([]byte, reqHdrLen)
		_, err = io.ReadFull(b, rawReqHdr)
		if err != nil {
			return nil, err
		}
	}
	if respHdrLen > 0 {
		rawRespHdr = make([]byte, respHdrLen)
		_, err = io.ReadFull(b, rawRespHdr)
		if err != nil {
			return nil, err
		}
	}

	var bodyReader io.ReadCloser = emptyReader(0)
	if hasBody {
		if p := resp.Header.Get("Preview"); p != "" {
			moreBody := true
			resp.Preview, err = ioutil.ReadAll(newChunkedReader(b))
			if err != nil {
				if strings.Contains(err.Error(), "ieof") {
					// The data ended with "0; ieof", which the HTTP chunked reader doesn't understand.
					moreBody = false
					err = nil
				} else {
					return nil, err
				}
			}
			var r io.Reader = bytes.NewBuffer(resp.Preview)
			if moreBody {
				r = io.MultiReader(r, &continueReader{buf: b})
			}
			bodyReader = ioutil.NopCloser(r)
		} else {
			bodyReader = ioutil.NopCloser(newChunkedReader(b))

			/*filepath := "client/re_now.txt"
			samplefile, _ := os.Create(filepath)

			defer samplefile.Close()

			//	io.Copy(samplefile, resp.ContentResponse.Body)
			samplefile.ReadFrom(bodyReader)*/
			// check errors
		}
	}

	// Construct the http.Request.
	if rawReqHdr != nil {

		resp.ContentRequest, err = http.ReadRequest(bufio.NewReader(bytes.NewBuffer(rawReqHdr)))
		if err != nil {
			return nil, fmt.Errorf("error while parsing HTTP request: %v", err)
		}

		if resp.Method == "REQMOD" {
			resp.ContentRequest.Body = bodyReader
		} else {
			resp.ContentRequest.Body = emptyReader(0)
		}
	}

	// Construct the http.Response.

	if rawRespHdr != nil {
		request := resp.ContentRequest
		if request == nil {
			request, _ = http.NewRequest("GET", "/", nil)
		}

		resp.ContentResponse, err = http.ReadResponse(bufio.NewReader(bytes.NewBuffer(rawRespHdr)), request)
		if err != nil {
			return nil, fmt.Errorf("error while parsing HTTP response: %v", err)
		}

		if hasBody {
			resp.ContentResponse.Body = bodyReader
		} else {
			resp.ContentResponse.Body = emptyReader(0)
		}
	}

	return
}
