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

// StartHTTPS starts the server with HTTPS enabled
func StartHTTPS(cfg *Config, router *Router) {
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

// Start starts the server with the appropriate configurations
func Start(cfg *Config, router *Router) {
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
