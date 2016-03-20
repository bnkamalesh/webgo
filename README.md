#GoLang web framework.

A lightweight & simple web framework for GoLang.

###Requirements

1. `GoLang 1.6`

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

The default database driver available is `MongoDB`, and its handler can be accessed from
the configuration.

Any data retrieved using this handler will be in [`bson.M`](https://godoc.org/labix.org/v2/mgo/bson#M) format.
You can also pass a struct pointer get data from the handler instead of `bson.M`.

P.S: API documentation is not avaialble yet, will update as soon as possible.

`Sample: https://github.com/bnkamalesh/webgo-sample`