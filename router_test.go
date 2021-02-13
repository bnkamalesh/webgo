package webgo

import (
	"net/http"
	"regexp"
	"testing"
)

func TestRoute_computePatternStr(t *testing.T) {
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
			name: "",
			fields: fields{
				Pattern: "/hello/:key1/:key2",
				uriKeys: []string{},
			},
			args: args{
				patternString: "/hello/:key1/:key2",
				hasWildcard:   false,
				key:           "key1",
			},
			want:    "/hello/([^/]+)/:key2",
			wantErr: false,
		},
		{
			name: "",
			fields: fields{
				Pattern: "/hello/:key1/:key2",
			},
			args: args{
				patternString: "/hello/:key1/:key2",
				hasWildcard:   false,
				key:           "key2",
			},
			want:    "/hello/:key1/([^/]+)",
			wantErr: false,
		},
		{
			name: "",
			fields: fields{
				Pattern: "/hello/:key1/:key2",
				uriKeys: []string{"key1"},
			},
			args: args{
				patternString: "/hello/:key1/:key2",
				hasWildcard:   false,
				key:           "key1",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "",
			fields: fields{
				Pattern: "/hello/:key1/:key2*",
			},
			args: args{
				patternString: "/hello/:key1/:key2*",
				hasWildcard:   true,
				key:           "key2",
			},
			want:    "/hello/:key1/(.*)",
			wantErr: false,
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
