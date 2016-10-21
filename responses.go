package webgo

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

//ErrorData used to render the error page
type ErrorData struct {
	ErrCode        int
	ErrDescription string
}

type dOutput struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
}

type errOutput struct {
	Errors interface{} `json:"errors"`
	Status int         `json:"status"`
}

const (
	//HeaderContentType is the key for mentioning the response header content type
	HeaderContentType = "Content-Type"
	//JSONContentType is the MIME type when the response is JSON
	JSONContentType = "application/json"
	//HTMLContentType is the MIME type when the response is HTML
	HTMLContentType = "text/html; charset=UTF-8"

	//ErrInternalServer to send when there's an internal server error
	ErrInternalServer = "Internal server error."
)

//SendResponse is used to respond to any request (JSON response) based on the code, data etc.
func SendResponse(w http.ResponseWriter, data interface{}, rCode int) {
	w.Header().Set(HeaderContentType, JSONContentType)

	w.WriteHeader(rCode)

	// Encode data to json and send response
	if err := json.NewEncoder(w).Encode(&dOutput{data, rCode}); err != nil {
		/*
			In case of encoding error, send "internal server error" after
			logging the actual error
		*/
		Log.Println(err)
		R500(w, struct {
			errors []string
		}{
			[]string{ErrInternalServer},
		})
	}
}

//SendError is used to respond to any request with an error
func SendError(w http.ResponseWriter, data interface{}, rCode int) {
	w.Header().Set(HeaderContentType, JSONContentType)

	w.WriteHeader(rCode)

	if err := json.NewEncoder(w).Encode(&errOutput{data, rCode}); err != nil {
		/*
			In case of encoding error, send "internal server error" after
			logging the actual error
		*/
		Log.Println(err)
		R500(w, struct {
			errors []string
		}{
			[]string{ErrInternalServer},
		})
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
	w.WriteHeader(204)
	fmt.Fprint(w)
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
