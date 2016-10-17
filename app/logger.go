package app

import (
	"os"

	"github.com/op/go-logging"
)

var (
	l          = logging.MustGetLogger("castro")
	utilFormat = logging.MustStringFormatter(
		`%{color}%{time:15:04:05} - %{color:reset}%{message}`,
	)
)

func setLogger() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendUtil := logging.NewBackendFormatter(backend, utilFormat)
	logging.SetBackend(backendUtil)
}
