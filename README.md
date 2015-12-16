#GoLang web framework.

A lightweight & simple web framework for GoLang.

###Requirements

1. `GoLang 1.5.2`, with `vendoring` enabled.

``` 
# Enable vendoring for Go1.5
$ export GO15VENDOREXPERIMENT=1
# ===
```

### Third party libraries used

1. [simplejson](github.com/bitly/go-simplejson), to read `config.json` file.
2. [mgo/mango](http://gopkg.in/mgo.v2), MongoDB driver.
3. [HttpRouter](github.com/julienschmidt/httprouter), multiplexer.
4. [Stack](https://github.com/alexedwards/stack), for chaining request handlers.


### Usage

The default database driver available is `MongoDB`, and its handler can be accessed from
the configuration. There's a function `DummyDb` inside `router.go` which has a sample of how to 
use the database handler.

Any data retrieved using this handler will be in [`bson.M`](https://godoc.org/labix.org/v2/mgo/bson#M) format.

```
package main

import (
	"bitbucket.org/kamaleshbn/webgo"
	"net/http"
)

// This is a sample handler function
func Dummy(ctx *webgo.Context, w http.ResponseWriter, r *http.Request) {
	
	params := ctx.Get("params").(map[string]string) // Get multiplexer params
	webgo.R200(w, params) // Sends a json response in the format `{"data": <params>, "status": 200}`

	// g := ctx.Get("globals").(webgo.Globals) // Get the global settings, and other functions shared
	// webgo.Render404(w, g.Templates["Error"]) // Render an html error page

	// aC := g.App["config"].(map[string]AppConfig) // Get the app specific configuration saved inside the context
	// webgo.Err.AppErr["E1"] // Get custom error set in the app
}
// ===

func main() {

	// Initialize errors, assigns a set of default errors to the struct's fields
	// And add the new ones passed to it
	webgo.Err.Init(map[string]error{
			"E1": "Error 1",
			"E2": "Error 2",
		})

	// Load configuration from file
	var cfg webgo.Config
	cfg.Load("config.json")

	// Load HTML templates
	var t webgo.Templates
	t.Load(map[string]string{
		"Error": cfg.TemplatesBasePath + "/error.html",
		"Test":  cfg.TemplatesBasePath + "/test.html",
	})

	// Initialize database handler
	dbh := webgo.InitDB(cfg.DBConfig.Username,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.Name)

	// Initializing context for the app
	var g webgo.Globals
	g.Init(&cfg, t.Tpls, dbh)
	// Add app specific configuration into globals, AFTER initialising the globals.
	g.App["config"] = "hello world" // this is a map[string]interface{}

	var mws webgo.Middlewares
	// Initializing router with all the required routes
	router := webgo.InitRouter(
		[]webgo.Route{
			webgo.Route{
				"Dummy",                      // A label for the API/URI, this is not used anywhere.
				"GET",                        // request type
				"/:page",                          // Pattern for the route
				webgo.NewStack().Then(Dummy), // route handler
				g,
			},
		})

	webgo.Start(&cfg, router)
	// ====
}

}

```