package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bnkamalesh/webgo/v4"
)

const baseapi = "http://127.0.0.1:9696/"

func TestMiddleware(t *testing.T) {
	router, respRec := setup(getRoutes())
	router.Use(CorsWrap("*"))
	router.Use(AccessLog)
	url := baseapi

	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err, url)
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != 200 {
		t.Fatal(err, respRec.Code, url)
	}

	h := respRec.Header().Get(headerAllowHeaders)
	if h != allowHeaders {
		t.Fatalf("\nExpected '%s'\ngot:'%s'", allowHeaders, h)
	}

	req, err = http.NewRequest(http.MethodOptions, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err, url)
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != 200 {
		t.Fatal(err, respRec.Code, url)
	}

	h = respRec.Header().Get(headerAllowHeaders)
	if h != allowHeaders {
		t.Fatalf("\nExpected '%s'\ngot:'%s'", allowHeaders, h)
	}
}

func TestCorsChainHandler(t *testing.T) {
	router, respRec := setup(getRoutesWithCorsChain())
	router.Use(AccessLog)
	url := baseapi
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err, url)
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != 200 {
		t.Fatal(err, respRec.Code, url)
	}

	h := respRec.Header().Get(headerAllowHeaders)
	if h != allowHeaders {
		t.Fatal("Expected ", allowHeaders, "\ngot", h)
	}
}

func TestCorsOptionsChain(t *testing.T) {
	router, respRec := setup(getRoutesWithCorsChain())
	router.Use(AccessLog)
	url := baseapi
	req, err := http.NewRequest(http.MethodOptions, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err, url)
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != 200 {
		t.Fatal(err, respRec.Code, url)
	}

	h := respRec.Header().Get(headerAllowHeaders)
	if h != allowHeaders {
		t.Fatal("Expected ", allowHeaders, "\ngot", h)
	}

	req, err = http.NewRequest(http.MethodOptions, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Error(err, url)
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != 200 {
		t.Error(err, respRec.Code, url)
	}

	h = respRec.Header().Get(headerAllowHeaders)
	if h != allowHeaders {
		t.Error("Expected ", allowHeaders, "\ngot", h)
	}
}

func TestChain(t *testing.T) {
	router, respRec := setup(getChainedRoutes())
	url := baseapi
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Error(err, url)
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != 200 {
		t.Error(err, respRec.Code, url)
	}

	h := respRec.Header().Get(headerAllowHeaders)
	if h != allowHeaders {
		t.Error("Expected ", allowHeaders, "\ngot", h)
	}

	req, err = http.NewRequest(http.MethodOptions, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Error(err, url)
	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != 200 {
		t.Error(err, respRec.Code, url)
	}

	h = respRec.Header().Get(headerAllowHeaders)
	if h != allowHeaders {
		t.Error("Expected ", allowHeaders, "\ngot", h)
	}
}

func setup(routes []*webgo.Route) (*webgo.Router, *httptest.ResponseRecorder) {
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
			// A label for the API/URI
			Name: "root",
			// request type
			Method:                  http.MethodGet,
			Pattern:                 "/",
			FallThroughPostResponse: true,
			TrailingSlash:           true,
			// route handler
			Handlers: []http.HandlerFunc{handler},
		},
	}
}
func getRoutesWithCorsChain() []*webgo.Route {
	return []*webgo.Route{
		{
			// A label for the API/URI
			Name: "OPTIONS",
			// request type
			Method:                  http.MethodOptions,
			Pattern:                 "/:w*",
			FallThroughPostResponse: true,
			TrailingSlash:           true,
			// route handler
			Handlers: []http.HandlerFunc{CorsOptions("*"), handler},
		},
		{
			// A label for the API/URI
			Name: "root",
			// request type
			Method:                  http.MethodGet,
			Pattern:                 "/",
			FallThroughPostResponse: true,
			TrailingSlash:           true,
			// route handler
			Handlers: []http.HandlerFunc{Cors("*"), handler},
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
			Handlers:                []http.HandlerFunc{CorsOptions("*")}, // route handler
		},
		{
			Name:                    "root",         // A label for the API/URI
			Method:                  http.MethodGet, // request type
			Pattern:                 "/",
			FallThroughPostResponse: true, // Pattern for the route
			TrailingSlash:           true,
			Handlers:                []http.HandlerFunc{Cors("*"), handler}, // route handler
		},
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	webgo.R200(w, "hello world")
}
