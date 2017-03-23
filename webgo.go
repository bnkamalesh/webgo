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

	println("Starting server in production mode, listening on `" + host + "`\n")

	httpServer := &http.Server{
		Addr:         host,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		Log.Println("Could not start http server -> ", err)
	}
}
