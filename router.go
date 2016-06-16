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

func NewRouter(handler Handler) *Router {
	root := NewNode()
	return &Router{tree: root, handler: handler}
}

func (r *Router) Route(method, path string, handlers ...Handler) {
	r.tree.Add(method, path, handlers...)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	params := url.Values{}
	keys := strings.Split(req.URL.Path, "/")[1:]
	node, _ := r.tree.Find(strings.Join(keys, "/"), params)
	context.Set(req, "path_params", params)
	handlers := node.handlers[req.Method]
	if handlers == nil {
		r.handler(w, req)
		return
	}
	for _, handler := range handlers {
		handler(w, req)
	}
	context.Clear(req)
}
