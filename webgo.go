package webgo

import (
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Start the server with the appropriate configurations
func Start(cfg *Config, router *httprouter.Router) {
	if cfg.Env == "production" {
		print("Starting server, listening on `http://localhost:" + cfg.Port + "`\n")
		err := http.ListenAndServe(":"+cfg.Port, router)
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
		n := negroni.Classic()
		n.UseHandler(router)
		n.Run(":" + cfg.Port)
	}

}

// ===
