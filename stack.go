// This is part of the library that is being used directly inside this package.
// https://github.com/alexedwards/stack
package webgo

import "net/http"

type chainHandler func(*Context) http.Handler
type ChainMiddleware func(*Context, http.Handler) http.Handler

type Chain struct {
	mws []ChainMiddleware
	h   chainHandler
}

func NewStack(mws ...ChainMiddleware) Chain {
	return Chain{mws: mws}
}

func (c Chain) Append(mws ...ChainMiddleware) Chain {
	newMws := make([]ChainMiddleware, len(c.mws)+len(mws))
	copy(newMws[:len(c.mws)], c.mws)
	copy(newMws[len(c.mws):], mws)
	c.mws = newMws
	return c
}

func (c Chain) Then(chf func(ctx *Context, w http.ResponseWriter, r *http.Request)) HandlerChain {
	c.h = adaptContextHandlerFunc(chf)
	return newHandlerChain(c)
}

func (c Chain) ThenHandler(h http.Handler) HandlerChain {
	c.h = adaptHandler(h)
	return newHandlerChain(c)
}

func (c Chain) ThenHandlerFunc(fn func(http.ResponseWriter, *http.Request)) HandlerChain {
	c.h = adaptHandlerFunc(fn)
	return newHandlerChain(c)
}

type HandlerChain struct {
	context *Context
	Chain
}

func newHandlerChain(c Chain) HandlerChain {
	return HandlerChain{context: NewContext(), Chain: c}
}

func (hc HandlerChain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Always take a copy of context (i.e. pointing to a brand new memory location)
	ctx := hc.context.copy()

	final := hc.h(ctx)
	for i := len(hc.mws) - 1; i >= 0; i-- {
		final = hc.mws[i](ctx, final)
	}
	final.ServeHTTP(w, r)
}

func StackInject(hc HandlerChain, key string, val interface{}) HandlerChain {
	hc.context = hc.context.copy().Put(key, val)
	return hc
}

// Adapt third party middleware with the signature
// func(http.Handler) http.Handler into ChainMiddleware
func Adapt(fn func(http.Handler) http.Handler) ChainMiddleware {
	return func(ctx *Context, h http.Handler) http.Handler {
		return fn(h)
	}
}

// Adapt http.Handler into a chainHandler
func adaptHandler(h http.Handler) chainHandler {
	return func(ctx *Context) http.Handler {
		return h
	}
}

// Adapt a function with the signature
// func(http.ResponseWriter, *http.Request) into a chainHandler
func adaptHandlerFunc(fn func(w http.ResponseWriter, r *http.Request)) chainHandler {
	return adaptHandler(http.HandlerFunc(fn))
}

// Adapt a function with the signature
// func(Context, http.ResponseWriter, *http.Request) into a chainHandler
func adaptContextHandlerFunc(fn func(ctx *Context, w http.ResponseWriter, r *http.Request)) chainHandler {
	return func(ctx *Context) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fn(ctx, w, r)
		})
	}
}
