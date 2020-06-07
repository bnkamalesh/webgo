package webgo

import (
	"net/http"
	"reflect"
	"regexp"
	"testing"
)

func TestRoute_matchAndGet(t *testing.T) {
	router, _ := setup()

	route := router.getHandlers[2]
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
		requestURI string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  map[string]string
	}{
		{
			name: "valid params",
			fields: fields{
				Method:                  route.Method,
				Pattern:                 route.Pattern,
				TrailingSlash:           route.TrailingSlash,
				FallThroughPostResponse: route.FallThroughPostResponse,
				Handlers:                route.Handlers,
				uriKeys:                 route.uriKeys,
				uriPatternString:        route.uriPatternString,
				uriPattern:              route.uriPattern,
				serve:                   route.serve,
			},
			args: args{
				requestURI: "/wparams/world/goblin/spiderman",
			},
			want: true,
			want1: map[string]string{
				"p1": "world",
				"p2": "spiderman",
			},
		},
		{
			name: "invalid URI",
			fields: fields{
				Method:                  route.Method,
				Pattern:                 route.Pattern,
				TrailingSlash:           route.TrailingSlash,
				FallThroughPostResponse: route.FallThroughPostResponse,
				Handlers:                route.Handlers,
				uriKeys:                 route.uriKeys,
				uriPatternString:        route.uriPatternString,
				uriPattern:              route.uriPattern,
				serve:                   route.serve,
			},
			args: args{
				requestURI: "/params/world/goblin/spiderman",
			},
			want:  false,
			want1: nil,
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
			got, got1 := r.matchAndGet(tt.args.requestURI)
			if got != tt.want {
				t.Errorf(
					"Route.matchAndGet() got = %v, want %v. route: %v",
					got,
					tt.want,
					r,
				)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Route.matchAndGet() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
