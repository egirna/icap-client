package icapclient

// the debug mode determiner
var (
	DEBUG = false
)

// SetDebugMode sets the debug mode for the entire package depending on the bool
func SetDebugMode(debug bool) {
	DEBUG = debug
}
