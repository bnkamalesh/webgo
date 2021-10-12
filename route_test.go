package webgo

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"
)

func TestRoute_computePatternStr(t *testing.T) {
	t.Parallel()
	type fields struct {
		Name                    string
		Method                  string
		Pattern                 string
		TrailingSlash           bool
		FallThroughPostResponse bool
		Handlers                []http.HandlerFunc
		uriKeys                 []string
		uriPatternString        string
		uriPattern              *regexp.Regexp
		serve                   http.HandlerFunc
	}
	type args struct {
		patternString string
		hasWildcard   bool
		key           string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "duplicate URIs",
			fields: fields{
				Pattern: "/a/b/:c/:c",
				// uriKeys is initialized with a key, so as to detect duplicate key
				uriKeys: []string{"c"},
			},
			args: args{
				patternString: strings.Replace("/a/b/:c/:c", ":c", urlchars, 2),
				hasWildcard:   false,
				key:           "c",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Route{
				Name:                    tt.fields.Name,
				Method:                  tt.fields.Method,
				Pattern:                 tt.fields.Pattern,
				TrailingSlash:           tt.fields.TrailingSlash,
				FallThroughPostResponse: tt.fields.FallThroughPostResponse,
				Handlers:                tt.fields.Handlers,
				uriKeys:                 tt.fields.uriKeys,
				uriPatternString:        tt.fields.uriPatternString,
				uriPattern:              tt.fields.uriPattern,
				serve:                   tt.fields.serve,
			}
			got, err := r.computePatternStr(tt.args.patternString, tt.args.hasWildcard, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Route.computePatternStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Route.computePatternStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
