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

var (
	errLogger  = log.New(os.Stderr, "Error ", log.LstdFlags|log.Lshortfile)
	stdLogger  = log.New(os.Stdout, "", log.LstdFlags)
	infoLogger = log.New(os.Stdout, "Info ", log.LstdFlags)
	warnLogger = log.New(os.Stdout, "Warning ", log.LstdFlags)
)
