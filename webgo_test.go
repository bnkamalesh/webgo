package webgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type response struct {
	Data   map[string]string `json:"data"`
	Status int               `json:"status"`
}

// appConfig is a sample struct to hold app configurations
type appConfig struct {
	Name string
}

const p1 = "world"
const p2 = "spiderman"

const baseapi = "http://127.0.0.1:9696"
const baseapiHTTPS = "http://127.0.0.1:9696"

var GETAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

var POSTAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

var PUTAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

var DELETEAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

var PATCHAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

var OPTIONSAPI = []string{
	strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/"),
	strings.Join([]string{baseapiHTTPS, "hello", p1, "goblin", p2}, "/"),
}

func withrequestbody(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		R400(w, err)
		return
	}
	r.Body.Close()

	output := string(b)
	wctx := Context(r)
	params := wctx.Params()
	R200(
		w,
		map[string]string{
			"p1":      params["p1"],
			"p2":      params["p2"],
			"payload": output,
			"pattern": r.URL.Path,
			"method":  r.Method,
		},
	)
}

func withuriparams(w http.ResponseWriter, r *http.Request) {
	wctx := Context(r)
	params := wctx.Params()
	R200(
		w,
		map[string]string{
			"p1":      params["p1"],
			"p2":      params["p2"],
			"pattern": r.URL.Path,
			"method":  r.Method,
		},
	)
}

func httpresponsewriter(w http.ResponseWriter, r *http.Request) {
	payload, _ := json.Marshal(
		map[string]string{
			"pattern": r.URL.Path,
			"method":  r.Method,
		},
	)
	w.Write(payload)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	R200(w, "Hello world")
}

func getRoutes() []*Route {
	return []*Route{
		{
			Name:                    "root",
			Method:                  http.MethodGet,
			Pattern:                 "/",
			FallThroughPostResponse: true,
			TrailingSlash:           true,
			Handlers:                []http.HandlerFunc{helloWorld},
		},
		{
			Name:     "hw-noparams",
			Method:   http.MethodGet,
			Pattern:  "/nparams",
			Handlers: []http.HandlerFunc{helloWorld},
		},
		{
			Name:          "hw-withparams",
			Method:        http.MethodGet,
			TrailingSlash: true,
			Pattern:       "/wparams/:p1/goblin/:p2",
			Handlers:      []http.HandlerFunc{withuriparams},
		},
		{
			Name:     "params-get",
			Method:   http.MethodGet,
			Pattern:  "/hello/:p1/goblin/:p2",
			Handlers: []http.HandlerFunc{withuriparams},
		},
		{
			Name:     "params-head",
			Method:   http.MethodHead,
			Pattern:  "/hello/:p1/goblin/:p2",
			Handlers: []http.HandlerFunc{withuriparams},
		},

		{
			Name:     "params-post-sameuri",
			Method:   http.MethodPost,
			Pattern:  "/hello/:p1/goblin/:p2",
			Handlers: []http.HandlerFunc{withrequestbody},
		},
		{
			Name:     "params-put-sameuri",
			Method:   http.MethodPut,
			Pattern:  "/hello/:p1/goblin/:p2",
			Handlers: []http.HandlerFunc{withrequestbody},
		},
		{
			Name:     "params-patch-sameuri",
			Method:   http.MethodPatch,
			Pattern:  "/hello/:p1/goblin/:p2",
			Handlers: []http.HandlerFunc{withrequestbody},
		},
		{
			Name:     "params-delete-sameuri",
			Method:   http.MethodDelete,
			Pattern:  "/hello/:p1/goblin/:p2",
			Handlers: []http.HandlerFunc{withrequestbody},
		},
		{
			Name:     "params-options-sameuri",
			Method:   http.MethodOptions,
			Pattern:  "/hello/:p1/goblin/:p2",
			Handlers: []http.HandlerFunc{withuriparams},
		},
		{
			Name:     "httpresponsewriter",
			Method:   http.MethodGet,
			Pattern:  "/httpresponsewriter",
			Handlers: []http.HandlerFunc{httpresponsewriter},
		},
	}
}

func mware(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	rw.Header().Add("k1", "v1")
	next(rw, req)
}

func setup() (*Router, *httptest.ResponseRecorder) {
	// Initializing router with all the required routes
	router := NewRouter(&Config{
		Host:               "127.0.0.1",
		Port:               "9696",
		HTTPSPort:          "8443",
		CertFile:           "tests/ssl/server.crt",
		KeyFile:            "tests/ssl/server.key",
		ReadTimeout:        15,
		WriteTimeout:       60,
		ShutdownTimeout:    4 * time.Second,
		InsecureSkipVerify: true,
	}, getRoutes())

	return router, httptest.NewRecorder()
}

func TestInvalidHTTPMethod(t *testing.T) {
	router, respRec := setup()

	for _, url := range GETAPI {
		req, err := http.NewRequest("ABC", url, bytes.NewBuffer(nil))
		if err != nil {
			t.Fatal(err, url)
			continue
		}

		router.ServeHTTP(respRec, req)
		if respRec.Code != http.StatusNotImplemented {
			t.Fatalf(`Expected response HTTP status code %d, received %d`, http.StatusNotImplemented, respRec.Code)
		}
	}
}

func TestGet(t *testing.T) {
	router, respRec := setup()

	for _, url := range GETAPI {
		req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
		if err != nil {
			t.Fatal(err, url)
			continue
		}

		router.ServeHTTP(respRec, req)

		resp := response{}

		err = json.NewDecoder(respRec.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)
			continue
		}

		if resp.Data["method"] != http.MethodGet {
			t.Fatal("URL:", url, "response method:", resp.Data["method"], " required method:", http.MethodGet)
		}

		if resp.Data["p1"] != p1 {
			t.Fatal("p1:", resp.Data["p1"])
		}

		if resp.Data["p2"] != p2 {
			t.Fatal("p2:", resp.Data["p2"])
		}
	}
}

func TestMiddleware(t *testing.T) {
	router, respRec := setup()
	router.Use(mware)
	url := baseapi + "/"
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err, url)
	}
	router.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusOK {
		t.Fatalf("Expected status '200', got '%d'", respRec.Code)
	}

	wctx := Context(req)
	if wctx == nil {
		t.Fatalf("Expected webgo context, got nil")
	}

	v := respRec.Header().Get("k1")
	if respRec.Header().Get("k1") != "v1" {
		t.Fatal("Expected response header value `v1` for key `k1`, received", v)
	}

	// test middleware with 404 request
	router, respRec = setup()
	router.UseOnSpecialHandlers(mware)
	url = fmt.Sprintf("%s/random/unimplemented/path", baseapi)
	req, err = http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err, url)
	}

	router.ServeHTTP(respRec, req)
	if respRec.Code != http.StatusNotFound {
		t.Fatalf(
			"Expected status '404', got '%d'",
			respRec.Code,
		)
	}
	v = respRec.Header().Get("k1")
	if respRec.Header().Get("k1") != "v1" {
		t.Fatalf("Expected response header value `v1` for key `k1`, received '%s'", v)
	}

	// test middleware with 501 request
	router, respRec = setup()
	router.UseOnSpecialHandlers(mware)
	req, err = http.NewRequest("UNIMPLEMENTED", baseapi, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err, url)
	}
	router.ServeHTTP(respRec, req)
	if respRec.Code != http.StatusNotImplemented {
		t.Fatalf(
			"Expected status '%d', got '%d'",
			http.StatusNotImplemented,
			respRec.Code,
		)
	}
	v = respRec.Header().Get("k1")
	if respRec.Header().Get("k1") != "v1" {
		t.Fatal("Expected response header value `v1` for key `k1`, received", v)
	}
}
func TestGetPostResponse(t *testing.T) {
	router, respRec := setup()
	url := baseapi + "/"
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err, url)

	}
	router.ServeHTTP(respRec, req)

	if respRec.Code != http.StatusOK {
		t.Fatal(err, respRec.Code, url)
	}
}
func TestGet404(t *testing.T) {
	router, respRec := setup()
	url := baseapi + "/random"
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatal(err, url)

	}

	router.ServeHTTP(respRec, req)

	if respRec.Code != 404 {
		t.Fatal(err)

	}
}
func TestHead(t *testing.T) {
	router, respRec := setup()

	for _, url := range GETAPI {
		req, err := http.NewRequest(http.MethodHead, url, bytes.NewBuffer(nil))
		if err != nil {
			t.Fatal(err, url)

			continue
		}

		router.ServeHTTP(respRec, req)

		resp := response{}

		err = json.NewDecoder(respRec.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)

			continue
		}

		if resp.Data["method"] != http.MethodHead {
			t.Fatal("URL:", url, "response method:", resp.Data["method"], " required method:", http.MethodGet)

		}

		if resp.Data["p1"] != p1 {
			t.Fatal("p1:", resp.Data["p1"])

		}

		if resp.Data["p2"] != p2 {
			t.Fatal("p2:", resp.Data["p2"])

		}
	}
}

func TestPost(t *testing.T) {
	router, respRec := setup()
	var payload = []byte(`{"payload": "nothing"}`)

	for _, url := range POSTAPI {
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err, url)

			continue
		}
		router.ServeHTTP(respRec, req)
		resp := response{}
		err = json.NewDecoder(respRec.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)

			continue
		}

		if resp.Data["method"] != http.MethodPost {
			t.Fatal("response method:", resp.Data["method"], " required method:", http.MethodPost)

		}

		if resp.Data["p1"] != p1 {
			t.Fatal("p1:", resp.Data["p1"])

		}

		if resp.Data["p2"] != p2 {
			t.Fatal("p2:", resp.Data["p2"])

		}

		if resp.Data["payload"] != string(payload) {
			t.Fatal("payload:", resp.Data["payload"])

		}
	}
}

func TestPut(t *testing.T) {
	router, respRec := setup()
	var payload = []byte(`{"payload": "nothing"}`)
	for _, url := range PUTAPI {

		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err, url)

			continue
		}
		router.ServeHTTP(respRec, req)
		resp := response{}
		err = json.NewDecoder(respRec.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)

			continue
		}

		if resp.Data["method"] != http.MethodPut {
			t.Fatal("response method:", resp.Data["method"], " required method:", http.MethodPut)

		}

		if resp.Data["p1"] != p1 {
			t.Fatal("p1:", resp.Data["p1"])

		}

		if resp.Data["p2"] != p2 {
			t.Fatal("p2:", resp.Data["p2"])

		}

		if resp.Data["payload"] != string(payload) {
			t.Fatal("payload:", resp.Data["payload"])

		}
	}
}

func TestPatch(t *testing.T) {
	router, respRec := setup()
	var payload = []byte(`{"payload": "nothing"}`)
	for _, url := range PATCHAPI {
		req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err, url)

			continue
		}

		resp := response{}
		router.ServeHTTP(respRec, req)
		err = json.NewDecoder(respRec.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)

			continue
		}

		if resp.Data["method"] != http.MethodPatch {
			t.Fatal("response method:", resp.Data["method"], " required method:", http.MethodPatch)

		}

		if resp.Data["p1"] != p1 {
			t.Fatal("p1:", resp.Data["p1"])

		}

		if resp.Data["p2"] != p2 {
			t.Fatal("p2:", resp.Data["p2"])

		}

		if resp.Data["payload"] != string(payload) {
			t.Fatal("payload:", resp.Data["payload"])

		}
	}
}

func TestDelete(t *testing.T) {
	router, respRec := setup()
	var payload = []byte(`{"payload": "nothing"}`)
	for _, url := range DELETEAPI {
		req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err, url)

			continue
		}

		resp := response{}
		router.ServeHTTP(respRec, req)
		err = json.NewDecoder(respRec.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)

			continue
		}

		if resp.Data["method"] != http.MethodDelete {
			t.Fatal("response method:", resp.Data["method"], " required method:", http.MethodDelete)

		}

		if resp.Data["p1"] != p1 {
			t.Fatal("p1:", resp.Data["p1"])

		}

		if resp.Data["p2"] != p2 {
			t.Fatal("p2:", resp.Data["p2"])

		}

		if resp.Data["payload"] != string(payload) {
			t.Fatal("payload:", resp.Data["payload"])

		}
	}
}

func TestOptions(t *testing.T) {
	router, respRec := setup()
	var payload = []byte(`{"payload": "nothing"}`)

	for _, url := range OPTIONSAPI {
		req, err := http.NewRequest(http.MethodOptions, url, bytes.NewBuffer(payload))
		if err != nil {
			t.Fatal(err, url)

			continue
		}

		resp := response{}
		router.ServeHTTP(respRec, req)
		err = json.NewDecoder(respRec.Body).Decode(&resp)
		if err != nil {
			t.Fatal(err)
			continue
		}

		if resp.Data["method"] != http.MethodOptions {
			t.Fatal("response method:", resp.Data["method"], " required method:", http.MethodOptions)
		}

		if resp.Data["p1"] != p1 {
			t.Fatal("p1:", resp.Data["p1"])
		}

		if resp.Data["p2"] != p2 {
			t.Fatal("p2:", resp.Data["p2"])
		}
	}
}

func TestHTTPResponseWriter(t *testing.T) {
	router, respRec := setup()
	path := "httpresponsewriter"
	url := fmt.Sprintf("%s/%s", baseapi, path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err, url)
	}

	resp := response{}
	router.ServeHTTP(respRec, req)
	if respRec.Result().StatusCode != http.StatusOK {
		t.Fatalf(
			"expected status 200, got '%d', url '%s'",
			respRec.Result().StatusCode,
			url,
		)
	}

	err = json.NewDecoder(respRec.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err, url)
	}

	if resp.Data["pattern"] != path {
		fmt.Sprintf(
			"expected pattern '%s', got '%s'",
			path,
			resp.Data["pattern"],
		)
	}
}
func TestStart(t *testing.T) {
	router, _ := setup()
	go router.Start()
	time.Sleep(time.Second * 3)
	err := router.Shutdown()
	if err != nil {
		t.Fatal(err)
	}
}
func TestStartHTTPS(t *testing.T) {
	router, _ := setup()
	go router.StartHTTPS()
	time.Sleep(time.Second * 3)
	err := router.ShutdownHTTPS()
	if err != nil {
		t.Fatal(err)
	}
}
