package webgo

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

// Start the server with the appropriate configurations
func Start(cfg *Config, router *httprouter.Router) {
	host := cfg.Host

	host += ":" + cfg.Port

	if cfg.Env == "production" {
		print("Starting server in production mode, listening on `" + host + "`\n")
		err := http.ListenAndServe(host, router)
		if err != nil {
			Err.Log(err)
		}
	} else {
		print("Starting server in development mode")
		n := negroni.Classic()
		n.UseHandler(router)
		n.Run(host)
	}

}

// ===
