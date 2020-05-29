package icapclient

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// getStatusWithCode prepares the status code and status text from two given strings
func getStatusWithCode(str1, str2 string) (int, string, error) {

	statusCode, err := strconv.Atoi(str1)

	if err != nil {
		return 0, "", err
	}

	status := strings.TrimSpace(str2)

	return statusCode, status, nil
}

// getHeaderVal parses the header and its value from a tcp message string
func getHeaderVal(str string) (string, string) {

	headerVals := strings.SplitN(str, ":", 2)
	header := headerVals[0]
	val := ""

	if len(headerVals) >= 2 {
		val = strings.TrimSpace(headerVals[1])
	}

	return header, val

}

// isRequestLine determines if the tcp message string is a request line, i.e the first line of the message or not
func isRequestLine(str string) bool {
	return strings.Contains(str, ICAPVersion) || strings.Contains(str, HTTPVersion)
}

// setEncapsulatedHeaderValue generates the Encapsulated  values and assigns to the ICAP request string
func setEncapsulatedHeaderValue(icapReqStr, httpReqStr, httpRespStr string) string {
	encpVal := " "

	if strings.HasPrefix(icapReqStr, MethodOPTIONS) { // if the request method is OPTIONS
		if httpReqStr == "" && httpRespStr == "" { // the most common case for OPTIONS method, no Encapsulated body
			encpVal += "null-body=0"
		} else {
			encpVal += "opt-body=0" // if there is an Encapsulated body
		}
	}

	if strings.HasPrefix(icapReqStr, MethodREQMOD) || strings.HasPrefix(icapReqStr, MethodRESPMOD) { // if the request method is RESPMOD or REQMOD
		re := regexp.MustCompile(DoubleCRLF)                // looking for the match of the string \r\n\r\n, as that is the expression that seperates each blocks, i.e headers and bodies
		reqIndices := re.FindAllStringIndex(httpReqStr, -1) // getting the offsets of the matches, tells us the starting/ending point of headers and bodies

		reqEndsAt := 0 // this is needed to calculate the response headers by adding the last offset of the request block
		if reqIndices != nil {
			encpVal += "req-hdr=0"
			reqEndsAt = reqIndices[0][1]
			if len(reqIndices) > 1 { // indicating there is a body present for the request block, as length would have been 1 for a single match of \r\n\r\n
				encpVal += fmt.Sprintf(", req-body=%d", reqIndices[0][1]) // assigning the starting point of the body
				reqEndsAt = reqIndices[1][1]
			} else if httpRespStr == "" {
				encpVal += fmt.Sprintf(", null-body=%d", reqIndices[0][1])
			}
			encpVal += ", "
		}

		respIndices := re.FindAllStringIndex(httpRespStr, -1)

		if respIndices != nil {
			encpVal += fmt.Sprintf("res-hdr=%d", reqEndsAt)
			if len(respIndices) > 1 {
				encpVal += fmt.Sprintf(", res-body=%d", reqEndsAt+respIndices[0][1])
			} else {
				encpVal += fmt.Sprintf(", null-body=%d", reqEndsAt+respIndices[0][1])
			}
		}

	}

	return fmt.Sprintf(icapReqStr, encpVal) // formatting the ICAP request Encapsulated header with the value
}

func addFullBodyInPreviewIndicator(str string) string {
	str = strings.TrimSuffix(str, DoubleCRLF)
	str += fullBodyEndIndicatorPreviewMode
	return str
}

func addHexaBodyByteNotations(str *string) {

	if str == nil {
		return
	}

	ss := strings.SplitN(*str, DoubleCRLF, 2)

	if len(ss) < 2 || ss[1] == "" {
		return
	}

	bodyBytes := []byte(ss[1])

	*str = fmt.Sprintf("%s%s%x%s%s%s", ss[0], DoubleCRLF, len(bodyBytes), CRLF, ss[1], bodyEndIndicator)
}

func chunkBodyInPreviewMode(str *string, pb, cl int, rb []byte) {

	ss := strings.SplitN(*str, DoubleCRLF, 2)

	if len(ss) < 2 || ss[1] == "" {
		return
	}

	bodyStr := ss[1]

	bodyBytes := []byte(bodyStr)

	previewPart := bodyBytes[:pb+1]
	chunkedBodyStr := fmt.Sprintf("%x%s%s%s", pb, CRLF, string(previewPart), CRLF)

	restChunkedBody := chunkBodyByBytes(rb, cl)
	chunkedBodyStr += string(restChunkedBody)

	*str = fmt.Sprintf("%s%s%s%s", ss[0], DoubleCRLF, chunkedBodyStr, CRLF+"0"+CRLF)

}

func chunkBodyByBytes(bdyByte []byte, cl int) []byte {

	newBytes := []byte{}

	spew.Dump("body byte", string(bdyByte))
	for i := 0; i < len(bdyByte); i += cl {
		end := i + cl
		if end > len(bdyByte) {
			end = len(bdyByte)
		}

		newBytes = append(newBytes, []byte(fmt.Sprintf("%x\r\n", len(bdyByte[i:end]))+string(bdyByte[i:end]))...)
	}

	newBytes = append(newBytes, []byte(bodyEndIndicator)...)

	return newBytes
}
