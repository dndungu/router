package router

import (
	"github.com/gorilla/context"
	"net/http"
	"net/url"
	"strings"
)

type Handler func(http.ResponseWriter, *http.Request)

type Router struct {
	tree    *node
	handler Handler
}

func New(handler Handler) *Router {
	root := newNode()
	return &Router{tree: root, handler: handler}
}

func (r *Router) Add(method, path string, handlers ...Handler) {
	r.tree.insert(method, path, handlers...)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	params := url.Values{}
	context.Set(req, "params", params)
	keys := strings.Split(req.URL.Path, "/")[1:]
	node, _ := r.tree.search(strings.Join(keys, "/"), params)
	handlers := node.handlers[req.Method]
	if handlers == nil {
		r.handler(w, req)
		return
	}
	for _, h := range handlers {
		h(w, req)
	}
	//	context.Clear(req)
}
