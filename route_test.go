package webgo

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRoute_init(t *testing.T) {
	type fields struct {
		Name          string
		Method        string
		Pattern       string
		TrailingSlash bool
		Handlers      []http.HandlerFunc
	}
	tests := []struct {
		name           string
		fields         fields
		wantPatternStr string
		wantErr        bool
	}{
		{
			name: "without trailing slash, without params",
			fields: fields{
				Name:          "",
				Method:        http.MethodGet,
				Pattern:       "/user",
				TrailingSlash: false,
				Handlers:      []http.HandlerFunc{http.NotFound},
			},
			wantErr:        false,
			wantPatternStr: "^/user$",
		},
		{
			name: "with trailing slash, without params",
			fields: fields{
				Name:          "",
				Method:        http.MethodGet,
				Pattern:       "/user",
				TrailingSlash: true,
				Handlers:      []http.HandlerFunc{http.NotFound},
			},
			wantErr:        false,
			wantPatternStr: fmt.Sprintf("^/user%s$", trailingSlash),
		},
		{
			name: "without trailing slash, with params",
			fields: fields{
				Name:          "",
				Method:        http.MethodGet,
				Pattern:       "/user/:id",
				TrailingSlash: false,
				Handlers:      []http.HandlerFunc{http.NotFound},
			},
			wantErr:        false,
			wantPatternStr: fmt.Sprintf("^/user/%s$", "([^/]+)"),
		},
		{
			name: "with trailing slash, with params",
			fields: fields{
				Name:          "",
				Method:        http.MethodGet,
				Pattern:       "/user/:id",
				TrailingSlash: true,
				Handlers:      []http.HandlerFunc{http.NotFound},
			},
			wantErr:        false,
			wantPatternStr: fmt.Sprintf("^/user/%s%s$", "([^/]+)", trailingSlash),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Route{
				Name:          tt.fields.Name,
				Method:        tt.fields.Method,
				Pattern:       tt.fields.Pattern,
				TrailingSlash: tt.fields.TrailingSlash,
				Handlers:      tt.fields.Handlers,
			}
			err := r.init()
			if err != nil {
				t.Error(err.Error())
			}
			if r.uriPatternString != tt.wantPatternStr {
				t.Errorf("Expected pattern '%s', got '%s'", tt.wantPatternStr, r.uriPatternString)
			}
		})
	}
}
