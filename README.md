#Go web framework/boilerplate.

A lightweight & simple web framework for Go.
[GoDoc webgo](https://godoc.org/github.com/bnkamalesh/webgo)


###Requirements

1. `Go 1.8` or higher


### Usage
Please refer to the Sample app built using webgo: [https://github.com/bnkamalesh/webgo-sample](https://github.com/bnkamalesh/webgo-sample) to see how webgo's capabilities can be used.

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
	
### Available HTTP response functions (JSON)
1. `R200(http.ResponseWriter, payload)` to send a JSON response with status 200

2. `R201(http.ResponseWriter, payload)` to send a JSON response with status 201

3. `R204(http.ResponseWriter)` to send a response header with status 204

4. `R302(http.ResponseWriter, payload)` to send a JSON response with status 302

5. `R400(http.ResponseWriter, payload)` to send a JSON response with status 400

6. `R403(http.ResponseWriter, payload)` to send a JSON response with status 403

7. `R404(http.ResponseWriter, payload)` to send a JSON response with status 404

7. `R406(http.ResponseWriter, payload)` to send a JSON response with status 406

8. `R451(http.ResponseWriter, payload)` to send a JSON response with status 451


### Functions to send customized responses

1. `SendResponse(http.ResponseWriter, payload, responseCode)` function in [responses.go](https://github.com/bnkamalesh/webgo/blob/master/responses.go) can be used to send a payload wrapped in the `data` struct, with any status code required. 

2. `SendError(http.ResponseWriter, payload, errorCode)` function in [responses.go](https://github.com/bnkamalesh/webgo/blob/master/responses.go) can be used to send a payload wrapped in the `errors` struct, with any status code required. 

3. `SendHeader(http.ResponseWriter, responseCode)` function in [responses.go](https://github.com/bnkamalesh/webgo/blob/master/responses.go) can be used to send a response header alone, with any HTTP status code required.

4. `Send(http.ResponseWriter, contentType, payload, responseCode)` function in [responses.go](https://github.com/bnkamalesh/webgo/blob/master/responses.go) can be used to send a completely custom response.

5. `Render(http.ResponseWriter, payload, responseCode, *template.Template)` can be used to render any template.

All HTTP responses are in [JSON](https://en.wikipedia.org/wiki/JSON) (if not rendering HTML templates and not using `Send`).

### Configuration

The app starts with configuration set in `config.json`. Configuration path(relative or absolute) can be provided as follows:

```
var cfg webgo.Config
cfg.Load("path/to/config.json")
```

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
