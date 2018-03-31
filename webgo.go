/*
Package webgo is a lightweight framework for building web apps. It has a multiplexer,
middleware plugging mechanism & context management of its own. The primary goal
of webgo is to get out of the developer's way as much as possible. i.e. it does
not enforce you to build your app in any particular pattern, instead just helps you
get all the trivial things done faster and easier.

e.g.
1. Getting named URI parameters.
2. Multiplexer for regex matching of URI and such.
3. Inject special app level configurations or any such objects to the request context as required.
*/
package webgo

import (
	"crypto/tls"
	"net/http"
	"time"
)

// WC is the webgocontext
type WC struct {
	Params     map[string]string
	Route      *Route
	AppContext map[string]interface{}
}

// Context returns the WebgoContext injected inside the HTTP request context
func Context(r *http.Request) *WC {
	wc, _ := r.Context().Value(wgoCtxKey).(*WC)
	return wc
}

// StartHTTPS starts the server with HTTPS enabled
func (router *Router) StartHTTPS() {
	cfg := router.config
	if cfg.CertFile == "" {
		errLogger.Fatalln("No certificate provided for HTTPS")
	}

	if cfg.KeyFile == "" {
		errLogger.Fatalln("No key file provided for HTTPS")
	}

	host := cfg.Host
	if len(cfg.HTTPSPort) > 0 {
		host += ":" + cfg.HTTPSPort
	}

	httpsServer := &http.Server{
		Addr:         host,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout * time.Second,
		WriteTimeout: cfg.WriteTimeout * time.Second,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: cfg.InsecureSkipVerify,
		},
	}

	infoLogger.Println("HTTPS server, listening on", host)
	err := httpsServer.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		errLogger.Fatalln("HTTPS server exited with error:", err.Error())
	}
}

// Start starts the HTTP server with the appropriate configurations
func (router *Router) Start() {
	cfg := router.config
	host := cfg.Host

	if len(cfg.Port) > 0 {
		host += ":" + cfg.Port
	}

	httpServer := &http.Server{
		Addr:         host,
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout * time.Second,
		WriteTimeout: cfg.WriteTimeout * time.Second,
	}

	infoLogger.Println("HTTP server, listening on '" + host + "'")
	err := httpServer.ListenAndServe()
	if err != nil {
		errLogger.Fatalln("HTTP server exited with error:", err.Error())
	}
}
