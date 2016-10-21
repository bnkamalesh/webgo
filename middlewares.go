package webgo

import (
	"net/http"
)

//Middlewares has all the default middlewares provided by webgo
type Middlewares struct{}

const (
	headerOrigin       = "Access-Control-Allow-Origin"
	headerMethods      = "Access-Control-Allow-Methods"
	headerCreds        = "Access-Control-Allow-Credentials"
	headerHeaders      = "Access-Control-Allow-Headers"
	headerReqHeaders   = "Access-Control-Request-Headers"
	headerGetOrigin    = "Origin"
	allowMethods       = "HEAD,GET,POST,PUT,PATCH,DELETE"
	allowMethodOptions = "OPTIONS"
	allowHeaders       = "Accept,Content-Type,Content-Length,Accept-Encoding,Access-Control-Request-Headers,"
)

//Cors is a basic Cors middleware definition.
func (m *Middlewares) Cors(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set response appropriate headers required for CORS
		w.Header().Set(headerOrigin, r.Header.Get(headerGetOrigin))
		w.Header().Set(headerMethods, allowMethods)
		w.Header().Set(headerCreds, "true")

		// If there are any extra keys to be added to the header
		w.Header().Set(headerHeaders, allowHeaders+r.Header.Get(headerReqHeaders))

		next.ServeHTTP(w, r)
	})
}

//CorsOptions is a cors middleware just for Options request - adding this helped remove the request method check (an `if` block to check the request type) from Cors middleware
func (m *Middlewares) CorsOptions(ctx *Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set response appropriate headers required for CORS
		w.Header().Set(headerOrigin, r.Header.Get(headerGetOrigin))
		w.Header().Set(headerMethods, allowMethodOptions)
		w.Header().Set(headerCreds, "true")
		w.Header().Set(headerHeaders, allowHeaders+r.Header.Get(headerReqHeaders))
		return
	})
}
