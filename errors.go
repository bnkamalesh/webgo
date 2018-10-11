package webgo

import (
	"errors"
	"log"
	"os"
)

var (
	// ErrInvalidPort is the error returned when the port number provided in the config file is invalid
	ErrInvalidPort = errors.New("Port number not provided or is invalid (should be between 0 - 65535)")
)

// Service defines all the logging methods to be implemented
type Logger interface {
	Debug(data ...interface{})
	Info(data ...interface{})
	Warn(data ...interface{})
	Error(data ...interface{})
	Fatal(data ...interface{})
}

// logHandler has all the log writer handlers
type logHandler struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
	fatal *log.Logger
}

// Debug prints log of severity 5
func (lh *logHandler) Debug(data ...interface{}) {
	lh.debug.Println(data...)
}

// Info prints logs of severity 4
func (lh *logHandler) Info(data ...interface{}) {
	lh.info.Println(data...)
}

// Warn prints log of severity 3
func (lh *logHandler) Warn(data ...interface{}) {
	lh.warn.Println(data...)
}

//  Error prints log of severity 2
func (lh *logHandler) Error(data ...interface{}) {
	lh.err.Println(data...)
}

// Fatal prints log of severity 1
func (lh *logHandler) Fatal(data ...interface{}) {
	lh.fatal.Fatalln(data...)
}

var LOGHANDLER Logger

func init() {
	LOGHANDLER = &logHandler{
		fatal: log.New(os.Stderr, "Fatal ", log.LstdFlags|log.Llongfile),
		err:   log.New(os.Stderr, "Error ", log.LstdFlags),
		warn:  log.New(os.Stderr, "Warning ", log.LstdFlags),
		info:  log.New(os.Stdout, "Info ", log.LstdFlags),
		debug: log.New(os.Stdout, "Debug ", log.LstdFlags),
	}
}
