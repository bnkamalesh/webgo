#GoLang web framework.

A lightweight & simple web framework for GoLang.
[GoDoc webgo](https://godoc.org/github.com/bnkamalesh/webgo)

###Requirements

1. `GoLang 1.5` or higher

``` 
# Enable vendoring for Go1.5
$ export GO15VENDOREXPERIMENT=1
# ===
```

### Third party libraries used

1. [mgo/mango](http://gopkg.in/mgo.v2), MongoDB driver.
2. [HttpRouter](github.com/julienschmidt/httprouter), multiplexer.
3. [Stack](https://github.com/alexedwards/stack), for chaining request handlers.


### Usage
Most of the usage can be seen in this sample app.
Sample app built using webgo: [https://github.com/bnkamalesh/webgo-sample](https://github.com/bnkamalesh/webgo-sample)

This framework does not force you to follow any architecture (e.g. MVC), instead is more of a configuration over convention based framework. There are a very limited set of HHTP responses available by default in the framework. They are
`200, 201, 302, 400, 403, 404, 406, 451, 500`. If you need more, you can directly call the `SendResponse` function in [responses.go](https://github.com/bnkamalesh/webgo/blob/master/responses.go) with any status code you like.

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
	"templatePath": "", // Folder containing all the templates

	"dbConfig": { // This is by default meant for MongoDB. You can use it for any database
		"name":     "",
		"host":     "",
		"port":     "",
		"username": "",
		"password": "",
		"authSource": "",
		"mgoDialString": "" // The full dial string, can be provided instead of filling all the other fields. It uses `mgo` //driver.
	}
}
```
