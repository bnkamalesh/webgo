package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bnkamalesh/webgo/v5"
	"github.com/bnkamalesh/webgo/v5/middleware/accesslog"
	"github.com/bnkamalesh/webgo/v5/middleware/cors"
)

var (
	lastModified = time.Now().Format(http.TimeFormat)
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

func errorSetter(w http.ResponseWriter, r *http.Request) {
	err := errors.New("oh no, server error")
	webgo.SetError(r, err)

	webgo.R500(w, err.Error())
}

func originalResponseWriter(w http.ResponseWriter, r *http.Request) {
	rw := webgo.OriginalResponseWriter(w)
	if rw == nil {
		webgo.Send(w, "text/html", "got nil", http.StatusPreconditionFailed)
		return
	}

	webgo.Send(w, "text/html", "success", http.StatusOK)
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

// StaticFiles is used to serve static files
func StaticFiles(rw http.ResponseWriter, r *http.Request) {
	wctx := webgo.Context(r)
	// '..' is replaced to prevent directory traversal which could go out of static directory
	path := strings.ReplaceAll(wctx.Params()["w"], "..", "-")

	rw.Header().Set("Last-Modified", lastModified)
	http.ServeFile(rw, r, fmt.Sprintf("./cmd/static/%s", path))
}

func getRoutes() []*webgo.Route {
	return []*webgo.Route{
		{
			Name:          "root",
			Method:        http.MethodGet,
			Pattern:       "/",
			Handlers:      []http.HandlerFunc{helloHandler},
			TrailingSlash: true,
		},
		{
			Name:          "matchall",
			Method:        http.MethodGet,
			Pattern:       "/matchall/:wildcard*",
			Handlers:      []http.HandlerFunc{helloWorld},
			TrailingSlash: true,
		},
		{
			Name:                    "api",
			Method:                  http.MethodGet,
			Pattern:                 "/api/:param",
			Handlers:                []http.HandlerFunc{chain, helloWorld},
			TrailingSlash:           true,
			FallThroughPostResponse: true,
		},
		{
			Name:          "invalidjson",
			Method:        http.MethodGet,
			Pattern:       "/invalidjson",
			Handlers:      []http.HandlerFunc{invalidJSON},
			TrailingSlash: true,
		},
		{
			Name:          "error-setter",
			Method:        http.MethodGet,
			Pattern:       "/error-setter",
			Handlers:      []http.HandlerFunc{errorSetter},
			TrailingSlash: true,
		},
		{
			Name:          "original-responsewriter",
			Method:        http.MethodGet,
			Pattern:       "/original-responsewriter",
			Handlers:      []http.HandlerFunc{originalResponseWriter},
			TrailingSlash: true,
		},
		{
			Name:          "static",
			Method:        http.MethodGet,
			Pattern:       "/static/:w*",
			Handlers:      []http.HandlerFunc{StaticFiles},
			TrailingSlash: true,
		},
	}
}

func main() {
	port := strings.TrimSpace(os.Getenv("HTTP_PORT"))
	if port == "" {
		port = "8080"
	}
	cfg := &webgo.Config{
		Host:         "",
		Port:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	router := webgo.NewRouter(cfg, getRoutes())
	router.UseOnSpecialHandlers(accesslog.AccessLog)
	router.Use(errLogger, accesslog.AccessLog, cors.CORS(nil))

	webgo.GlobalLoggerConfig(
		nil, nil,
		webgo.LogCfgDisableDebug,
	)

	router.Start()
}
