package webgo

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

// Start the server with the appropriate configurations
func Start(cfg *Config, router *httprouter.Router) {
	host := cfg.Host
	if len(host) <= 0 {
		host = "127.0.0.1"
	}

	host += ":" + cfg.Port

	if cfg.Env == "production" {
		print("Starting server in production mode, listening on `http://" + host + "`\n")
		err := http.ListenAndServe(host, router)
		if err != nil {
			// Log error
			Err.Log("webgo.go", "Start()", err)

			// Stay idle for 2 seconds
			// time.Sleep(1000 * 1000 * 2)

			// Start server again
			// This is not a graceful restart, since all the clients connected to the server will be disconnected
			// Graceful restart of the server should be implemented.
			print("Server exited, no restart implemented")
		}
	} else {
		print("Starting server in development mode")
		n := negroni.Classic()
		n.UseHandler(router)
		n.Run(host)
	}

}

// ===
