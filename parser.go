package icapclient

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func getStatusWithCode(str1, str2 string) (int, string, error) {

	statusCode, err := strconv.Atoi(str1)

	if err != nil {
		return 0, "", err
	}

	status := strings.TrimSpace(str2)

	return statusCode, status, nil
}

func getHeaderVal(str string) (string, string) {

	headerVals := strings.SplitN(str, ":", 2)
	header := headerVals[0]
	val := ""

	if len(headerVals) >= 2 {
		val = strings.TrimSpace(headerVals[1])
	}

	return header, val

}

func isRequestLine(str string) bool {
	return strings.Contains(str, ICAPVersion) || strings.Contains(str, HTTPVersion)
}

func setEncapsulatedHeaderValue(icapReqStr, httpReqStr, httpRespStr string) string {
	encpVal := " "

	if strings.HasPrefix(icapReqStr, MethodOPTIONS) {
		if httpReqStr == "" && httpRespStr == "" {
			encpVal += "null-body=0"
		}
	}

	if strings.HasPrefix(icapReqStr, MethodREQMOD) || strings.HasPrefix(icapReqStr, MethodRESPMOD) {
		re, _ := regexp.Compile(DoubleCRLF)
		reqIndices := re.FindAllStringIndex(httpReqStr, -1)

		reqEndsAt := 0
		if reqIndices != nil {
			encpVal += "req-hdr=0"
			reqEndsAt = reqIndices[0][1]
			if len(reqIndices) > 1 {
				encpVal += fmt.Sprintf(", req-body=%d", reqIndices[0][1])
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

	return fmt.Sprintf(icapReqStr, encpVal)
}

// func setEncapsulatedHeaderValue(icapReqStr, httpReqStr, httpRespStr string) string {
// 	encpVal := " "
//
// 	re, _ := regexp.Compile("\r\n\r\n")
//
// 	spew.Dump(re.FindAllStringIndex(httpReqStr, -1), re.FindAllStringIndex(httpRespStr, -1))
// 	spew.Dump("req", httpReqStr, "resp", httpRespStr)
//
// 	if iss := strings.Split(icapReqStr, " "); len(iss) > 0 && iss[0] == MethodOPTIONS {
// 		encpVal += "null-body=0"
// 	} else if iss[0] == MethodREQMOD || iss[0] == MethodRESPMOD {
//
// 		// reqBodyLen := 0
// 		reqHeaderLen := 0
// 		if httpReqStr != "" {
// 			encpVal += "req-hdr=0"
// 			hss := strings.Split(httpReqStr, CRLF+CRLF)
// 			reqHeaderLen = len([]byte(hss[0]))
//
// 			if len(hss) > 1 && hss[1] != "" {
// 				reqBody := reqHeaderLen + 1
// 				encpVal += fmt.Sprintf(", req-body=%d", reqBody)
// 				// reqBodyLen = len([]byte(hss[1]))
// 			}
// 		}
//
// 		if httpRespStr != "" {
// 			if encpVal != " " {
// 				encpVal += ", "
// 			}
// 			hss := strings.Split(httpRespStr, CRLF+CRLF)
// 			respHdr := 112 //reqHeaderLen + reqBodyLen  + 1
// 			encpVal += fmt.Sprintf("res-hdr=%d", respHdr)
// 			if len(hss) > 1 && hss[1] != "" {
// 				// respHeaderLen := len([]byte(hss[0]))
// 				encpVal += fmt.Sprintf(", res-body=%d", 417) //respHdr+respHeaderLen+1)
// 			}
// 		}
//
// 	}
//
// 	return fmt.Sprintf(icapReqStr+httpReqStr+httpRespStr, encpVal)
// }
