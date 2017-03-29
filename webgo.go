package webgo

import (
	"net/http"
	"time"
)

// Start the server with the appropriate configurations
func Start(cfg *Config, router *Router, readTimeout, writeTimeout time.Duration) {
	host := cfg.Host

	if len(cfg.Port) > 0 {
		host += ":" + cfg.Port
	}

	println("\nStarting server, listening on '" + host + "'")

	httpServer := &http.Server{
		Addr:         host,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		println("Server exited with error:", err.Error())
	}
}
