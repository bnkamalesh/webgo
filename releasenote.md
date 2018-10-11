# Release notes

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