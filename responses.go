package webgo

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

// ErrorData used to render the error page
type ErrorData struct {
	ErrCode        int
	ErrDescription string
}

// dOutput is the standard/valid output wrapped in `{data: <payload>, status: <http response status>}`
type dOutput struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
}

// errOutput is the error output wrapped in `{errors:<errors>, status: <http response status>}`
type errOutput struct {
	Errors interface{} `json:"errors"`
	Status int         `json:"status"`
}

// responseWriter is a custom HTTP response writer for JSON response
type responseWriter struct {
	http.ResponseWriter
	code int
}

func (rw responseWriter) Write(data []byte) (int, error) {
	rw.WriteHeader(rw.code)
	return rw.ResponseWriter.Write(data)
}

func (rw responseWriter) WriteHeader(code int) {
	rw.ResponseWriter.Header().Set(HeaderContentType, JSONContentType)
	rw.ResponseWriter.WriteHeader(code)
}

const (
	// HeaderContentType is the key for mentioning the response header content type
	HeaderContentType = "Content-Type"
	// JSONContentType is the MIME type when the response is JSON
	JSONContentType = "application/json"
	// HTMLContentType is the MIME type when the response is HTML
	HTMLContentType = "text/html; charset=UTF-8"

	// ErrInternalServer to send when there's an internal server error
	ErrInternalServer = "Internal server error."
)

// SendHeader is used to send only a response header, i.e no response body
func SendHeader(w http.ResponseWriter, rCode int) {
	w.WriteHeader(rCode)
}

// Send sends a completely custom response without wrapping in the
// `{data: <data>, status: <int>` struct
func Send(w http.ResponseWriter, contentType string, data interface{}, rCode int) {
	w.Header().Set(HeaderContentType, contentType)
	w.WriteHeader(rCode)
	_, err := fmt.Fprint(w, data)
	if err != nil {
		R500(w, ErrInternalServer)
	}
}

// SendResponse is used to respond to any request (JSON response) based on the code, data etc.
func SendResponse(w http.ResponseWriter, data interface{}, rCode int) {
	rw := responseWriter{
		ResponseWriter: w,
		code:           rCode,
	}

	err := json.NewEncoder(rw).Encode(dOutput{Data: data, Status: rCode})
	if err != nil {
		/*
			In case of encoding error, send "internal server error" after
			logging the actual error.
		*/
		errLogger.Println(err)
		R500(w, ErrInternalServer)
	}
}

// SendError is used to respond to any request with an error
func SendError(w http.ResponseWriter, data interface{}, rCode int) {
	rw := responseWriter{
		ResponseWriter: w,
		code:           rCode,
	}

	err := json.NewEncoder(rw).Encode(errOutput{data, rCode})
	if err != nil {
		/*
			In case of encoding error, send "internal server error" after
			logging the actual error.
		*/
		errLogger.Println(err)
		R500(w, ErrInternalServer)
	}
}

// Render is used for rendering templates (HTML)
func Render(w http.ResponseWriter, data interface{}, rCode int, tpl *template.Template) {
	// In case of HTML response, setting appropriate header type for text/HTML response
	w.Header().Set(HeaderContentType, HTMLContentType)
	w.WriteHeader(rCode)

	// Rendering an HTML template with appropriate data
	tpl.Execute(w, data)
}

// Render404 - used to render a 404 page
func Render404(w http.ResponseWriter, tpl *template.Template) {
	Render(w, ErrorData{
		404,
		"Sorry, the URL you requested was not found on this server... Or you're lost :-/",
	},
		404,
		tpl,
	)
}

// R200 - Successful/OK response
func R200(w http.ResponseWriter, data interface{}) {
	SendResponse(w, data, 200)
}

// R201 - New item created
func R201(w http.ResponseWriter, data interface{}) {
	SendResponse(w, data, 201)
}

// R204 - empty, no content
func R204(w http.ResponseWriter) {
	SendHeader(w, 204)
}

// R302 - Temporary redirect
func R302(w http.ResponseWriter, data interface{}) {
	SendResponse(w, data, 302)
}

// R400 - Invalid request, any incorrect/erraneous value in the request body
func R400(w http.ResponseWriter, data interface{}) {
	SendError(w, data, 400)
}

// R403 - Unauthorized access
func R403(w http.ResponseWriter, data interface{}) {
	SendError(w, data, 403)
}

// R404 - Resource not found
func R404(w http.ResponseWriter, data interface{}) {
	SendError(w, data, 404)
}

// R406 - Unacceptable header. For any error related to values set in header
func R406(w http.ResponseWriter, data interface{}) {
	SendError(w, data, 406)
}

// R451 - Resource taken down because of a legal request
func R451(w http.ResponseWriter, data interface{}) {
	SendError(w, data, 451)
}

// R500 - Internal server error
func R500(w http.ResponseWriter, data interface{}) {
	SendError(w, data, 500)
}
