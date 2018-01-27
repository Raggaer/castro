package util

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
)

var (
	// Logger main application logging entry point
	Logger = &ApplicationLogger{}
)

// ApplicationLogger struct for application logging to files
type ApplicationLogger struct {
	rw sync.RWMutex

	// Logger main logger instance of the app
	Logger *logrus.Logger

	// LoggerOutput output file
	LoggerOutput *os.File

	// LastLoggerDay save last day the logger was created
	LastLoggerDay time.Time
}

// Custom logrus formatter interface
type castroFormatter struct {
}

// Format converts a logrus text into a valid byte array for castro logging
func (c *castroFormatter) Format(e *logrus.Entry) ([]byte, error) {
	buff := &bytes.Buffer{}
	buff.WriteString(
		fmt.Sprintf("[%s] (%s) %s \r\n", e.Level, e.Time.Format("2006-01-02 15:04:05"), e.Message),
	)
	return buff.Bytes(), nil
}

// CreateLogFile creates a log file with the current time
func CreateLogFile() (*os.File, time.Time, error) {
	// Get current time
	t := time.Now()

	// Create logs folder if needed
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		if err := os.Mkdir("logs", os.ModeDir); err != nil {
			return nil, t, err
		}
	}

	// Create log file
	f, err := os.OpenFile(filepath.Join(
		"logs",
		fmt.Sprintf("%v-%v-%v.txt", t.Year(), t.Month(), t.Day()),
	), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)

	if err != nil {
		return nil, t, err
	}

	// Get file stat
	loggerFileInfo, err := f.Stat()
	if err != nil {
		return nil, t, err
	}

	return f, loggerFileInfo.ModTime(), nil
}

// CreateLogger creates a new logrus logger with the given output
func CreateLogger(out io.Writer) *logrus.Logger {
	// Set logger instance
	l := logrus.New()

	// Set logger output
	l.Out = out

	// Set logger format
	l.Formatter = &castroFormatter{}

	// Set fatal handler
	logrus.RegisterExitHandler(func() {

		// Show panic message
		log.Printf(
			"Fatal error encountered. For more information check %v",
			filepath.Join("logs", fmt.Sprintf("%v-%v-%v.txt", Logger.LastLoggerDay.Year(), Logger.LastLoggerDay.Month(), Logger.LastLoggerDay.Day())),
		)
	})

	return l
}

// RenewLogger runs a routine to check if the logger needs to be renewed if true a new logger file is created
func RenewLogger() {
	// Create time ticker
	ticker := time.NewTicker(time.Second)

	// Stop ticker
	defer ticker.Stop()

	// Wait for the ticker
	for {
		select {
		case <-ticker.C:

			// Check if file is outdated
			if time.Date(
				Logger.LastLoggerDay.Year(),
				Logger.LastLoggerDay.Month(),
				Logger.LastLoggerDay.Day(),
				0,
				0,
				0,
				0,
				Logger.LastLoggerDay.Location(),
			).Add(time.Hour * 24).After(time.Now()) {
				continue
			}

			// Lock mutex
			Logger.rw.Lock()
			defer Logger.rw.Unlock()

			Logger.Logger.Info("Creating new log file")

			// Create new log file
			f, day, err := CreateLogFile()

			if err != nil {
				Logger.Logger.Fatalf("Cannot renew log file: %v", err)
			}

			// Save old logger file handle
			old := Logger.LoggerOutput

			// Create logger
			Logger = &ApplicationLogger{
				Logger:        CreateLogger(f),
				LastLoggerDay: day,
				LoggerOutput:  f,
			}

			// Close old logger file handle
			old.Close()

			Logger.Logger.Info("Created new log file")
		}
	}
}
