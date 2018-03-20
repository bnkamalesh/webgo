// Package middlewares implements the middleware interface which wraps around the whole request.
// i.e. Starting of request till end of execution of the request (after all the chained handlers
// are processed)
// This package also provides 2 chainable middlewares to handle CORS
package middlewares
