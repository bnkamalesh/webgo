package main

import (
	"net/http"
	"time"

	"github.com/bnkamalesh/webgo/v4"
	"github.com/bnkamalesh/webgo/v4/middleware"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	// WebGo context
	wctx := webgo.Context(r)
	// URI parameters, map[string]string
	params := wctx.Params()
	// route, the webgo.Route which is executing this request
	route := wctx.Route
	webgo.R200(
		w,
		map[string]interface{}{
			"route":   route.Name,
			"params":  params,
			"chained": r.Header.Get("chained"),
		},
	)
}

func invalidJSON(w http.ResponseWriter, r *http.Request) {
	webgo.R200(w, make(chan int))
}

func chain(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("chained", "true")
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func getRoutes() []*webgo.Route {
	return []*webgo.Route{
		&webgo.Route{
			Name:          "root",
			Method:        http.MethodGet,
			Pattern:       "/",
			Handlers:      []http.HandlerFunc{helloHandler},
			TrailingSlash: true,
		},
		&webgo.Route{
			Name:          "matchall",
			Method:        http.MethodGet,
			Pattern:       "/matchall/:wildcard*",
			Handlers:      []http.HandlerFunc{helloWorld},
			TrailingSlash: true,
		},
		&webgo.Route{
			Name:                    "api",
			Method:                  http.MethodGet,
			Pattern:                 "/api/:param",
			Handlers:                []http.HandlerFunc{chain, helloWorld},
			TrailingSlash:           true,
			FallThroughPostResponse: true,
		},
		&webgo.Route{
			Name:          "invalidjson",
			Method:        http.MethodGet,
			Pattern:       "/invalidjson",
			Handlers:      []http.HandlerFunc{invalidJSON},
			TrailingSlash: true,
		},
	}
}

func main() {
	cfg := &webgo.Config{
		Host:         "",
		Port:         "8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	router := webgo.NewRouter(cfg, getRoutes())

	router.UseOnSpecialHandlers(middleware.AccessLog)
	router.Use(middleware.AccessLog)
	router.Use(middleware.CorsWrap())
	webgo.GlobalLoggerConfig(
		nil, nil,
		webgo.LogCfgDisableDebug,
	)
	router.Start()
}
