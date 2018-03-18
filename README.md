[![Build Status](https://travis-ci.org/bnkamalesh/webgo.svg?branch=master)](https://travis-ci.org/bnkamalesh/webgo)
[![](https://goreportcard.com/badge/github.com/bnkamalesh/webgo)](https://goreportcard.com/report/github.com/bnkamalesh/webgo)

# WebGo v2.0 (planned)

1. Current implementation of logging middleware will be deprecated
 - Logging middleware will be converted to wrapper
 e.g. middlewares.Log(Handler)
 This is a clean implementation, to avoid unnecessary complexity and 
 computation from ServeHTTP()
 
 2. No more capability of enabling access log on selected routes. Access log
 can either be turned on for all or not.

# WebGo v1.6.0

A lightweight & simple web framework for Go.
[GoDoc webgo](https://godoc.org/github.com/bnkamalesh/webgo)

### Update 18 March 2018 (v1.6.0)
 - Added Trailing slash feature. Added a new option `TrailingSlash` boolean
 to `Route` definition. If true, the provided URI pattern will be matched
 with or without the trailing slash. Default is false.

2. All middlewares will be moved to a nested folder `middlewares` rather than
being the methods of the type middlewares.

These changes would cause backward incompatibility

### Requirements

1. `Go 1.8` or higher


### Usage

Please refer to the Sample app built using webgo: [webgo-sample](https://github.com/bnkamalesh/webgo-sample) to see how webgo's capabilities can be used.

Supported HTTP methods are: `OPTIONS, HEAD, GET, POST, PUT, PATCH, DELETE`.

This framework does not force you to follow any architecture (e.g. MVC), instead is more of a configuration over convention based framework. While using any of the default HTTP response function available, for any status less than 400, the JSON response is wrapped as follows:

```
{
	data: <payload, any valid JSON data>,
	status: <status code, integer>
}
```

While using any of the default HTTP response function available, for any status greater than 400, the JSON response is wrapped as follows:

```
{
	errors: <payload, any valid JSON data>,
	status: <status code, integer>
}
```
### URI patterns
While defining the path of an HTTP route; you can choose from the following 
4 different options:

1. Static path e.g. `/v1/api/users`
2. Path with named parameters e.g. `/v1/api/users/:userID/photos/:photoID`
 - Here `userID` and `photoID` are the parameter names.
3. Path with wildcard e.g. `/v2/*`
4. Path with named wildcard parameter. e.g. `/v2/:param*`


### [Middlewares](https://github.com/bnkamalesh/webgo/blob/master/middlewares.go)

Middlewares and HTTP handlers have the same function signature (same as HTTP standard library's handler function). An execution chain (1 request passing through a set of middlewares and handler function) is stopped immediately after a response is sent. If you'd like the execution to continue even after a response is sent. Each `Route` specified has a property `FallThroughPostResponse` which if set to true will continue executing the chain, but no further responses will be written. You can see a sample [here](https://github.com/bnkamalesh/webgo-sample/blob/master/routes.go).
	
### [Available HTTP response functions (JSON)](https://github.com/bnkamalesh/webgo/blob/master/responses.go)

1. `R200(http.ResponseWriter, payload)` to send a JSON response with status 200

2. `R201(http.ResponseWriter, payload)` to send a JSON response with status 201

3. `R204(http.ResponseWriter)` to send a response header with status 204

4. `R302(http.ResponseWriter, payload)` to send a JSON response with status 302

5. `R400(http.ResponseWriter, payload)` to send a JSON response with status 400

6. `R403(http.ResponseWriter, payload)` to send a JSON response with status 403

7. `R404(http.ResponseWriter, payload)` to send a JSON response with status 404

7. `R406(http.ResponseWriter, payload)` to send a JSON response with status 406

8. `R451(http.ResponseWriter, payload)` to send a JSON response with status 451


### [Functions to send customized responses](https://github.com/bnkamalesh/webgo/blob/master/responses.go)

1. `SendResponse(http.ResponseWriter, payload, responseCode)` function in [responses.go](https://github.com/bnkamalesh/webgo/blob/master/responses.go) can be used to send a payload wrapped in the `data` struct, with any status code required. 

2. `SendError(http.ResponseWriter, payload, errorCode)` function in [responses.go](https://github.com/bnkamalesh/webgo/blob/master/responses.go) can be used to send a payload wrapped in the `errors` struct, with any status code required. 

3. `SendHeader(http.ResponseWriter, responseCode)` function in [responses.go](https://github.com/bnkamalesh/webgo/blob/master/responses.go) can be used to send a response header alone, with any HTTP status code required.

4. `Send(http.ResponseWriter, contentType, payload, responseCode)` function in [responses.go](https://github.com/bnkamalesh/webgo/blob/master/responses.go) can be used to send a completely custom response.

5. `Render(http.ResponseWriter, payload, responseCode, *template.Template)` can be used to render any template.

All HTTP responses are in [JSON](https://en.wikipedia.org/wiki/JSON) (if not rendering HTML templates and not using `Send`).


### Configuration

WebGo configuration can be loaded directly from a JSON file using the helper `Load("/path/to/json")` 
of the struct `Config`.
```
var cfg webgo.Config
cfg.Load("path/to/config.json")
```

Following options can be provided in the JSON file
```
{
	"environment":  "", // running mode, it can be "production" or "development"
	"host": "", // Host on which the app should run
	"port": "", // Port on which the app should listen to
	"httpsOnly": false, // If true, only HTTPS server is started
	"httpsPort":  "", // Port on which the HTTPS server should listen to
	"certFile": "", // Certificate file path for HTTPS
	"keyFile": "", // Private key file path of the certificate
	"templatePath": "" // Folder containing all the templates
}
```

### Bencmark

Simple hello world, JSON response `{data: "Hello world", status: 200}`.

Middlewares: [CORS](https://github.com/bnkamalesh/webgo/blob/master/middlewares.go)

Options    : Logging turned *off*.

Source     : [WebGo Sample](https://github.com/bnkamalesh/webgo-sample)

#### Specs: 

Machine   : MacBook Pro (Retina, 13-inch, Early 2015)

Processor : 2.7 GHz Intel Core i5

Memory    : 8 GB 1867 MHz DDR3

[Ulimit](http://wiki.linuxquestions.org/wiki/Ulimit)    : 50,000

```
$ ab -k -n 25000 -c 500 http://127.0.0.1:8000/
This is ApacheBench, Version 2.3 <$Revision: 1757674 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 2500 requests
Completed 5000 requests
Completed 7500 requests
Completed 10000 requests
Completed 12500 requests
Completed 15000 requests
Completed 17500 requests
Completed 20000 requests
Completed 22500 requests
Completed 25000 requests
Finished 25000 requests


Server Software:        
Server Hostname:        127.0.0.1
Server Port:            8000

Document Path:          /nparams
Document Length:        36 bytes

Concurrency Level:      500
Time taken for tests:   0.873 seconds
Complete requests:      25000
Failed requests:        0
Keep-Alive requests:    25000
Total transferred:      4200000 bytes
HTML transferred:       900000 bytes
Requests per second:    28621.57 [#/sec] (mean)
Time per request:       17.469 [ms] (mean)
Time per request:       0.035 [ms] (mean, across all concurrent requests)
Transfer rate:          4695.73 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    1   3.9      0      37
Processing:     4   17   5.8     15      43
Waiting:        4   16   5.8     15      43
Total:          4   17   7.8     15      66

Percentage of the requests served within a certain time (ms)
  50%     15
  66%     17
  75%     19
  80%     21
  90%     25
  95%     29
  98%     47
  99%     55
 100%     66 (longest request)
 ```

## How to run the tests?

You need to first start the HTTP server, you need to do the following:

```
$ cd /path/to/webgo/tests
$ go run main.go
```

After starting the server, you can run the tests by:

```
$ cd /path/to/webgo/tests
$ go test
```