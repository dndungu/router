// Package router - this is a golang HTTP router that implements a trie data structure for optimal performance
package router

import (
	"github.com/gorilla/context"
	"net/http"
	"net/url"
	"sync"
)

// Handler - this a http handler middleware function
type Handler func(http.ResponseWriter, *http.Request)

// Router - the router keeps a record of different paths and their handlers, it will direct incoming http calls to the right handler
type Router struct {
	tree    *node
	handler Handler
}

// New - this method creates a new instance of a router
func New(handler Handler) *Router {
	root := newNode()
	return &Router{tree: root, handler: handler}
}

// Connect - add a CONNECT handler for the specified path
func (r *Router) Connect(path string, handlers ...Handler) {
	r.tree.insert("CONNECT", path, handlers...)
}

// Delete - add a DELETE handler for the specified path
func (r *Router) Delete(path string, handlers ...Handler) {
	r.tree.insert("DELETE", path, handlers...)
}

// Get - add a GET handler for the specified path
func (r *Router) Get(path string, handlers ...Handler) {
	r.tree.insert("GET", path, handlers...)
}

// Post - add a POST handler for the specified path
func (r *Router) Post(path string, handlers ...Handler) {
	r.tree.insert("POST", path, handlers...)
}

// Put - add a PUT handler for the specified path
func (r *Router) Put(path string, handlers ...Handler) {
	r.tree.insert("PUT", path, handlers...)
}

// Patch - add a PATCH handler for the specified path
func (r *Router) Patch(path string, handlers ...Handler) {
	r.tree.insert("PATCH", path, handlers...)
}

// Trace - add a TRACE handler for the specified path
func (r *Router) Trace(path string, handlers ...Handler) {
	r.tree.insert("TRACE", path, handlers...)
}

// Options - add a OPTIONS handler for the specified path
func (r *Router) Options(path string, handlers ...Handler) {
	r.tree.insert("OPTIONS", path, handlers...)
}

// Add - this method adds a path and it's handlers to the router
func (r *Router) Add(method, path string, handlers ...Handler) {
	r.tree.insert(method, r.Path(path), handlers...)
}

// Params - this method returns any URL parameter if it exists
func (r *Router) Params(req *http.Request) url.Values {
	return Params(req)
}

// Params - this function returns any URL parameter if it exists
func Params(req *http.Request) url.Values {
	return context.Get(req, "params").(url.Values)
}

// ServeHTTP - this method is called every time a new request comes in
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var wg sync.WaitGroup
	params := url.Values{}
	context.Set(req, "params", params)
	path := r.Path(req.URL.Path)
	node, _ := r.tree.search(path, params)
	handlers := node.handlers[req.Method]
	if handlers == nil {
		handlers = []Handler{r.handler}
	}
	for _, h := range handlers {
		wg.Add(1)
		go func() {
			h(w, req)
			wg.Done()
		}()
	}
	wg.Wait()
	context.Clear(req)
}

// Path - append a front slash to paths without it
func (r *Router) Path(path string) string {
	if path[0] == '/' {
		path = path[1:]
	}
	return path
}
