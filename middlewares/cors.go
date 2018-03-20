package middlewares

import (
	"net/http"

	"github.com/bnkamalesh/webgo"
)

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

// Cors is a basic CORS middleware which can be added to any of the handlers
func Cors(rw http.ResponseWriter, req *http.Request) {
	// Set appropriate response headers required for CORS
	rw.Header().Set(headerOrigin, req.Header.Get(headerGetOrigin))
	rw.Header().Set(headerMethods, allowMethods)
	rw.Header().Set(headerCreds, "true")

	// Adding allowed headers
	rw.Header().Set(headerAllowHeaders, allowHeaders+req.Header.Get(headerReqHeaders))
}

// CorsOptions is a CORS middleware only for OPTIONS request method
func CorsOptions(rw http.ResponseWriter, req *http.Request) {
	// Set appropriate response headers required for CORS
	rw.Header().Set(headerOrigin, req.Header.Get(headerGetOrigin))
	rw.Header().Set(headerMethods, allowMethods)
	rw.Header().Set(headerCreds, "true")
	rw.Header().Set(headerAllowHeaders, allowHeaders+req.Header.Get(headerReqHeaders))
	webgo.SendHeader(rw, 200)
}
