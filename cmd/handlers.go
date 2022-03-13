package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bnkamalesh/webgo/v6"
	"github.com/bnkamalesh/webgo/v6/extensions/sse"
)

// StaticFilesHandler is used to serve static files
func StaticFilesHandler(rw http.ResponseWriter, r *http.Request) {
	wctx := webgo.Context(r)
	// '..' is replaced to prevent directory traversal which could go out of static directory
	path := strings.ReplaceAll(wctx.Params()["w"], "..", "-")

	rw.Header().Set("Last-Modified", lastModified)
	http.ServeFile(rw, r, fmt.Sprintf("./cmd/static/%s", path))
}

func OriginalResponseWriterHandler(w http.ResponseWriter, r *http.Request) {
	rw := webgo.OriginalResponseWriter(w)
	if rw == nil {
		webgo.Send(w, "text/html", "got nil", http.StatusPreconditionFailed)
		return
	}

	webgo.Send(w, "text/html", "success", http.StatusOK)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fs, err := os.OpenFile("./cmd/static/index.html", os.O_RDONLY, 0600)
	if err != nil {
		webgo.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	info, err := fs.Stat()
	if err != nil {
		webgo.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out := make([]byte, info.Size())
	_, err = fs.Read(out)
	if err != nil {
		webgo.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(out)
}

func SSEHandler(sse *sse.SSE) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := webgo.Context(r).Params()
		r.Header.Set(sse.ClientIDHeader, params["clientID"])

		err := sse.Handler(w, r)
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Println("errorLogger:", err.Error())
			return
		}
	}
}

func ErrorSetterHandler(w http.ResponseWriter, r *http.Request) {
	err := errors.New("oh no, server error")
	webgo.SetError(r, err)

	webgo.R500(w, err.Error())
}

func ParamHandler(w http.ResponseWriter, r *http.Request) {
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

func InvalidJSONHandler(w http.ResponseWriter, r *http.Request) {
	webgo.R200(w, make(chan int))
}
