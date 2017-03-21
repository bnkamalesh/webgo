package webgo

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const urlchars = `([a-zA-Z0-9_%.]+)`

var l *log.Logger
var validHTTPMethods = []string{http.MethodOptions, http.MethodHead, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (crw *customResponseWriter) WriteHeader(code int) {
	crw.statusCode = code
	crw.ResponseWriter.WriteHeader(code)
}

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

	//ShowLog if enabled will print basic access log to the console
	ShowLog bool

	Handler HandlerChain // Handler function with middlewares
	G       *Globals     // App globals

	//uriKeys is the list of URI params
	uriKeys []string

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

		for i := 0; i < len(r.Pattern); i++ {
			char := string(r.Pattern[i])

			if char == ":" {
				hasKey = true
			} else if hasKey && char != "/" {
				key += char
			} else if hasKey && len(key) > 0 {
				patternString = strings.Replace(patternString, ":"+key, urlchars, 1)
				r.uriKeys = append(r.uriKeys, key)

				hasKey = false
				key = ""
			}
		}

		if hasKey && len(key) > 0 {
			patternString = strings.Replace(patternString, ":"+key, urlchars, 1)
			r.uriKeys = append(r.uriKeys, key)
		}

	}

	patternString = "^" + patternString + "/??$"
	println("\npatternString:", patternString)
	//compile the regex for the pattern string calculated
	reg, err := regexp.Compile(patternString)
	if err != nil {
		return err
	}
	r.uriPattern = reg

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
	Handlers map[string][]Route
}

func (rtr *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	for _, route := range rtr.Handlers[req.Method] {
		startTime := time.Now()
		if ok, params := route.matchAndGet(req.RequestURI); ok {
			endTime := time.Now()
			crw := &customResponseWriter{rw, http.StatusOK}

			// using request context
			// ctx := context.WithValue(req.Context(), "params", params)
			// ctx = context.WithValue(req.Context(), "routeHandler", &route)

			l.Println("params:", params, "route.Pattern:", route.Pattern)
			newHandlerChain := StackInject(route.Handler, "params", params)
			newHandlerChain = StackInject(newHandlerChain, "routeHandler", &route)
			newHandlerChain.ServeHTTP(crw, req)

			out := endTime.Format("2006-01-02 15:04:05 -0700 MST") + " " + req.Method + " " + req.URL.String() + " " + endTime.Sub(startTime).String()
			l.Println(out, crw.statusCode)
			return
		}

	}

	//serve 404
}

// InitRouter initializes Router settings
func InitRouter(routes []Route) *Router {
	var handlers = make(map[string][]Route, len(validHTTPMethods))

	for _, validMethod := range validHTTPMethods {
		handlers[validMethod] = []Route{}
	}

	for _, route := range routes {
		found := false
		for _, validMethod := range validHTTPMethods {
			if route.Method == validMethod {
				found = true
			}
		}

		if found == false {
			log.Fatal("Unsupported HTTP request method provided. Method:", route.Method)
		}

		route.init()
		handlers[route.Method] = append(handlers[route.Method], route)
	}

	return &Router{
		Handlers: handlers,
	}
}

// ===
