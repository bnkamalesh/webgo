package webgo

import (
	"net/http"
)

// Middlewares has all the default middlewares provided by webgo
type Middlewares struct{}

const (
	headerOrigin       = "Access-Control-Allow-Origin"
	headerMethods      = "Access-Control-Allow-Methods"
	headerCreds        = "Access-Control-Allow-Credentials"
	headerAllowHeaders = "Access-Control-Allow-Headers"
	headerReqHeaders   = "Access-Control-Request-Headers"
	headerGetOrigin    = "Origin"
	allowMethods       = "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS"
	allowHeaders       = "Accept,Content-Type,Content-Length,Accept-Encoding,Access-Control-Request-Headers,"
)

// Cors is a basic Cors middleware definition.
func (m *Middlewares) Cors(rw http.ResponseWriter, req *http.Request) {
	// Set response appropriate headers required for CORS
	rw.Header().Set(headerOrigin, req.Header.Get(headerGetOrigin))
	rw.Header().Set(headerMethods, allowMethods)
	rw.Header().Set(headerCreds, "true")

	// If there are any extra keys to be added to the header
	rw.Header().Set(headerAllowHeaders, allowHeaders+req.Header.Get(headerReqHeaders))
}

// CorsOptions is a cors middleware just for Options request - adding this helped remove the request method check (an `if` block to check the request type) from Cors middleware
func (m *Middlewares) CorsOptions(rw http.ResponseWriter, req *http.Request) {
	// Set response appropriate headers required for CORS
	rw.Header().Set(headerOrigin, req.Header.Get(headerGetOrigin))
	rw.Header().Set(headerMethods, allowMethods)
	rw.Header().Set(headerCreds, "true")
	rw.Header().Set(headerAllowHeaders, allowHeaders+req.Header.Get(headerReqHeaders))
	SendHeader(rw, 200)
}
