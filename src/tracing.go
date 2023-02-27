package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/trace"
)

// TraceLevel represents a log level that is used for tracing.
const TraceLevel = 5

// logger represents a custom logger that logs messages with different log levels.
var logger = log.New(os.Stderr, "myapp: ", log.Ldate|log.Ltime|log.Lshortfile)

// SetOutput sets the output destination for the logger.
func SetOutput(w io.Writer) {
	logger.SetOutput(w)
}

// SetFlags sets the formatting flags for the logger.
func SetFlags(flags int) {
	logger.SetFlags(flags)
}

// Debug logs a message with debug level.
func Debug(msg string, args ...interface{}) {
	logger.Printf("DEBUG: "+msg, args...)
}

// Info logs a message with info level.
func Info(msg string, args ...interface{}) {
	logger.Printf("INFO: "+msg, args...)
}

// Warn logs a message with warning level.
func Warn(msg string, args ...interface{}) {
	logger.Printf("WARN: "+msg, args...)
}

// Error logs a message with error level.
func Error(msg string, args ...interface{}) {
	logger.Printf("ERROR: "+msg, args...)
}

// Trace logs a message with trace level.
func Trace(msg string, args ...interface{}) {
	logger.Printf("TRACE: "+msg, args...)
}

// FmtSubscriber represents a subscriber that formats log events and writes them to a file.
type FmtSubscriber struct {
	w io.Writer
}

// NewFmtSubscriber creates a new FmtSubscriber with the specified output writer.
func NewFmtSubscriber(w io.Writer) *FmtSubscriber {
	return &FmtSubscriber{w}
}

// Printf formats and writes a log event to the output writer.
func (s *FmtSubscriber) Printf(format string, a ...interface{}) {
	msg := format
	if len(a) > 0 {
		msg = format + " " + fmt.Sprintf("%v", a)
	}
	if trace.IsEnabled() {
		trace.Log(context.Background(), msg, "")
	}
}
