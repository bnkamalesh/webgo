[![Build Status](https://travis-ci.org/bnkamalesh/webgo.svg?branch=master)](https://travis-ci.org/bnkamalesh/webgo)
[![](https://goreportcard.com/badge/github.com/bnkamalesh/webgo)](https://goreportcard.com/report/github.com/bnkamalesh/webgo)
[![](https://cover.run/go/github.com/bnkamalesh/webgo.svg?tag=golang-1.10)](https://cover.run/go/github.com/bnkamalesh/webgo)
[![](https://godoc.org/github.com/nathany/looper?status.svg)](http://godoc.org/github.com/bnkamalesh/webgo)
[![](https://raw.githubusercontent.com/istio/fortio/master/docs/mentioned-badge.svg?sanitize=true)](https://github.com/avelino/awesome-go)


# WebGo v2.2.3

WebGo is a minimalistic web framework for Go. It gives you the following features.

1. Multiplexer
2. Chaining handlers
3. Middleware
4. Webgo context
5. Helper functions
6. HTTPS ready
7. Graceful shutdown

WebGo's route handlers have the same signature as the standard libraries' HTTP handler.
i.e. `http.HandlerFunc`

### Multiplexer

WebGo multiplexer lets you define URIs with or without parameters. Following are the ways you can
define a URI for a route.

1. `api/users` - no URI parameters
2. `api/users/:userID` - named URI parameter; the value would be available in the key `userID`
3. `api/:wildcard*` - wildcard URIs; every router confirming to `/api/path/to/anything` would be
matched by this route, and the path would be available inside the named URI parameter `wildcard`

If there are multiple routes which have matching path patterns, the first one defined in the list of routes takes precedence and is executed.

### Chaining

Chaining lets you handle a route with multiple handlers, and are executed in sequence.
All handlers in the chain are `http.HandlerFunc`s.

### Middleware

WebGo middleware signature is `func(http.ResponseWriter, *http.Request, http.HandlerFunc)`.
Its `router` exposes a method `Use` to add a middleware to the app. e.g.

```go
func accessLog(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(rw, req)
	end := time.Now()
	log := end.Format("2006-01-02 15:04:05 -0700 MST") + " " + req.Method + " " + req.URL.String() + " " + end.Sub(start).String()
	fmt.Println(log)
}

router := webgo.NewRouter(*webgo.Config, []*Route)
router.Use(accessLog)
```

You can add as many middleware as you want, middleware execution order is in the same order as
they were added. Execution of the middlware is stopped if you do not call `next()`.

### Context

Any app specific context can be set inside the router, and is available inside every request
handler via HTTP request context. The Webgo Context (injected into HTTP request context) also provides the route configuration itself as well as the URI parameters.

```go
router := webgo.NewRouter(*webgo.Config, []*Route)
router.AppContext = map[string]interface{}{
	"dbHandler": <db handler>
}

func API1(rw http.ResponseWriter, req *http.Request) {
	wctx := webgo.Context(req)
	wctx.Params // URI paramters, a map of string of string, where the keys are the named URI parameters
	wctx.Route // is the route configuration itself for which the handler matched
	wctx.AppContext // is the app context configured in `router.AppContext`
}
```

### Helper functions

WebGo has certain helper functions for which the responses are wrapped in a struct according to
error or successful response. All success responses are wrapped in the struct

```go
type dOutput struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
}
```

JSON output looks like
```go
{
	"data": <any data>,
	"status: <HTTP response status code, integer>
}
```

All the error responses are wrapped in the struct

```go
type errOutput struct {
	Errors interface{} `json:"errors"`
	Status int         `json:"status"`
}
```

JSON output for error response looks like

```go
{
	"errors": <any data>,
	"status": <HTTP response status code, integer>
}
```

It provides the following helper functions.

1. SendHeader(w http.ResponseWriter, rCode int) - Send only an HTTP response header with the required
response code.
2. Send(w http.ResponseWriter, contentType string, data interface{}, rCode int) - Send any response
3. SendResponse(w http.ResponseWriter, data interface{}, rCode int) - Send response wrapped in
WebGo's default response struct
4. SendError(w http.ResponseWriter, data interface{}, rCode int) - Send an error wrapped in WebGo's
default error response struct
5. Render(w http.ResponseWriter, data interface{}, rCode int, tpl *template.Template) - Render renders
a Go template, with the requierd data & response code.
6. Render404(w http.ResponseWriter, tpl *template.Template) - Renders a static 404 message
7. R200(w http.ResponseWriter, data interface{}) - Responds with the provided data, with HTTP 200
status code
8. R201(w http.ResponseWriter, data interface{}) - Same as R200, but status code 201
9. R204(w http.ResponseWriter) - Sends a reponse header with status code 201 and no response body.
10. R302(w http.ResponseWriter, data interface{}) - Same as R200, but with status code 302
11. R400(w http.ResponseWriter, data interface{}) - Same as R200, but with status code 400 and
response wrapped in error struct
12. R403(w http.ResponseWriter, data interface{}) - Same as R400, but with status code 403
13. R404(w http.ResponseWriter, data interface{}) - Same as R400, but with status code 404
14. R406(w http.ResponseWriter, data interface{}) - Same as R400, but with status code 406
15. R451(w http.ResponseWriter, data interface{}) - Same as R400, but with status code 451
16. R500(w http.ResponseWriter, data interface{}) - Same as R400, but with status code 500

### HTTPS ready

HTTPS server can be started easily, just by providing the key & cert file required.
You can also have both HTTP & HTTPS servers running side by side.

Start HTTPS server

```go
cfg := webgo.Config{
	Port: "80",
	HTTPSPort: "443",
	CertFile: "/path/to/certfile",
	KeyFile: "/path/to/keyfile",
}
router := webgo.NewRouter(&cfg, []*Route)
router.StartHTTPS()
```

Starting both HTTP & HTTPS server

```go
router := webgo.NewRouter(*webgo.Config, []*Route)
go router.StartHTTPS()
router.Start()
```
### Graceful shutdown

Graceful shutdown lets you shutdown your server without affecting any live connections/clients
connected to your server.

```go
func main() {
	osSig := make(chan os.Signal, 5)

	cfg := webgo.Config{
		Host:            "",
		Port:            "8080",
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    60 * time.Second,
		ShutdownTimeout: 15 * time.Second,
	}
	router := webgo.NewRouter(&cfg, getRoutes())

	go func() {
		<-osSig
		// Initiate HTTP server shutdown
		err := router.Shutdown()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			fmt.Println("shutdown complete")
			os.Exit(0)
		}

		// err := router.ShutdownHTTPS()
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// } else {
		// 	fmt.Println("shutdown complete")
		// 	os.Exit(0)
		// }
	}()

	signal.Notify(osSig, os.Interrupt, syscall.SIGTERM)

	router.Start()

	for {
		// Prevent main thread from exiting, and wait for shutdown to complete
		time.Sleep(time.Second * 1)
	}
}

```

## Full sample

```go
package main

import (
	"net/http"
	"time"

	"github.com/bnkamalesh/webgo/middleware"

	"github.com/bnkamalesh/webgo"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	wctx := webgo.Context(r)
	webgo.R200(
		w,
		wctx.Params, // URI parameters
	)
}

func getRoutes() []*webgo.Route {
	return []*webgo.Route{
		&webgo.Route{
			Name:     "helloworld",                   // A label for the API/URI, this is not used anywhere.
			Method:   http.MethodGet,                 // request type
			Pattern:  "/",                            // Pattern for the route
			Handlers: []http.HandlerFunc{helloWorld}, // route handler
		},
		&webgo.Route{
			Name:     "helloworld",                                     // A label for the API/URI, this is not used anywhere.
			Method:   http.MethodGet,                                   // request type
			Pattern:  "/api/:param",                                    // Pattern for the route
			Handlers: []http.HandlerFunc{middleware.Cors, helloWorld}, // route handler
		},
	}
}

func main() {
	cfg := webgo.Config{
		Host:         "",
		Port:         "8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	router := webgo.NewRouter(&cfg, getRoutes())
	router.Use(middleware.AccessLog)
	router.Start()
}

```

### Testing

URI parameters, route matching, HTTP server, HTTPS server, graceful shutdown and context fetching are tested in the tests provided.

```bash
$ cd $GOPATH/src/github.com/bnkamalesh/webgo
$ go test
```
