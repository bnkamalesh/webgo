package webgo

import (
	"context"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

//urlchars is the set of characters which are allowed in a URI param
//Regex prepared based on http://stackoverflow.com/a/4669750/1359163, https://tools.ietf.org/html/rfc3986
//const urlchars = `([a-zA-Z0-9\*\-+._~!$()=&',;:@%]+)`
//Though this allows invalid characters in the URI parameter, it has better performance.
const (
	urlchars            = `([^/]+)`
	urlwildcard         = `(.*)`
	errMultiHeaderWrite = `http: multiple response.WriteHeader calls`
	errMultiWrite       = `http: multiple response.Write calls`
	errDuplicateKey     = `Error: Duplicate URI keys found`
)

type ctxkey string

const wgoCtxKey = ctxkey("webgocontext")

var l *log.Logger
var validHTTPMethods = []string{http.MethodOptions, http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

//customResponseWriter is a custom HTTP response writer
type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

//WriteHeader is the interface implementation to get HTTP response code and add it to the custom response writer
func (crw *customResponseWriter) WriteHeader(code int) {
	if crw.written == false {
		crw.statusCode = code
		crw.ResponseWriter.WriteHeader(code)
		return
	}

	l.Println(errMultiHeaderWrite)
}

//Write is the interface implementation to respond to the HTTP request, but check if a response was already sent
func (crw *customResponseWriter) Write(body []byte) (int, error) {
	if crw.written {
		l.Println(errMultiWrite)
		return 0, nil
	}

	crw.written = true
	return crw.ResponseWriter.Write(body)
}

//init initializes the logging variable
func init() {
	l = log.New(os.Stdout, "", 0)
}

// Route struct defines a route for each API
type Route struct {
	//Name is unique identifier for the route
	Name string
	//Method is the HTTP request method/type
	Method string
	//Pattern is the URI pattern to match
	Pattern string

	//HideAccessLog if enabled will not print the basic access log to console
	HideAccessLog bool

	//FallThroughPostResponse if enabled will execute all the handlers even if a response was already sent to the client
	FallThroughPostResponse bool

	//Handler is a slice of http.HandlerFunc which can be middlewares or anything else. Though only 1 of them will be allowed to respond to client.
	//subsquent writes from the following handlers will be ignored
	Handler []http.HandlerFunc
	G       *Globals // App globals

	//uriKeys is the list of URI params
	uriKeys []string

	//uriPatternString is the pattern string which is compiled to regex object
	uriPatternString string
	//uriPattern is the compiled regex to match the URI pattern
	uriPattern *regexp.Regexp
}

//WC is the webgocontext
type WC struct {
	Params map[string]string
	Route  *Route
}

//init intializes prepares the URIKeys, compile regex for the provided pattern
func (r *Route) init() error {
	var patternString = r.Pattern

	if strings.Contains(r.Pattern, ":") {
		//uriValues is a map of URI Key and it's respective value, this is calculated per request
		var key = ""
		var hasKey = false
		var hasWildcard = false

		for i := 0; i < len(r.Pattern); i++ {
			char := string(r.Pattern[i])

			if char == ":" {
				hasKey = true
			} else if char == "*" {
				hasWildcard = true
			} else if hasKey && char != "/" {
				key += char
			} else if hasKey && len(key) > 0 {
				if hasWildcard {
					patternString = strings.Replace(patternString, ":"+key+"*", urlwildcard, 1)
				} else {
					patternString = strings.Replace(patternString, ":"+key, urlchars, 1)
				}

				for idx, k := range r.uriKeys {
					if key == k {
						l.Fatal(errDuplicateKey, "\nURI: ", r.Pattern, "\nKey:", k, ", Position:", idx+1)
					}
				}

				r.uriKeys = append(r.uriKeys, key)

				hasWildcard, hasKey = false, false
				key = ""
			}
		}

		if hasKey && len(key) > 0 {
			if hasWildcard {
				patternString = strings.Replace(patternString, ":"+key+"*", urlwildcard, 1)
			} else {
				patternString = strings.Replace(patternString, ":"+key, urlchars, 1)
			}

			for idx, k := range r.uriKeys {
				if key == k {
					l.Fatal(errDuplicateKey, "\nURI: ", r.Pattern, "\nKey:", k, ", Position:", idx+1)
				}
			}
			r.uriKeys = append(r.uriKeys, key)
		}

	}

	patternString = "^" + patternString + "$"

	//compile the regex for the pattern string calculated
	reg, err := regexp.Compile(patternString)
	if err != nil {
		return err
	}

	r.uriPattern = reg
	r.uriPatternString = patternString
	return nil
}

//matchAndGet will match the given requestURI with its pattern and set its URI params accordingly
func (r *Route) matchAndGet(requestURI string) (bool, map[string]string) {
	if ok := r.uriPattern.Match([]byte(requestURI)); !ok {
		return false, nil
	}

	// Getting URI parameters
	values := r.uriPattern.FindStringSubmatch(requestURI)

	var uriValues map[string]string
	if len(values) > 0 {
		uriValues = make(map[string]string, len(values)-1)
		for j := 1; j < len(values); j++ {
			uriValues[r.uriKeys[j-1]] = values[j]
		}
	}

	return true, uriValues
}

//Router is the HTTP router
type Router struct {
	handlers map[string][]*Route

	optHandlers    []*Route
	headHandlers   []*Route
	getHandlers    []*Route
	postHandlers   []*Route
	putHandlers    []*Route
	patchHandlers  []*Route
	deleteHandlers []*Route

	HideAccessLog bool
	NotFound      http.HandlerFunc
}

func (rtr *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	startTime := time.Now()

	crw := &customResponseWriter{
		ResponseWriter: rw,
	}

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

	var params map[string]string

	ok := false

	path := req.URL.EscapedPath()
	for _, route := range handlers {
		if ok, params = route.matchAndGet(path); !ok {
			continue
		}

		//webgo context object created and is injected to the request context
		reqwc := req.WithContext(
			context.WithValue(
				req.Context(),
				wgoCtxKey,
				&WC{
					Params: params,
					Route:  route,
				},
			),
		)

		for _, handler := range route.Handler {
			if crw.written == false {
				// If there has been no write to response writer yet
				handler(crw, reqwc)
			} else if route.FallThroughPostResponse {
				//run a handler post response write, only if fall through is enabled
				handler(crw, reqwc)
			} else {
				//Do not run any more handlers if already responded and no fall through enabled
				break
			}
		}

		if rtr.HideAccessLog == false && route.HideAccessLog == false {
			endTime := time.Now()
			l.Println(
				endTime.Format("2006-01-02 15:04:05 -0700 MST")+" "+req.Method+" "+req.RequestURI+" "+endTime.Sub(startTime).String(),
				crw.statusCode,
			)
		}

		return
	}

	//serve 404 when there are no matching routes
	rtr.NotFound(rw, req)
	if rtr.HideAccessLog == false {
		endTime := time.Now()
		l.Println(
			endTime.Format("2006-01-02 15:04:05 -0700 MST")+" "+req.Method+" "+req.URL.String()+" "+endTime.Sub(startTime).String(),
			http.StatusNotFound,
		)
	}
}

//Context returns the WebgoContext saved inside the HTTP request context
func Context(r *http.Request) *WC {
	return r.Context().Value(wgoCtxKey).(*WC)
}

// InitRouter initializes Router settings
func InitRouter(routes []*Route) *Router {
	var handlers = make(map[string][]*Route, len(validHTTPMethods))

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

		if found == false {
			l.Fatal("Unsupported HTTP request method provided. Method:", route.Method)
		}

		if route.Handler == nil || len(route.Handler) == 0 {
			l.Fatal("No handlers provided for the route '", route.Pattern, "', method '", route.Method, "'")
		}

		err := route.init()
		if err != nil {
			l.Fatal("Unsupported URI pattern.", route.Pattern, err)
		}

		//checking if the URI pattern is duplicated
		for i := 0; i < idx; i++ {
			rt := routes[i]

			if rt.Name == route.Name {
				l.Println("Warning: Duplicate route name(\"" + rt.Name + "\") detected. Route name should be unique.")
			}

			if rt.Method == route.Method {
				// regex pattern match
				if ok, _ := rt.matchAndGet(route.Pattern); ok {
					l.Println("Warning: Duplicate URI pattern detected.\nPattern: '" + rt.Pattern + "'\nDuplicate pattern: '" + route.Pattern + "'")
					l.Println("Note: Only the first route to match the URI pattern would handle the request")
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
	}
}
