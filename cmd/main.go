package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bnkamalesh/webgo/v6"
	"github.com/bnkamalesh/webgo/v6/extensions/sse"
	"github.com/bnkamalesh/webgo/v6/middleware/accesslog"
	"github.com/bnkamalesh/webgo/v6/middleware/cors"
)

var (
	lastModified = time.Now().Format(http.TimeFormat)
)

func chain(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("chained", "true")
}

// errLogger is a middleware which will log all errors returned/set by a handler
func errLogger(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)

	err := webgo.GetError(r)
	if err != nil {
		// log only server errors
		if webgo.ResponseStatus(w) > 499 {
			log.Println("errorLogger:", err.Error())
		}
	}
}

func routegroupMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Add("routegroup", "true")
	next(w, r)
}

func getRoutes(sse *sse.SSE) []*webgo.Route {
	return []*webgo.Route{
		{
			Name:          "root",
			Method:        http.MethodGet,
			Pattern:       "/",
			Handlers:      []http.HandlerFunc{HomeHandler},
			TrailingSlash: true,
		},
		{
			Name:          "matchall",
			Method:        http.MethodGet,
			Pattern:       "/matchall/:wildcard*",
			Handlers:      []http.HandlerFunc{ParamHandler},
			TrailingSlash: true,
		},
		{
			Name:                    "api",
			Method:                  http.MethodGet,
			Pattern:                 "/api/:param",
			Handlers:                []http.HandlerFunc{chain, ParamHandler},
			TrailingSlash:           true,
			FallThroughPostResponse: true,
		},
		{
			Name:          "invalidjson",
			Method:        http.MethodGet,
			Pattern:       "/invalidjson",
			Handlers:      []http.HandlerFunc{InvalidJSONHandler},
			TrailingSlash: true,
		},
		{
			Name:          "error-setter",
			Method:        http.MethodGet,
			Pattern:       "/error-setter",
			Handlers:      []http.HandlerFunc{ErrorSetterHandler},
			TrailingSlash: true,
		},
		{
			Name:          "original-responsewriter",
			Method:        http.MethodGet,
			Pattern:       "/original-responsewriter",
			Handlers:      []http.HandlerFunc{OriginalResponseWriterHandler},
			TrailingSlash: true,
		},
		{
			Name:          "static",
			Method:        http.MethodGet,
			Pattern:       "/static/:w*",
			Handlers:      []http.HandlerFunc{StaticFilesHandler},
			TrailingSlash: true,
		},
		{
			Name:          "sse",
			Method:        http.MethodGet,
			Pattern:       "/sse/:clientID",
			Handlers:      []http.HandlerFunc{SSEHandler(sse)},
			TrailingSlash: true,
		},
	}
}

func setup() (*webgo.Router, *sse.SSE) {
	port := strings.TrimSpace(os.Getenv("HTTP_PORT"))
	if port == "" {
		port = "8080"
	}
	cfg := &webgo.Config{
		Host:         "",
		Port:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 1 * time.Hour,
	}

	webgo.GlobalLoggerConfig(
		nil, nil,
		webgo.LogCfgDisableDebug,
	)

	routeGroup := webgo.NewRouteGroup("/v6.2", false)
	routeGroup.Add(webgo.Route{
		Name:     "router-group-prefix-v6.2_api",
		Method:   http.MethodGet,
		Pattern:  "/api/:param",
		Handlers: []http.HandlerFunc{chain, ParamHandler},
	})
	routeGroup.Use(routegroupMiddleware)

	sseService := sse.New()
	routes := getRoutes(sseService)
	routes = append(routes, routeGroup.Routes()...)

	router := webgo.NewRouter(cfg, routes...)
	router.UseOnSpecialHandlers(accesslog.AccessLog)
	router.Use(errLogger, accesslog.AccessLog, cors.CORS(nil))
	return router, sseService
}

func main() {
	router, sseService := setup()
	// broadcast server time to all SSE listeners
	go func() {
		retry := time.Millisecond * 500
		for {
			now := time.Now().Format(time.RFC1123Z)
			sseService.Broadcast(sse.Message{
				Data:  now + fmt.Sprintf(" (%d)", sseService.ActiveClients()),
				Retry: retry,
			})
			time.Sleep(time.Second)
		}
	}()

	router.Start()
}
