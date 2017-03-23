package webgo

import (
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
const urlchars = `([^/]+)`
const urlwildcard = `(.+)`

var l *log.Logger
var validHTTPMethods = []string{http.MethodOptions, http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

//customResponseWriter is a custom HTTP response writer
type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

//WriteHeader is the interface implementation to get HTTP response code and add it to the custom response writer
func (crw *customResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
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

	Handler HandlerChain // Handler function with middlewares
	G       *Globals     // App globals

	//uriKeys is the list of URI params
	uriKeys []string

	//uriPatternString is the pattern string which is compiled to regex object
	uriPatternString string
	//uriPattern is the compiled regex to match the URI pattern
	uriPattern *regexp.Regexp
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
	values := r.uriPattern.FindStringSubmatch(requestURI)
	if len(values) > 0 {
		var uriValues = make(map[string]string, len(values)-1)
		for j := 1; j < len(values); j++ {
			uriValues[r.uriKeys[j-1]] = values[j]
		}
		return true, uriValues
	}

	return r.uriPattern.Match([]byte(requestURI)), map[string]string{}
}

//Router is the HTTP router
type Router struct {
	handlers      map[string][]*Route
	HideAccessLog bool
	NotFound      http.HandlerFunc
}

func (rtr *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	startTime := time.Now()

	for _, route := range rtr.handlers[req.Method] {
		if ok, params := route.matchAndGet(req.RequestURI); ok {

			crw := &customResponseWriter{rw, http.StatusOK}

			//injecting URI parameters and the route handler itself to the context
			newHandlerChain := StackInject(route.Handler, "params", params)
			newHandlerChain = StackInject(newHandlerChain, "routeHandler", route)

			newHandlerChain.ServeHTTP(crw, req)

			if rtr.HideAccessLog == false && route.HideAccessLog == false {
				endTime := time.Now()

				l.Println(
					endTime.Format("2006-01-02 15:04:05 -0700 MST")+" "+req.Method+" "+req.URL.String()+" "+endTime.Sub(startTime).String(),
					crw.statusCode,
				)
			}
			return
		}

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
					l.Fatal("Duplicate URI pattern detected.\nPattern:", rt.Pattern, "\nDuplicate pattern:", route.Pattern)
				}
			}
		}

		handlers[route.Method] = append(handlers[route.Method], route)
	}

	return &Router{
		handlers: handlers,
		NotFound: http.NotFound,
	}
}
