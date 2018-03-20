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

// Start starts the server with the appropriate configurations
func Start(cfg *Config, router *Router, readTimeout, writeTimeout time.Duration) {
	host := cfg.Host
	httpshost := cfg.Host

	if len(cfg.Port) > 0 {
		host += ":" + cfg.Port
	}

	if len(cfg.HTTPSPort) > 0 {
		httpshost += ":" + cfg.HTTPSPort
	}

	if cfg.HTTPSOnly {
		if cfg.CertFile == "" {
			errLogger.Fatalln("No certificate provided for HTTPS")
		}

		if cfg.KeyFile == "" {
			errLogger.Fatalln("No key file provided for HTTPS")
		}

		httpsServer := &http.Server{
			Addr:         httpshost,
			Handler:      router,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: cfg.Env == "development",
			},
		}

		infoLogger.Println("\nStarting HTTPS server, listening on '" + httpshost + "'")
		err := httpsServer.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			infoLogger.Println("HTTPS Server exited with error:", err.Error())
		}
		return
	}

	if len(cfg.HTTPSPort) > 0 {
		if cfg.CertFile == "" {
			errLogger.Fatalln("No certificate provided for HTTPS")
		}

		if cfg.KeyFile == "" {
			errLogger.Fatalln("No key file provided for HTTPS")
		}

		if cfg.Port == cfg.HTTPSPort {
			errLogger.Fatalln("HTTP & HTTPS cannot listen on the same port. [" + cfg.Port + "]")
		}

		//Starting HTTPS server
		go func() {
			httpsServer := &http.Server{
				Addr:         httpshost,
				Handler:      router,
				ReadTimeout:  readTimeout,
				WriteTimeout: writeTimeout,
				TLSConfig: &tls.Config{
					InsecureSkipVerify: cfg.Env == "development",
				},
			}

			infoLogger.Println("Starting HTTPS server, listening on '" + httpshost + "'")
			err := httpsServer.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
			if err != nil {
				infoLogger.Println("HTTPS Server exited with error:", err.Error())
			}
			return
		}()
	}

	httpServer := &http.Server{
		Addr:         host,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	infoLogger.Println("Starting HTTP server, listening on '" + host + "'")
	err := httpServer.ListenAndServe()
	if err != nil {
		infoLogger.Println("HTTP Server exited with error:", err.Error())
	}

}
