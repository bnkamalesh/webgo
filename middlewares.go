package webgo

import (
	"net/http"
)

type Middlewares struct{}

// A basic Cors middleware definition.
func (m *Middlewares) Cors(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set response appropriate headers required for CORS
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Access-Control-Request-Headers,"+r.Header.Get("Access-Control-Request-Headers"))

		next.ServeHTTP(w, r)
	})
}

// Cors middleware just for Options request - adding this helped remove the
// request method check from Cors middleware
func (m *Middlewares) CorsOptions(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set response appropriate headers required for CORS
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Access-Control-Request-Headers,"+r.Header.Get("Access-Control-Request-Headers"))
		return
	})
}
