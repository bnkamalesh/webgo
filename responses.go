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

// SendReponse is used to respond to any request (JSON response) based on the code, data etc.
func SendResponse(w http.ResponseWriter, data interface{}, rCode int) {
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

	if rCode < 400 {
		// wrap the response data inside {"data": <data>}
		data = &dOutput{data, rCode}
	} else {
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

// Render is used for rendering templates (HTML)
func Render(w http.ResponseWriter, data interface{}, rCode int, tpl *template.Template) {
	w.WriteHeader(rCode)
	// In case of HTML response, setting appropriate header type for text/HTML response
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
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

// For JSON response with appropriate response code

// R200 - Successful/OK response
func R200(w http.ResponseWriter, data interface{}) {
	SendResponse(w, data, 200)
}

// R201 - New item created
func R201(w http.ResponseWriter, data interface{}) {
	SendResponse(w, data, 201)
}

// R302 - Temporary redirect
func R302(w http.ResponseWriter, data interface{}) {
	SendResponse(w, data, 302)
}

// R400 - Invalid request, any incorrect/erraneous value in the request body
func R400(w http.ResponseWriter, data interface{}) {
	SendResponse(w, data, 400)
}

// R403 - Unauthorized access
func R403(w http.ResponseWriter, data interface{}) {
	SendResponse(w, data, 403)
}

// R404 - Resource not found
func R404(w http.ResponseWriter, data interface{}) {
	SendResponse(w, data, 404)
}

// R406 - Unacceptable header. For any error related to values set in header
func R406(w http.ResponseWriter, data interface{}) {
	SendResponse(w, data, 406)
}

// R451 - Resource taken down because of a legal request
func R451(w http.ResponseWriter, data interface{}) {
	SendResponse(w, data, 451)
}

// R500 - Internal server error
func R500(w http.ResponseWriter, data interface{}) {
	SendResponse(w, data, 500)
}
