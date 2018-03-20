package webgo

import (
	"errors"
	"log"
	"os"
)

var (
	// ErrInvalidPort is the error returned when the port number provided in the config file is invalid
	ErrInvalidPort = errors.New("Port number not provided or is invalid")
)

var (
	errLogger  = log.New(os.Stderr, "Error ", log.LstdFlags|log.Lshortfile)
	stdLogger  = log.New(os.Stdout, "", log.LstdFlags)
	infoLogger = log.New(os.Stdout, "Info ", log.LstdFlags)
	warnLogger = log.New(os.Stdout, "Warning ", log.LstdFlags)
)

// Errors is the custom error for webgo error handling
type Errors struct {
	msg string
}

func (e *Errors) Error() string {
	return e.msg
}

// New returns a new instance of Errors struct
func New(str string) *Errors {
	return &Errors{
		msg: str,
	}
}
