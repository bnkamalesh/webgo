// This is part of the library that is being used directly inside this package.
// https://github.com/alexedwards/stack
package webgo

import (
	"sync"
)

type Context struct {
	mu sync.RWMutex
	m  map[string]interface{}
}

// Create a new context
func NewContext() *Context {
	m := make(map[string]interface{})
	return &Context{m: m}
}

// Get context, and lock it while it's being accessed
func (c *Context) Get(key string) interface{} {
	if !c.Exists(key) {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.m[key]
}

// Add/Update a key,val to the context
func (c *Context) Put(key string, val interface{}) *Context {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = val
	return c
}

// Delete a key from the context
func (c *Context) Delete(key string) *Context {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.m, key)
	return c
}

// Check if a key exists in the context
func (c *Context) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.m[key]
	return ok
}

// Copy a context
func (c *Context) copy() *Context {
	nc := NewContext()
	c.mu.RLock()
	c.mu.RUnlock()
	for k, v := range c.m {
		nc.m[k] = v
	}
	return nc
}
