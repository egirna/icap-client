package icapclient

import (
	"io"
	"log"
	"os"
)

// the debug mode determiner & the writer to the write the debug output to
var (
	DEBUG       = false
	debugWriter io.Writer
)

// SetDebugMode sets the debug mode for the entire package depending on the bool
func SetDebugMode(debug bool) {
	DEBUG = debug

	if DEBUG { // setting os.Stdout as the default debug writer if debug mode is enabled
		debugWriter = os.Stdout
		log.SetOutput(debugWriter)
	}
}

// SetDebugOutput sets writer to write the debug outputs (default: os.Stdout)
func SetDebugOutput(w io.Writer) {
	debugWriter = w
	log.SetOutput(debugWriter)
}
