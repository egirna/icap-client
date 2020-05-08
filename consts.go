package icapclient

// the icap request methods
const (
	MethodOPTIONS = "OPTIONS"
	MethodRESPMOD = "RESPMOD"
	MethodREQMOD  = "REQMOD"
)

// the error messages
const (
	ErrInvalidScheme       = "the url scheme must be icap://"
	ErrMethodNotRegistered = "the requested method is not registered"
	ErrInvalidHost         = "the requested host is invalid"
	ErrConnectionNotOpen   = "no open connection to close"
	ErrInvalidICAPResponse = "invalid response received from icap server"
	ErrInvalidHTTPResponse = "invalid response received from http server"
)

// general constants required for the package
const (
	SchemeICAP  = "icap"
	ICAPVersion = "ICAP/1.0"
	HTTPVersion = "HTTP/1.1"
)
