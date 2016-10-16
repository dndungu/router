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

func newNode() *node {
	return &node{key: "", param: false, handlers: make(map[string][]Handler)}
}

func (n *node) insert(method, path string, handlers ...Handler) {
	keys := strings.Split(path, "/")
	count := len(keys)
	for {
		child, key := n.search(path[1:], nil)
		if strings.Compare(child.key, key) == 0 && count == 1 {
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

func (n *node) search(path string, params url.Values) (*node, string) {
	keys := strings.Split(path, "/")
	if len(n.children) == 0 {
		return n, keys[0]
	}
	for _, child := range n.children {
		if strings.Compare(keys[0], child.key) != 0 && child.param == false {
			continue
		}
		if child.param && params != nil {
			params.Add(child.key[1:], keys[0])
		}
		next := strings.Join(keys[1:], "/")
		if len(next) > 0 {
			return child.search(next, params)
		}
		return child, keys[0]
	}
	return n, keys[0]
}
