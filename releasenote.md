# Release notes

## Release - 28/01/2019

### v3.3.0
#### Note: all changes from v3.0.0 to now are consolidated and listed

1. Updated tests to use .Erorr instead of Log & Fail
2. Updated the sample app to use CORS middleware
3. 'Method not implemented' status is now executed using a handler similar to 'NotFound'
4. overwriting the actual request instance
5. refactored for better readability
6. removed some variable declarations to avoid memory allocation
7. updated how middleware is used on NotFound & NotImplemented handlers, new method `UseOnSpecialHandlers` should be used to apply middleware on 'NotFound' & 'NotImplemented' handlers
	- added test to check if webgo context is available in middleware
8. Updated middleware package to use webgo v3
9. updated sample app to show chaining
10. updated ServeHTTP to inject webgo context into the request earlier, to make it available to middleware as well
11. Updated travis settings
12. Using codecov.io for coverage
13. patched test command to ignore `cmd` directory
14. updated test to use v3 instead of latest
15. Implemented `http.Hijacker` for the custom response writer used in the router

## Release - 06/10/2019

### v3.0.0

1. Using fmt.Errorf instead of using fmt.Sprintf to form string and creating error
2. Renamed webgo context from "WC" to "ContextPayload" for readability

## Release - 11/10/2018

### v2.4.0

1. Updated Readme
2. Updated Logger
	- Added a new logging interface
3. Updated tests to cover more code
4. Updated responses to use constants defined in Go's http standard library instead of using
integers for respective HTTP response codes

## Release - 09/10/2018

### v2.3.2

1. Updated Readme
2. Updated Middleware
	- Backward incompatible updates
	- CORS middleware functions are updated and now accepts list of supported/allowed domains
	- The middleware functions will default to "*" if no domains are passed to the functions

## Release - 03/05/2018

### v2.2.4

1. 501 response for unsupported HTTP methods

## Release - 12/04/2018

### v2.2.3

1. Fixed the long standing invalid http status code [bug](https://github.com/bnkamalesh/webgo/issues/7)
2. Fixed bug in access-log middleware which caused invalid HTTP status code
3. Updated docs with all latest updates

## Release - 08/04/2018

### v2.2.0

1. Graceful shutdown added
2. Updated readme with details of how to use the shutdown

## Release - 02/04/2018

### v2.1.0

1. Updated Readme to include godoc badge
2. Renamed `middlewares` to `middleware`

### v2.1.1

1. Initializing `AppContext` in NewRouter to avoid nil pointer assignment

## Release v2.0.0 - 01/04/2018

1. Log levels
	1. Error logs are now printed to `os.Stderr` with a prefix `Error`
	2. Info logs are now printed to `os.Stdout` with a prefix `Info`
	3. Warning logs are now printed to `os.Stdout` with a prefix `Warning`

2. Removed per route access log control
3. Removed AccessLog option from router
4. `Globals` is now removed, and router holds the app context now
5. Removed templates, as it can just be added to Globals' app context 
(`App` which is a map[string]interface{})
6. Removed configuration `HTTPSOnly` as it can be started by calling `StartHTTPS`
7. Read and write timeout are now added in configuration instead of passing to Start
8. `Start` & `StartHTTPS` are now methods of router