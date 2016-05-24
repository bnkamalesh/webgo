package webgo

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Struct to define a route for each API
type Route struct {
	Name    string       // Just a label to name the route/API, this is not used anywhere
	Method  string       // Request type
	Pattern string       // URI
	Handler HandlerChain // Handler function with middlewares
	G       *Globals     // App globals
}

// Inject httprouter params to the context
func InjectParams(route Route) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// convert httprouter params to map of string
		params := make(map[string]string)
		for i := range ps {
			params[ps[i].Key] = ps[i].Value
		}

		// Injecting multiplexer params to every request context
		newHandlerChain := StackInject(route.Handler, "params", params)
		// Injecting globals to every request context
		route.G.Params["routeName"] = route.Name
		newHandlerChain = StackInject(newHandlerChain, "globals", route.G)
		newHandlerChain.ServeHTTP(w, r)
	}
}

// Initiate Router settings
func InitRouter(routes []Route) *httprouter.Router {
	router := httprouter.New()

	// Handles all the route types
	for _, route := range routes {
		switch route.Method {
		case "OPTIONS":
			router.OPTIONS(
				route.Pattern,
				InjectParams(route))
		case "GET":
			router.GET(
				route.Pattern,
				InjectParams(route))
		case "POST":
			router.POST(
				route.Pattern,
				InjectParams(route))
		case "PUT":
			router.PUT(
				route.Pattern,
				InjectParams(route))
		case "DELETE":
			router.DELETE(
				route.Pattern,
				InjectParams(route))
		case "PATCH":
			router.PATCH(
				route.Pattern,
				InjectParams(route))
		}
	}
	return router
}

// ===
