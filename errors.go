package webgo

import (
	Er "errors"
	"log"
)

type Errors struct {
	// Errors in JSON response is a array of string containing all error messages
	Errs []string

	// This is used to add custom error types
	AppErr map[string]error

	// C for `Code`
	C001 error
	C002 error
	C003 error
	C004 error
	C005 error
	C006 error
}

func (e *Errors) Init(errTypes map[string]error) {

	// Error codes which are to be used through out the app.
	// App configuration errors
	e.C001 = Er.New("Invalid number of arguments provided")
	e.C002 = Er.New("Could not unmarshal JSON config file")
	e.C003 = Er.New("App environment not provided in config file. Accepted values are `production` or `development`")
	e.C004 = Er.New("App port not provided in config file")
	e.C005 = Er.New("Invalid JSON")

	e.AppErr = errTypes
}

/*
 Function to log error to the console
 "location": file from which the error is being logged
 "fname": function from which the error is being logged
 "err": the error generated
 "info": A simplified error message
*/

// Error logging in any module should use this log function for consistency
func (e *Errors) Log(location, fname string, err error) {
	log.Println("Error: ", location, " -> ", fname, "\n  -", err)
}

// ===

// Fatal errors which will exit the app after printig on console
func (e *Errors) Fatal(location, fname string, err error) {
	log.Fatal("Fatal Error: ", location, " -> ", fname, "\n  -", err)
}

// ===

// Global variable to access Error logging structure.
var Err Errors
