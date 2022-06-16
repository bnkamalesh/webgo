package webgo

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestRouteGroupsPathPrefix(t *testing.T) {
	t.Parallel()
	routes := []Route{
		{
			Name:     "r1",
			Pattern:  "/a",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{dummyHandler},
		},
		{
			Name:     "r2",
			Pattern:  "/b/:c",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{dummyHandler},
		},
		{
			Name:     "r3",
			Pattern:  "/:w*",
			Method:   http.MethodGet,
			Handlers: []http.HandlerFunc{dummyHandler},
		},
	}

	const prefix = "/v6.2"
	expectedSkipMiddleware := true
	rg := NewRouteGroup("/v6.2", expectedSkipMiddleware, routes...)

	list := rg.Routes()
	for idx := range list {
		route := list[idx]
		originalRoute := routes[idx]
		expectedPattern := fmt.Sprintf("%s%s", prefix, originalRoute.Pattern)
		if route.Pattern != expectedPattern {
			t.Errorf("Expected pattern %q, got %q", expectedPattern, route.Pattern)
		}
		if route.skipMiddleware != expectedSkipMiddleware {
			t.Errorf("Expected skip %v, got %v", expectedSkipMiddleware, route.skipMiddleware)
		}
	}
}

func dummyHandler(w http.ResponseWriter, r *http.Request) {}

func BenchmarkMatchWithWildcard(b *testing.B) {
	route := Route{
		Name:                    "widlcard",
		Method:                  http.MethodGet,
		TrailingSlash:           true,
		FallThroughPostResponse: true,
		Pattern:                 "/:w*/static1/:myvar/:w2*",
		Handlers:                []http.HandlerFunc{dummyHandler},
	}

	uri := "/hello/world/how/are/you/static1/hello2/world2/how2/are2/you2/static2"
	err := route.init()
	if err != nil {
		b.Error(err)
		return
	}

	for i := 0; i < b.N; i++ {
		ok, _ := route.matchPath(uri)
		if !ok {
			b.Errorf("Expected match, got no match")
			break
		}
	}
}

func TestMatchWithWildcard(t *testing.T) {
	route := Route{
		Name:                    "widlcard",
		Method:                  http.MethodGet,
		TrailingSlash:           true,
		FallThroughPostResponse: true,
		Pattern:                 "/:w*/static1/:myvar/:w2*/static2",
		Handlers:                []http.HandlerFunc{dummyHandler},
	}
	err := route.init()
	if err != nil {
		t.Error(err)
		return
	}

	uri := "/hello/world/how/are/you/static1/hello2/world2/how2/are2/you2/static2"
	wantParams := map[string]string{
		"w":     "hello/world/how/are/you",
		"myvar": "hello2",
		"w2":    "world2/how2/are2/you2",
	}
	matched, params := route.matchPath(uri)
	if !matched {
		t.Errorf("Expected match, got no match")
		return
	}
	if !reflect.DeepEqual(params, wantParams) {
		t.Errorf("Expected params %v, got %v", wantParams, params)
		return
	}

	t.Run("no match", func(t *testing.T) {
		route := Route{
			Name:                    "widlcard",
			Method:                  http.MethodGet,
			TrailingSlash:           true,
			FallThroughPostResponse: true,
			Pattern:                 "/:w*/static1/:myvar/:w2*/static2",
			Handlers:                []http.HandlerFunc{dummyHandler},
		}
		err := route.init()
		if err != nil {
			t.Error(err)
			return
		}

		uri := "/hello/world/how/are/you/static2/hello2/world2/how2/are2/you2/static2"
		matched, params := route.matchPath(uri)
		if matched {
			t.Errorf("Expected no match, got match")
			return
		}
		if params != nil {
			t.Errorf("Expected params %v, got %v", nil, params)
			return
		}
	})
	t.Run("match with more params", func(t *testing.T) {
		route := Route{
			Name:                    "widlcard",
			Method:                  http.MethodGet,
			TrailingSlash:           true,
			FallThroughPostResponse: true,
			Pattern:                 "/:w*/static1/:myvar/:w2*/static2/:myvar2/:w3*/static3",
			Handlers:                []http.HandlerFunc{dummyHandler},
		}
		err := route.init()
		if err != nil {
			t.Error(err)
			return
		}

		uri := "/hello/world/how/are/you/static1/hello2/world2/how2/are2/you2/static2/hello3/world3/how3/are3/you3/static3"
		wantParams := map[string]string{
			"w":      "hello/world/how/are/you",
			"myvar":  "hello2",
			"w2":     "world2/how2/are2/you2",
			"myvar2": "hello3",
			"w3":     "world3/how3/are3/you3",
		}
		matched, params := route.matchPath(uri)
		if !matched {
			t.Errorf("Expected match, got no match")
			return
		}
		if !reflect.DeepEqual(params, wantParams) {
			t.Errorf("Expected params %v, got %v", wantParams, params)
			return
		}
	})
	t.Run("match - end with wildcard", func(t *testing.T) {
		route := Route{
			Name:                    "widlcard",
			Method:                  http.MethodGet,
			TrailingSlash:           true,
			FallThroughPostResponse: true,
			Pattern:                 "/:w*/static1/:myvar/:w2*",
			Handlers:                []http.HandlerFunc{dummyHandler},
		}
		err := route.init()
		if err != nil {
			t.Error(err)
			return
		}

		uri := "/hello/world/how/are/you/static1/hello2/world2/how2/are2/you2/static2"
		wantParams := map[string]string{
			"w":     "hello/world/how/are/you",
			"myvar": "hello2",
			"w2":    "world2/how2/are2/you2/static2",
		}
		matched, params := route.matchPath(uri)
		if !matched {
			t.Errorf("Expected match, got no match")
			return
		}
		if !reflect.DeepEqual(params, wantParams) {
			t.Errorf("Expected params %v, got %v", wantParams, params)
			return
		}
	})
}
