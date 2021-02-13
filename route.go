package webgo

import (
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
			return "", fmt.Errorf(
				"%s\nURI:%s\nKey:%s, Position: %d",
				errDuplicateKey,
				r.Pattern,
				k,
				idx+1,
			)
		}
	}

	r.uriKeys = append(r.uriKeys, key)
	return patternString, nil
}

func (r *Route) parseURIWithParams(patternString string) (string, error) {
	if !strings.Contains(r.Pattern, ":") {
		return "", nil
	}

	var err error
	// uriValues is a map of URI Key and its respective value,
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
				return "", err
			}
			hasWildcard, hasKey = false, false
			key = ""
		}
	}

	if hasKey && len(key) > 0 {
		patternString, err = r.computePatternStr(patternString, hasWildcard, key)
		if err != nil {
			return "", err
		}
	}
	return patternString, nil
}

// init prepares the URIKeys, compile regex for the provided pattern
func (r *Route) init() error {
	patternString := r.Pattern

	patternString, err := r.parseURIWithParams(patternString)
	if err != nil {
		return err
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

// matchPath matches the requestURI with the URI pattern of the route.
// If the path is an exact match (i.e. no URI parameters), then the second parameter ('isExactMatch') is true
func (r *Route) matchPath(requestURI string) (bool, isExactMatch bool) {
	if r.Pattern == requestURI {
		return true, true
	}

	return r.uriPattern.Match([]byte(requestURI)), false
}

func (r *Route) params(requestURI string) map[string]string {
	params := r.uriPattern.FindStringSubmatch(requestURI)[1:]
	uriValues := make(map[string]string, len(params))

	for i := 0; i < len(params); i++ {
		uriValues[r.uriKeys[i]] = params[i]
	}

	return uriValues
}

func routeServeChainedHandlers(r *Route) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		crw, ok := rw.(*customResponseWriter)
		if !ok {
			crw = newCRW(rw, http.StatusOK)
		}

		for _, handler := range r.Handlers {
			if crw.written && !r.FallThroughPostResponse {
				break
			}
			handler(crw, req)
		}
	}
}

func defaultRouteServe(r *Route) http.HandlerFunc {
	if len(r.Handlers) > 1 {
		return routeServeChainedHandlers(r)
	}

	// when there is only 1 handler, custom response writer is not required to check if response
	// is already written or fallthrough is enabled
	return r.Handlers[0]
}
