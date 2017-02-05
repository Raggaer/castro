package util

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	// Logger main logger instance of the app
	Logger = logrus.New()

	// LoggerOutput output file
	LoggerOutput *os.File

	// LastLoggerDay save last day the logger was created
	LastLoggerDay time.Time
)

// CreateLogFile creates a log file with the current time
func CreateLogFile() (*os.File, time.Time, error) {
	// Get current time
	t := time.Now()

	// Create log file
	f, err := os.OpenFile(filepath.Join(
		"logs",
		fmt.Sprintf("%v-%v-%v.json", t.Year(), t.Month(), t.Day()),
	), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)

	if err != nil {
		return nil, t, err
	}

	return f, t, nil
}

// CreateLogger creates a new logrus logger with the given output
func CreateLogger(out io.Writer) *logrus.Logger {
	// Set logger instance
	l := logrus.New()

	// Set logger output
	l.Out = out

	// Set logger format
	l.Formatter = &logrus.JSONFormatter{}

	// Set fatal handler
	logrus.RegisterExitHandler(func() {

		// Show panic message
		log.Printf(
			"Fatal error encountered. Castro will now exit. For more information check %v",
			filepath.Join("logs", fmt.Sprintf("%v-%v-%v.json", LastLoggerDay.Year(), LastLoggerDay.Month(), LastLoggerDay.Day())),
		)
	})

	return l
}

// RenewLogger runs a routine to check if the logger needs to be renewed
// if true a new logger file is created
func RenewLogger() {
	// Create time ticker
	ticker := time.NewTicker(time.Hour * 24)

	// Stop ticker
	defer ticker.Stop()

	// Wait for the ticker
	for {
		select {
		case <-ticker.C:

			Logger.Info("Creating new log file")

			// Create new log file
			f, day, err := CreateLogFile()

			if err != nil {
				Logger.Fatalf("Cannot renew log file: %v", err)
			}

			// Set las logger day
			LastLoggerDay = day

			// Save olg logger file handle
			old := LoggerOutput

			// Create a new logger
			Logger = CreateLogger(f)

			// Close old logger file handle
			old.Close()

			Logger.Info("Created new log file")
		}
	}
}
