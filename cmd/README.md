# Webgo Sample

## How to run 

If you have Go installed on your computer, open the terminal and:

```bash
$ cd $GOPATH/src
$ mkdir -p github.com/bnkamalesh
$ cd github.com/bnkamalesh
$ git clone https://github.com/bnkamalesh/webgo.git
$ cd webgo
$ go run cmd/main.go

Info 2020/06/03 12:55:26 HTTP server, listening on :8080
```

Or if you have [Docker](https://www.docker.com/), open the terminal and:

```bash
$ git clone https://github.com/bnkamalesh/webgo.git
$ cd webgo
$ docker run \
-p 8080:8080 \
-v ${PWD}:/go/src/github.com/bnkamalesh/webgo/ \
-w /go/src/github.com/bnkamalesh/webgo/cmd \
--rm -ti golang:latest go run main.go

Info 2020/06/03 12:55:26 HTTP server, listening on :8080
```


You can try the following API calls with the sample app. It also uses all the features provided by webgo

1. `http://localhost:8080/`
	- Route with no named parameters configured
2. `http://localhost:8080/matchall/`
	- Route with wildcard parameter configured
	- All URIs which begin with `/matchall` will be matched because it has a wildcard variable
	- e.g. 
		- http://localhost:8080/matchall/hello
		- http://localhost:8080/matchall/hello/world
		- http://localhost:8080/matchall/hello/world/user
3. `http://localhost:8080/api/<param>`
	- Route with a named 'param' configured
	- It will match all requests which match `/api/<single parameter>`
	- e.g.
		- http://localhost:8080/api/hello
		- http://localhost:8080/api/world
4. `http://localhost:8080/error-setter`
	- Route which sets an error and sets response status 500
5. `http://localhost:8080/v5.4/api/<param>`
	- Route with a named 'param' configured
	- It will match all requests which match `/v5.4/api/<single parameter>`
	- e.g.
		- http://localhost:8080/v5.4/api/hello
		- http://localhost:8080/v5.4/api/world
