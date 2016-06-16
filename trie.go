package router

import (
	"net/url"
	"strings"
)

type node struct {
	children []*node
	key      string
	param    bool
	handlers map[string][]Handler
}

func NewNode() *node {
	return &node{key: "", param: false, handlers: make(map[string][]Handler)}
}

func (n *node) Add(method, path string, handlers ...Handler) {
	i := 0
	if path[0] == '/' {
		i = 1
	}
	keys := strings.Split(path[i:], "/")
	count := len(keys)
	for {
		child, key := n.Find(path[i:], nil)
		if child.key == key && count == 1 {
			child.handlers[method] = handlers
			return
		}
		param := false
		if len(key) > 0 && key[0] == ':' {
			param = true
		}
		infant := node{key: key, param: param, handlers: make(map[string][]Handler)}
		if count == 1 {
			infant.handlers[method] = handlers
		}
		child.children = append(child.children, &infant)
		count--
		if count == 0 {
			break
		}
	}
}

func (n *node) Find(path string, params url.Values) (*node, string) {
	keys := strings.Split(path, "/")
	if len(n.children) == 0 {
		return n, keys[0]
	}
	return n.search(path, params)
}

func (n *node) search(path string, params url.Values) (*node, string) {
	keys := strings.Split(path, "/")
	for _, child := range n.children {
		if keys[0] != child.key && child.param == false {
			continue
		}
		if child.param && params != nil {
			params.Add(child.key[1:], keys[0])
		}
		next := strings.Join(keys[1:], "/")
		if len(next) > 0 {
			return child.Find(next, params)
		}
		return child, keys[0]
	}
	return n, keys[0]
}
