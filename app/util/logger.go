package util

import (
	"os"

	"github.com/op/go-logging"
)

var (
	// Logger main logger instance of the app
	Logger = logging.MustGetLogger("castro")

	// Format for all util messages
	utilFormat = logging.MustStringFormatter(
		`%{color}%{time:15:04:05} - %{level:.5s}:%{color:reset} %{message}`,
	)
)

func init() {
	// Create backend for logging
	backend := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend we want to add some additional
	// information
	backendFormatter := logging.NewBackendFormatter(backend, utilFormat)

	// Set the backends to be used.
	logging.SetBackend(backendFormatter)
}
