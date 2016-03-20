package webgo

import (
	"encoding/json"
	"html/template"
	"net/http"
)

// Struct used to render the error page
type ErrorData struct {
	ErrCode        int
	ErrDescription string
}

//===

// The function to respond to any request (JSON response) based on the code, data etc.
func sendResponse(w http.ResponseWriter, data interface{}, rCode int) {
	w.WriteHeader(rCode)

	// renderFile = false means response should be in JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	type dOutput struct {
		Data   interface{} `json:"data"`
		Status int         `json:"status"`
	}

	type errOutput struct {
		Errors interface{} `json:"errors"`
		Status int         `json:"status"`
	}

	switch rCode {
	case 200, 201:
		// wrap the response data inside {"data": <data>}
		data = &dOutput{data, rCode}
	default:
		// wrap the response data inside {"errors": <data>}
		data = &errOutput{data, rCode}
	}

	// Encode data to json and send response
	if err := json.NewEncoder(w).Encode(data); err != nil {
		/*
			In case of encoding error, send "internal server error" after
			logging the actual error
		*/
		R500(w, struct {
			errors []string
		}{
			[]string{"Internal server error."},
		})
	}
}

// ===

// For rendering templates (HTML)
func Render(w http.ResponseWriter, data interface{}, rCode int, tpl *template.Template) {
	w.WriteHeader(rCode)
	// In case of HTML response, setting appropriate header type for text/HTML response
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	// Rendering an HTML template with appropriate data
	tpl.Execute(w, data)
}

// ===

// For rendering 404
func Render404(w http.ResponseWriter, tpl *template.Template) {
	Render(w, ErrorData{
		404,
		"Sorry, the URL you requested was not found on this server... Or you're lost :-/",
	},
		404,
		tpl,
	)
}

// ===

// For JSON response with appropriate response code

// Successful/OK response
func R200(w http.ResponseWriter, data interface{}) {
	sendResponse(w, data, 200)
}

// New item created
func R201(w http.ResponseWriter, data interface{}) {
	sendResponse(w, data, 201)
}

// ===

// Temporary redirect
func R302(w http.ResponseWriter, data interface{}) {
	sendResponse(w, data, 302)
}

// ===

// Invalid request, any incorrect/erraneous value in the request body
func R400(w http.ResponseWriter, data interface{}) {
	sendResponse(w, data, 400)
}

// ===

// Unauthorized access
func R403(w http.ResponseWriter, data interface{}) {
	sendResponse(w, data, 403)
}

// ===

// Unacceptable header. For any error related to values set in header
func R406(w http.ResponseWriter, data interface{}) {
	sendResponse(w, data, 406)
}

// Resource not found
func R404(w http.ResponseWriter, data interface{}) {
	sendResponse(w, data, 404)
}

// ===

// Internal server error
func R500(w http.ResponseWriter, data interface{}) {
	sendResponse(w, data, 500)
}

// ===
