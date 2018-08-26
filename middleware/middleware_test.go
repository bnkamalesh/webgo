package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bnkamalesh/webgo"
)

const baseapi = "http://127.0.0.1:9696/"

func TestMiddleware(t *testing.T) {
	router, respRec := setup(false)
	router.Use(CorsWrap)
	router.Use(AccessLog)
	url := baseapi
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Log(err, url)
		t.Fail()
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != 200 {
		t.Log(err, respRec.Code, url)
		t.Fail()
	}

	h := respRec.Header().Get(headerAllowHeaders)
	if h != allowHeaders {
		t.Log("Expected ", allowHeaders, "\ngot", h)
		t.Fail()
	}

	req, err = http.NewRequest(http.MethodOptions, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Log(err, url)
		t.Fail()
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != 200 {
		t.Log(err, respRec.Code, url)
		t.Fail()
	}

	h = respRec.Header().Get(headerAllowHeaders)
	if h != allowHeaders {
		t.Log("Expected ", allowHeaders, "\ngot", h)
		t.Fail()
	}
}

func TestChain(t *testing.T) {
	router, respRec := setup(true)
	url := baseapi
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Log(err, url)
		t.Fail()
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != 200 {
		t.Log(err, respRec.Code, url)
		t.Fail()
	}

	h := respRec.Header().Get(headerAllowHeaders)
	if h != allowHeaders {
		t.Log("Expected ", allowHeaders, "\ngot", h)
		t.Fail()
	}

	req, err = http.NewRequest(http.MethodOptions, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Log(err, url)
		t.Fail()
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != 200 {
		t.Log(err, respRec.Code, url)
		t.Fail()
	}

	h = respRec.Header().Get(headerAllowHeaders)
	if h != allowHeaders {
		t.Log("Expected ", allowHeaders, "\ngot", h)
		t.Fail()
	}
}

func setup(chain bool) (*webgo.Router, *httptest.ResponseRecorder) {
	routes := getRoutes()
	if chain {
		routes = getChainedRoutes()
	}
	// Initializing router with all the required routes
	router := webgo.NewRouter(&webgo.Config{
		Host:               "127.0.0.1",
		Port:               "9696",
		HTTPSPort:          "8443",
		CertFile:           "tests/ssl/server.crt",
		KeyFile:            "tests/ssl/server.key",
		ReadTimeout:        15,
		WriteTimeout:       60,
		InsecureSkipVerify: true,
	}, routes)

	return router, httptest.NewRecorder()
}

func getRoutes() []*webgo.Route {
	return []*webgo.Route{
		{
			Name:                    "root",         // A label for the API/URI
			Method:                  http.MethodGet, // request type
			Pattern:                 "/",
			FallThroughPostResponse: true, // Pattern for the route
			TrailingSlash:           true,
			Handlers:                []http.HandlerFunc{handler}, // route handler
		},
	}
}
func getChainedRoutes() []*webgo.Route {
	return []*webgo.Route{
		{
			Name:                    "root",             // A label for the API/URI
			Method:                  http.MethodOptions, // request type
			Pattern:                 "/:w*",
			FallThroughPostResponse: true, // Pattern for the route
			TrailingSlash:           true,
			Handlers:                []http.HandlerFunc{CorsOptions}, // route handler
		},
		{
			Name:                    "root",         // A label for the API/URI
			Method:                  http.MethodGet, // request type
			Pattern:                 "/",
			FallThroughPostResponse: true, // Pattern for the route
			TrailingSlash:           true,
			Handlers:                []http.HandlerFunc{Cors, handler}, // route handler
		},
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	webgo.R200(w, "hello world")
}
