package webgo

import (
	"net/http"

	"time"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

// Start the server with the appropriate configurations
func Start(cfg *Config, router *httprouter.Router, readTimeout, writeTimeout time.Duration) {
	host := cfg.Host

	if len(cfg.Port) > 0 {
		host += ":" + cfg.Port
	}

	if cfg.Env == "production" {
		print("Starting server in production mode, listening on `" + host + "`\n")
		httpServer := &http.Server{
			Addr:         host,
			Handler:      router,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		}
		// err := http.ListenAndServe(host, router)
		err := httpServer.ListenAndServe()
		if err != nil {
			Log.Println("Could not start http server -> ", err)
		}
	} else {
		// In development mode, it runs using Negroni, which provides basic logging like access logs,
		// panic handler etc.
		print("Starting server in development mode")
		n := negroni.Classic()
		n.UseHandler(router)
		n.Run(host)
	}
}
