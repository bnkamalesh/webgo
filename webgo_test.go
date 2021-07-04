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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestResponseStatus(t *testing.T) {
	w := newCRW(httptest.NewRecorder(), http.StatusOK)
	SendError(w, nil, http.StatusNotFound)
	if http.StatusNotFound != ResponseStatus(w) {
		t.Errorf(
			"Expected status '%d', got '%d'",
			http.StatusNotFound,
			ResponseStatus(w),
		)
	}

	// ideally we should get 200 from ResponseStatus; but it can get accurate status code only
	// when `customresponsewriter` is used
	rw := httptest.NewRecorder()
	SendError(rw, nil, http.StatusNotFound)
	if http.StatusOK != ResponseStatus(rw) {
		t.Errorf(
			"Expected status '%d', got '%d'",
			http.StatusOK,
			ResponseStatus(rw),
		)
	}
}

func TestStart(t *testing.T) {
	router, _ := setup("9696")
	go router.Start()
	time.Sleep(time.Second * 2)
	err := router.Shutdown()
	if err != nil {
		t.Fatal(err)
	}
}
func TestStartHTTPS(t *testing.T) {
	router, _ := setup("8443")
	go router.StartHTTPS()
	time.Sleep(time.Second * 2)
	err := router.ShutdownHTTPS()
	if err != nil {
		t.Fatal(err)
	}
}

func TestErrorHandling(t *testing.T) {
	err := errors.New("hello world, failed")
	router, _ := setup("7878")
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	router.ServeHTTP(w, r)

	SetError(r, err)
	gotErr := GetError(r)

	if !errors.Is(err, gotErr) {
		t.Fatalf("expected err %v, got %v", err, gotErr)
	}
}
