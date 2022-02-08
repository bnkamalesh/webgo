package webgo

import (
	"fmt"
	"net/http"
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
