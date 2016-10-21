package webgo

import (
	"log"
	"os"
)

const (
	//C001 Error Code 1
	C001 = "Invalid number of arguments provided"
	//C002 Error Code 2
	C002 = "Could not unmarshal JSON config file"
	//C003 Error Code 3
	C003 = "App environment not provided in config fil Accepted values are `production` or `development`"
	//C004 Error Code 4
	C004 = "App port not provided in config file"
	//C005 Error Code 5
	C005 = "Invalid JSON"
)

//Errors struct is the custom error for webgo error handling
type Errors struct {
	msg string
}

func (e *Errors) Error() string {
	return e.msg
}

//New returns a new instance of Errors struct
func New(str string) *Errors {
	return &Errors{
		msg: str,
	}
}

func init() {
	// Setting up Go log with custom flags
	Log = log.New(os.Stderr, "", log.LstdFlags|log.Llongfile)
}

// Log is used to log errors, which will print the filename and linenumber
var Log *log.Logger
