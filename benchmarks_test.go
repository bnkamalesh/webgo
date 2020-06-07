package webgo

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func makeBenchReq(b *testing.B,
	router *Router,
	respRec *httptest.ResponseRecorder,
	url string,
) error {
	req, err := http.NewRequest(http.MethodGet, url, bytes.NewBuffer(nil))
	if err != nil {
		return fmt.Errorf("%s %s", url, err.Error())
	}
	router.ServeHTTP(respRec, req)
	if respRec.Result().StatusCode != http.StatusOK {
		return fmt.Errorf(
			"%s %s, expected %d, got %d",
			err.Error(),
			url,
			http.StatusOK,
			respRec.Result().StatusCode,
		)
	}
	return nil
}

func runbench(b *testing.B, url string) {
	router, respRec := setup()
	var err error
	// b.RunParallel(func(pb *testing.PB) {
	// 	for pb.Next() {
	// 		respRec = httptest.NewRecorder()
	// 		err = makeBenchReq(b, router, respRec, url)
	// 		if err != nil {
	// 			b.Fatal(err)
	// 		}
	// 	}
	// })

	for i := 0; i < b.N; i++ {
		err = makeBenchReq(b, router, respRec, url)
		if err != nil {
			b.Fatal(err)
		}
	}
}
func BenchmarkGetNoParams(b *testing.B) {
	url := strings.Join([]string{baseapi, "nparams"}, "/")
	runbench(b, url)
}

func BenchmarkGetWithParams(b *testing.B) {
	url := strings.Join([]string{baseapi, "wparams", p1, "goblin", p2}, "/")
	runbench(b, url)
}

func BenchmarkPostWithParams(b *testing.B) {
	url := strings.Join([]string{baseapi, "hello", p1, "goblin", p2}, "/")
	runbench(b, url)
}

func Benchmark_MatchAndGet(b *testing.B) {
	router, _ := setup()
	r := router.getHandlers[2]

	path := "/hello/world/goblin/spiderman"
	for i := 0; i < b.N; i++ {
		r.matchAndGet(path)
	}
}
