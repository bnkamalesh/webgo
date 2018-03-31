// Package middlewares implements the middleware interface which wraps around the whole request.
// i.e. Starting of request till end of execution of the request (after all the chained handlers
// are processed)
// This package also provides 2 chainable middlewares to handle CORS
package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bnkamalesh/webgo"
)

// responseWriter is a custom HTTP response writer
type responseWriter struct {
	http.ResponseWriter
	code int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.code = code
}

// AccessLog is a middleware which prints access log
func AccessLog(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()
	w := &responseWriter{
		ResponseWriter: rw,
	}
	next(w, req)
	end := time.Now()
	log := end.Format("2006-01-02 15:04:05 -0700 MST") + " " + req.Method + " " + req.URL.String() + " " + end.Sub(start).String()
	fmt.Println(log, w.code)
}

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

// Cors is a basic CORS middleware which can be added to individual handlers
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

// CorsWrap is a single Cors middleware which can be applied to the whole app at once
func CorsWrap(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	rw.Header().Set(headerOrigin, req.Header.Get(headerGetOrigin))
	rw.Header().Set(headerMethods, allowMethods)
	rw.Header().Set(headerCreds, "true")
	rw.Header().Set(headerAllowHeaders, allowHeaders+req.Header.Get(headerReqHeaders))
	if req.Method == http.MethodOptions {
		webgo.SendHeader(rw, 200)
		return
	}

	next(rw, req)
}
