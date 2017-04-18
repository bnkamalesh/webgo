#Go web framework/boilerplate.

A lightweight & simple web framework for Go.
[GoDoc webgo](https://godoc.org/github.com/bnkamalesh/webgo)

###Requirements

1. `Go 1.7` or higher


### Usage
Most of the usage can be seen in this sample app.
Sample app built using webgo: [https://github.com/bnkamalesh/webgo-sample](https://github.com/bnkamalesh/webgo-sample)

This framework does not force you to follow any architecture (e.g. MVC), instead is more of a configuration over convention based framework. There are a very limited set of HHTP responses available by default in the framework. They are
`200, 201, 204, 302, 400, 403, 404, 406, 451, 500`. If you need more, you can directly call the `SendResponse` function in [responses.go](https://github.com/bnkamalesh/webgo/blob/master/responses.go) with any status code you like. `SendError` function in [responses.go](https://github.com/bnkamalesh/webgo/blob/master/responses.go) can be used to send a payload wrapped in error struct, with any error code you like. `SendHeader` function in [responses.go](https://github.com/bnkamalesh/webgo/blob/master/responses.go) can be used to send a response header alone, with any HTTP status code you require.

All HTTP responses are in [JSON](https://en.wikipedia.org/wiki/JSON) (if not rendering HTML templates). Any response with status code less than 400 will be wrapped in a JSON format `{data: "payload", status: 200/201}`. Every other response will be wrapped in `{errors: "payload", status: >= 400}`.

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

	"templatePath": "", // Folder containing all the templates
	"hideAccessLog": false // if true, access log will not be printed
}
```
