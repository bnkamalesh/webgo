package webgo

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Route defines a route for each API
type Route struct {
	// Name is unique identifier for the route
	Name string
	// Method is the HTTP request method/type
	Method string
	// Pattern is the URI pattern to match
	Pattern string
	// TrailingSlash if set to true, the URI will be matched with or without
	// a trailing slash. Note: It does not *do* a redirect.
	TrailingSlash bool

	// FallThroughPostResponse if enabled will execute all the handlers even if a response was already sent to the client
	FallThroughPostResponse bool

	// Handlers is a slice of http.HandlerFunc which can be middlewares or anything else. Though only 1 of them will be allowed to respond to client.
	// subsequent writes from the following handlers will be ignored
	Handlers []http.HandlerFunc

	// uriKeys is the list of URI parameter variables available for this route
	uriKeys []string
	// uriPatternString is the pattern string which is compiled to regex object
	uriPatternString string
	// uriPattern is the compiled regex to match the URI pattern
	uriPattern *regexp.Regexp

	serve http.HandlerFunc
}

// computePatternStr computes the pattern string required for generating the route's regex.
// It also adds the URI parameter key to the route's `keys` field
func (r *Route) computePatternStr(patternString string, hasWildcard bool, key string) (string, error) {
	regexPattern := ""
	patternKey := ""
	if hasWildcard {
		patternKey = fmt.Sprintf(":%s*", key)
		regexPattern = urlwildcard
	} else {
		patternKey = fmt.Sprintf(":%s", key)
		regexPattern = urlchars
	}

	patternString = strings.Replace(patternString, patternKey, regexPattern, 1)

	for idx, k := range r.uriKeys {
		if key == k {
			return "", errors.New(
				fmt.Sprintf(
					"%s\nURI:%s\nKey:%s, Position: %d",
					errDuplicateKey,
					r.Pattern,
					k,
					idx+1,
				),
			)
		}
	}

	r.uriKeys = append(r.uriKeys, key)
	return patternString, nil
}

// init prepares the URIKeys, compile regex for the provided pattern
func (r *Route) init() error {
	patternString := r.Pattern
	err := error(nil)

	if strings.Contains(r.Pattern, ":") {
		// uriValues is a map of URI Key and it's respective value,
		// this is calculated per request
		key := ""
		hasKey := false
		hasWildcard := false

		for i := 0; i < len(r.Pattern); i++ {
			char := string(r.Pattern[i])

			if char == ":" {
				hasKey = true
			} else if char == "*" {
				hasWildcard = true
			} else if hasKey && char != "/" {
				key += char
			} else if hasKey && len(key) > 0 {
				patternString, err = r.computePatternStr(patternString, hasWildcard, key)
				if err != nil {
					return err
				}
				hasWildcard, hasKey = false, false
				key = ""
			}
		}

		if hasKey && len(key) > 0 {
			patternString, err = r.computePatternStr(patternString, hasWildcard, key)
			if err != nil {
				return err
			}
		}
	}

	if r.TrailingSlash {
		patternString = fmt.Sprintf("^%s%s$", patternString, trailingSlash)
	} else {
		patternString = fmt.Sprintf("^%s$", patternString)
	}

	// compile the regex for the pattern string calculated
	reg, err := regexp.Compile(patternString)
	if err != nil {
		return err
	}

	r.uriPattern = reg
	r.uriPatternString = patternString
	r.serve = defaultRouteServe(r)

	return nil
}

// matchAndGet returns if the request URI matches the pattern defined in a Route as well as
// all the URI parameters configured for the route.
func (r *Route) matchAndGet(requestURI string) (bool, map[string]string) {
	if r.Pattern == requestURI {
		return true, nil
	}

	if !r.uriPattern.Match([]byte(requestURI)) {
		return false, nil
	}

	// Getting URI parameters
	values := r.uriPattern.FindStringSubmatch(requestURI)
	if len(values) == 0 {
		return true, nil
	}

	uriValues := make(map[string]string, len(values)-1)
	for i := 1; i < len(values); i++ {
		uriValues[r.uriKeys[i-1]] = values[i]
	}
	return true, uriValues

}

func defaultRouteServe(r *Route) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		crw, _ := rw.(*customResponseWriter)
		if crw == nil {
			crw = &customResponseWriter{
				ResponseWriter: rw,
			}
		}
		for _, handler := range r.Handlers {
			if !crw.written {
				// If there has been no write to response writer yet
				handler(crw, req)
			} else if r.FallThroughPostResponse {
				// run a handler post response write, only if fall through is enabled
				handler(crw, req)
			} else {
				// Do not run any more handlers if already responded and no fall through enabled
				break
			}
		}
	}
}
