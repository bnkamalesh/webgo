package webgo

import (
	Er "errors"
	"log"
	"os"
)

type Errors struct {
	// C for `Code`
	C001 error
	C002 error
	C003 error
	C004 error
	C005 error
	C006 error

	// Log is used to log errors, which will print the filename and linenumber
	Log    *log.Logger
	AppErr map[string]error
}

func (e *Errors) init(errTypes map[string]error) {
	// Error codes which are to be used through out the app.
	e.C001 = Er.New("Invalid number of arguments provided")
	e.C002 = Er.New("Could not unmarshal JSON config file")
	e.C003 = Er.New("App environment not provided in config file. Accepted values are `production` or `development`")
	e.C004 = Er.New("App port not provided in config file")
	e.C005 = Er.New("Invalid JSON")

	// Setting up Go log with custom flags
	Err.Log = log.New(os.Stderr, "", log.LstdFlags|log.Llongfile)

	// App configuration errors
	e.AppErr = errTypes
}

func init() {
	// Initializing Err variable with default values
	Err.init(nil)
}

// Global variable to access Error logging structure.
var Err Errors
