package webgo

import (
	"context"
	"net/http"
	"regexp"
	"strings"
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

type ctxkey string

const wgoCtxKey = ctxkey("webgocontext")

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
		warnLogger.Println(errMultiHeaderWrite)
		return
	}

	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

// Write is the interface implementation to respond to the HTTP request,
// but check if a response was already sent.
func (crw *customResponseWriter) Write(body []byte) (int, error) {
	if crw.written {
		warnLogger.Println(errMultiWrite)
		return 0, nil
	}

	crw.written = true
	return crw.ResponseWriter.Write(body)
}

// Route defines a route for each API
type Route struct {
	// Name is unique identifier for the route
	Name string
	// Method is the HTTP request method/type
	Method string
	// Pattern is the URI pattern to match
	Pattern string
	// TrailingSlash if set to true, the URI will be matched with or without
	// a trailing slash. Note: It does not *do* a redirect.
	TrailingSlash bool

	// FallThroughPostResponse if enabled will execute all the handlers even if a response was already sent to the client
	FallThroughPostResponse bool

	// Handler is a slice of http.HandlerFunc which can be middlewares or anything else. Though only 1 of them will be allowed to respond to client.
	// subsequent writes from the following handlers will be ignored
	Handler []http.HandlerFunc

	// uriKeys is the list of URI parameter variables available for this route
	uriKeys []string
	// uriPatternString is the pattern string which is compiled to regex object
	uriPatternString string
	// uriPattern is the compiled regex to match the URI pattern
	uriPattern *regexp.Regexp
}

// WC is the webgocontext
type WC struct {
	Params     map[string]string
	Route      *Route
	AppContext map[string]interface{}
}

// init prepares the URIKeys, compile regex for the provided pattern
func (r *Route) init() error {
	patternString := r.Pattern

	if strings.Contains(r.Pattern, ":") {
		// uriValues is a map of URI Key and it's respective value,
		// this is calculated per request
		key := ""
		hasKey := false
		hasWildcard := false

		for i := 0; i < len(r.Pattern); i++ {
			char := string(r.Pattern[i])

			if char == ":" {
				hasKey = true
			} else if char == "*" {
				hasWildcard = true
			} else if hasKey && char != "/" {
				key += char
			} else if hasKey && len(key) > 0 {
				regexPattern := ""
				patternKey := ""
				if hasWildcard {
					patternKey = ":" + key + "*"
					regexPattern = urlwildcard
				} else {
					patternKey = ":" + key
					regexPattern = urlchars
				}

				patternString = strings.Replace(patternString, patternKey, regexPattern, 1)

				for idx, k := range r.uriKeys {
					if key == k {
						errLogger.Fatalln(errDuplicateKey, "\nURI: ", r.Pattern, "\nKey:", k, ", Position:", idx+1)
					}
				}

				r.uriKeys = append(r.uriKeys, key)

				hasWildcard, hasKey = false, false
				key = ""
			}
		}

		if hasKey && len(key) > 0 {
			regexPattern := ""
			patternKey := ""
			if hasWildcard {
				patternKey = ":" + key + "*"
				regexPattern = urlwildcard
			} else {
				patternKey = ":" + key
				regexPattern = urlchars
			}

			patternString = strings.Replace(patternString, patternKey, regexPattern, 1)

			for idx, k := range r.uriKeys {
				if key == k {
					errLogger.Fatalln(errDuplicateKey, "\nURI: ", r.Pattern, "\nKey:", k, ", Position:", idx+1)
				}
			}
			r.uriKeys = append(r.uriKeys, key)
		}

	}

	if r.TrailingSlash {
		patternString = "^" + patternString + trailingSlash + "$"
	} else {

		patternString = "^" + patternString + "$"
	}

	// compile the regex for the pattern string calculated
	reg, err := regexp.Compile(patternString)
	if err != nil {
		return err
	}

	r.uriPattern = reg
	r.uriPatternString = patternString
	return nil
}

// matchAndGet will match the given requestURI with its pattern and set its URI params accordingly
func (r *Route) matchAndGet(requestURI string) (bool, map[string]string) {
	if r.Pattern == requestURI {
		return true, nil
	}

	if !r.uriPattern.Match([]byte(requestURI)) {
		return false, nil
	}

	// Getting URI parameters
	values := r.uriPattern.FindStringSubmatch(requestURI)
	if len(values) == 0 {
		return true, nil
	}

	uriValues := make(map[string]string, len(values)-1)
	for i := 1; i < len(values); i++ {
		uriValues[r.uriKeys[i-1]] = values[i]
	}
	return true, uriValues

}

// Router is the HTTP router
type Router struct {
	handlers map[string][]*Route

	optHandlers    []*Route
	headHandlers   []*Route
	getHandlers    []*Route
	postHandlers   []*Route
	putHandlers    []*Route
	patchHandlers  []*Route
	deleteHandlers []*Route

	NotFound http.HandlerFunc

	// config has all the app config
	config     *Config
	appContext map[string]interface{}
}

// ServeHTTP is the required `ServeHTTP` implementation to listen to HTTP requests
func (rtr *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var handlers []*Route

	switch req.Method {
	case http.MethodOptions:
		handlers = rtr.optHandlers
	case http.MethodHead:
		handlers = rtr.headHandlers
	case http.MethodGet:
		handlers = rtr.getHandlers
	case http.MethodPost:
		handlers = rtr.postHandlers
	case http.MethodPut:
		handlers = rtr.putHandlers
	case http.MethodPatch:
		handlers = rtr.patchHandlers
	case http.MethodDelete:
		handlers = rtr.deleteHandlers
	}

	var route *Route
	ok := false
	params := make(map[string]string)
	path := req.URL.EscapedPath()
	for _, h := range handlers {
		if ok, params = h.matchAndGet(path); ok {
			route = h
			break
		}
	}

	if route == nil || len(route.Handler) == 0 {
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
				AppContext: rtr.appContext,
			},
		),
	)

	for _, handler := range route.Handler {
		if crw.written == false {
			// If there has been no write to response writer yet
			handler(crw, reqwc)
		} else if route.FallThroughPostResponse {
			// run a handler post response write, only if fall through is enabled
			handler(crw, reqwc)
		} else {
			// Do not run any more handlers if already responded and no fall through enabled
			break
		}
	}
}

// Context returns the WebgoContext injected inside the HTTP request context
func Context(r *http.Request) *WC {
	wc, _ := r.Context().Value(wgoCtxKey).(*WC)
	return wc
}

// NewRouter initializes returns a new router instance with all the configurations and routes set
func NewRouter(cfg *Config, routes []*Route) *Router {
	handlers := make(map[string][]*Route, len(validHTTPMethods))

	for _, validMethod := range validHTTPMethods {
		handlers[validMethod] = []*Route{}
	}

	for idx, route := range routes {
		found := false
		for _, validMethod := range validHTTPMethods {
			if route.Method == validMethod {
				found = true
			}
		}

		if !found {
			errLogger.Fatalln("Unsupported HTTP request method provided. Method:", route.Method)
		}

		if route.Handler == nil || len(route.Handler) == 0 {
			errLogger.Fatalln("No handlers provided for the route '", route.Pattern, "', method '", route.Method, "'")
		}

		err := route.init()
		if err != nil {
			errLogger.Fatalln("Unsupported URI pattern.", route.Pattern, err)
		}

		// checking if the URI pattern is duplicated
		for i := 0; i < idx; i++ {
			rt := routes[i]

			if rt.Name == route.Name {
				warnLogger.Println("Duplicate route name(\"" + rt.Name + "\") detected. Route name should be unique.")
			}

			if rt.Method == route.Method {
				// regex pattern match
				if ok, _ := rt.matchAndGet(route.Pattern); ok {
					warnLogger.Println("Duplicate URI pattern detected.\nPattern: '" + rt.Pattern + "'\nDuplicate pattern: '" + route.Pattern + "'")
					infoLogger.Println("Only the first route to match the URI pattern would handle the request")
				}
			}
		}

		handlers[route.Method] = append(handlers[route.Method], route)
	}

	return &Router{
		handlers: handlers,

		optHandlers:    handlers[http.MethodOptions],
		headHandlers:   handlers[http.MethodHead],
		getHandlers:    handlers[http.MethodGet],
		postHandlers:   handlers[http.MethodPost],
		putHandlers:    handlers[http.MethodPut],
		patchHandlers:  handlers[http.MethodPatch],
		deleteHandlers: handlers[http.MethodDelete],

		NotFound: http.NotFound,

		config: cfg,
	}
}
