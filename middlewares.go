package webgo

import (
	"net/http"
)

type Middlewares struct{}

// Need to write more default middlewares, for logging  and others.

// A basic Cors middleware definition.
func (m Middlewares) Cors(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			allowedHeaders := r.Header.Get("Access-Control-Allow-Headers")
			if allowedHeaders != "" {
				allowedHeaders = "," + allowedHeaders
			}
			// Set response appropriate headers required for CORS
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Access-Control-Allow-Headers, Authorization, authKey, AuthKey, AuthToken, Token"+allowedHeaders)
			// Stop here & return if it's a preflighted `OPTIONS` request
		}
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ===
