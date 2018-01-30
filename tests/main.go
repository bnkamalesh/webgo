package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bnkamalesh/webgo"
)

var payload = struct {
	Name   string
	Place  string
	Animal string
	Things string
	Other  string
}{
	Name:   "Hello",
	Place:  "Place",
	Animal: "Animal",
	Things: "Things",
	Other: `Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.
	Why do we use it?
	
	It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout. The point of using Lorem Ipsum is that it has a more-or-less normal distribution of letters, as opposed to using 'Content here, content here', making it look like readable English. Many desktop publishing packages and web page editors now use Lorem Ipsum as their default model text, and a search for 'lorem ipsum' will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident, sometimes on purpose (injected humour and the like).
	`,
}

func dummy(w http.ResponseWriter, r *http.Request) {

	wctx := webgo.Context(r)

	var output string
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		webgo.Log.Println(err)
	} else {
		output = string(b)
	}
	defer r.Body.Close()

	webgo.R200(
		w,
		map[string]string{
			"p1":      wctx.Params["p1"],
			"p2":      wctx.Params["p2"],
			"payload": output,
			"pattern": r.URL.Path,
			"method":  r.Method,
		},
	)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	webgo.R200(w, "Hello world")
}

var l = log.New(os.Stdout, "", 0)

func getRoutes(g *webgo.Globals) []*webgo.Route {
	// var mws webgo.Middlewares

	return []*webgo.Route{
		&webgo.Route{
			Name:    "root",                    // A label for the API/URI, this is not used anywhere.
			Method:  http.MethodGet,            // request type
			Pattern: "/",                       // Pattern for the route
			Handler: []http.HandlerFunc{dummy}, // route handler
			G:       g,
		},
		&webgo.Route{
			Name:    "hw-noparams",                  // A label for the API/URI, this is not used anywhere.
			Method:  http.MethodGet,                 // request type
			Pattern: "/nparams",                     // Pattern for the route
			Handler: []http.HandlerFunc{helloWorld}, // route handler
			G:       g,
		},
		&webgo.Route{
			Name:    "hw-withparams",                // A label for the API/URI, this is not used anywhere.
			Method:  http.MethodGet,                 // request type
			Pattern: "/wparams/:p1/goblin/:p2",      // Pattern for the route
			Handler: []http.HandlerFunc{helloWorld}, // route handler
			G:       g,
		},
		&webgo.Route{
			Name:    "params-get",              // A label for the API/URI, this is not used anywhere.
			Method:  http.MethodGet,            // request type
			Pattern: "/hello/:p1/goblin/:p2",   // Pattern for the route
			Handler: []http.HandlerFunc{dummy}, // route handler
			G:       g,
		},

		&webgo.Route{
			Name:    "params-post-sameuri",     // A label for the API/URI, this is not used anywhere.
			Method:  http.MethodPost,           // request type
			Pattern: "/hello/:p1/goblin/:p2",   // Pattern for the route
			Handler: []http.HandlerFunc{dummy}, // route handler
			G:       g,
		},
		&webgo.Route{
			Name:    "params-put-sameuri",      // A label for the API/URI, this is not used anywhere.
			Method:  http.MethodPut,            // request type
			Pattern: "/hello/:p1/goblin/:p2",   // Pattern for the route
			Handler: []http.HandlerFunc{dummy}, // route handler
			G:       g,
		},
		&webgo.Route{
			Name:    "params-patch-sameuri",    // A label for the API/URI, this is not used anywhere.
			Method:  http.MethodPatch,          // request type
			Pattern: "/hello/:p1/goblin/:p2",   // Pattern for the route
			Handler: []http.HandlerFunc{dummy}, // route handler
			G:       g,
		},
		&webgo.Route{
			Name:    "params-delete-sameuri",   // A label for the API/URI, this is not used anywhere.
			Method:  http.MethodDelete,         // request type
			Pattern: "/hello/:p1/goblin/:p2",   // Pattern for the route
			Handler: []http.HandlerFunc{dummy}, // route handler
			G:       g,
		},
		&webgo.Route{
			Name:    "params-options-sameuri", // A label for the API/URI, this is not used anywhere.
			Method:  http.MethodOptions,       // request type
			Pattern: "/hello/:p1/goblin/:p2",  // Pattern for the route
			// Handler: []http.HandlerFunc{dummy}, // route handler
			Handler: []http.HandlerFunc{dummy}, // route handler
			G:       g,
		},
	}
}

func main() {
	// Load configuration from file
	var cfg webgo.Config
	cfg.Load("config.json")

	cfg.HTTPSPort = "8443"
	cfg.CertFile = "./ssl/server.crt"
	cfg.KeyFile = "./ssl/server.key"

	// Initializing context for the app
	var g webgo.Globals

	g.Init(&cfg, nil)

	// Initializing router with all the required routes
	router := webgo.InitRouter(getRoutes(&g))
	router.HideAccessLog = (cfg.Env == "production")

	webgo.Start(
		&cfg,
		router,
		time.Second*5,
		time.Second*5,
	)
}
