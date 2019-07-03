package webgo

import (
	"context"
	"net/http"
)

// urlchars is regex to validate characters in a URI parameter
// const urlchars = `([a-zA-Z0-9\*\-+._~!$()=&',;:@%]+)`
// Regex prepared based on http://stackoverflow.com/a/4669750/1359163,
// https://tools.ietf.org/html/rfc3986
// Though the current one allows invalid characters in the URI parameter, it has better performance.
const (
	urlchars            = `([^/]+)`
	urlwildcard         = `(.*)`
	trailingSlash       = `[\/]?`
	errMultiHeaderWrite = `http: multiple response.WriteHeader calls`
	errMultiWrite       = `http: multiple response.Write calls`
	errDuplicateKey     = `Error: Duplicate URI keys found`
)

var validHTTPMethods = []string{
	http.MethodOptions,
	http.MethodHead,
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
}

// customResponseWriter is a custom HTTP response writer
type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// WriteHeader is the interface implementation to get HTTP response code and add
// it to the custom response writer
func (crw *customResponseWriter) WriteHeader(code int) {
	if crw.written {
		LOGHANDLER.Warn(errMultiHeaderWrite)
		return
	}

	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

// Write is the interface implementation to respond to the HTTP request,
// but check if a response was already sent.
func (crw *customResponseWriter) Write(body []byte) (int, error) {
	if crw.written {
		LOGHANDLER.Warn(errMultiWrite)
		return 0, nil
	}

	crw.written = true
	return crw.ResponseWriter.Write(body)
}

// Router is the HTTP router
type Router struct {
	optHandlers    []*Route
	headHandlers   []*Route
	getHandlers    []*Route
	postHandlers   []*Route
	putHandlers    []*Route
	patchHandlers  []*Route
	deleteHandlers []*Route
	allHandlers    map[string][]*Route

	// NotFound is the generic handler for 404 resource not found response
	NotFound http.HandlerFunc
	// AppContext holds all the app specific context which is to be injected into all HTTP
	// request context
	AppContext map[string]interface{}

	// config has all the app config
	config *Config

	// httpServer is the server handler for the active HTTP server
	httpServer *http.Server
	// httpsServer is the server handler for the active HTTPS server
	httpsServer *http.Server
}

func (rtr *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rr := []*Route(nil)

	switch req.Method {
	case http.MethodOptions:
		rr = rtr.optHandlers
	case http.MethodHead:
		rr = rtr.headHandlers
	case http.MethodGet:
		rr = rtr.getHandlers
	case http.MethodPost:
		rr = rtr.postHandlers
	case http.MethodPut:
		rr = rtr.putHandlers
	case http.MethodPatch:
		rr = rtr.patchHandlers
	case http.MethodDelete:
		rr = rtr.deleteHandlers
	default:
		Send(rw, "", "501 Not Implemented", http.StatusNotImplemented)
		return
	}

	route := new(Route)
	ok := false
	params := make(map[string]string)
	path := req.URL.EscapedPath()
	for _, r := range rr {
		if ok, params = r.matchAndGet(path); ok {
			route = r
			break
		}
	}

	if !ok {
		// serve 404 when there are no matching routes
		rtr.NotFound(rw, req)
		return
	}

	crw := &customResponseWriter{
		ResponseWriter: rw,
	}
	// webgo context object created and is injected to the request context
	reqwc := req.WithContext(
		context.WithValue(
			req.Context(),
			wgoCtxKey,
			&WC{
				Params:     params,
				Route:      route,
				AppContext: rtr.AppContext,
			},
		),
	)

	route.serve(crw, reqwc)
}

// Use adds a middleware layer
func (rtr *Router) Use(f func(http.ResponseWriter, *http.Request, http.HandlerFunc)) {
	for _, handlers := range rtr.allHandlers {
		for _, route := range handlers {
			srv := route.serve
			route.serve = func(rw http.ResponseWriter, req *http.Request) {
				f(rw, req, srv)
			}
		}
	}
}

// NewRouter initializes returns a new router instance with all the configurations and routes set
func NewRouter(cfg *Config, routes []*Route) *Router {
	handlers := httpHandlers(routes)
	r := &Router{
		optHandlers:    handlers[http.MethodOptions],
		headHandlers:   handlers[http.MethodHead],
		getHandlers:    handlers[http.MethodGet],
		postHandlers:   handlers[http.MethodPost],
		putHandlers:    handlers[http.MethodPut],
		patchHandlers:  handlers[http.MethodPatch],
		deleteHandlers: handlers[http.MethodDelete],
		allHandlers:    handlers,

		NotFound:   http.NotFound,
		AppContext: make(map[string]interface{}),
		config:     cfg,
	}

	return r
}

// checkDuplicateRoutes checks if any of the routes have duplicate name or URI pattern
func checkDuplicateRoutes(idx int, route *Route, routes []*Route) {
	// checking if the URI pattern is duplicated
	for i := 0; i < idx; i++ {
		rt := routes[i]

		if rt.Name == route.Name {
			LOGHANDLER.Warn("Duplicate route name(\"" + rt.Name + "\") detected. Route name should be unique.")
		}

		if rt.Method == route.Method {
			// regex pattern match
			if ok, _ := rt.matchAndGet(route.Pattern); ok {
				LOGHANDLER.Warn("Duplicate URI pattern detected.\nPattern: '" + rt.Pattern + "'\nDuplicate pattern: '" + route.Pattern + "'")
				LOGHANDLER.Info("Only the first route to match the URI pattern would handle the request")
			}
		}
	}
}

// httpHandlers returns all the handlers in a map, for each HTTP method
func httpHandlers(routes []*Route) map[string][]*Route {
	handlers := make(map[string][]*Route, len(validHTTPMethods))

	for _, validMethod := range validHTTPMethods {
		handlers[validMethod] = []*Route{}
	}

	for idx, route := range routes {
		found := false
		for _, validMethod := range validHTTPMethods {
			if route.Method == validMethod {
				found = true
				break
			}
		}

		if !found {
			LOGHANDLER.Fatal("Unsupported HTTP request method provided. Method:", route.Method)
			return nil
		}

		if route.Handlers == nil || len(route.Handlers) == 0 {
			LOGHANDLER.Fatal("No handlers provided for the route '", route.Pattern, "', method '", route.Method, "'")
			return nil
		}

		err := route.init()
		if err != nil {
			LOGHANDLER.Fatal("Unsupported URI pattern.", route.Pattern, err)
			return nil
		}

		checkDuplicateRoutes(idx, route, routes)

		handlers[route.Method] = append(handlers[route.Method], route)
	}

	return handlers
}
