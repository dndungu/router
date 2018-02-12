package pkg

import (
	"fmt"
	"net/http"
)

type setFunc func(key, value string)

type node struct {
	nodes    map[string]*node
	handlers map[string][]http.HandlerFunc
}

func newNode() *node {
	return &node{
		nodes:    make(map[string]*node),
		handlers: make(map[string][]http.HandlerFunc),
	}
}

func (n *node) insert(keys []string, verb string, handlers ...http.HandlerFunc) error {
	if len(keys) == 0 {
		if h, ok := n.handlers[verb]; ok {
			n.handlers[verb] = append(h, handlers...)
		} else {
			n.handlers[verb] = handlers
		}
		return nil
	}

	k := keys[0]

	if _, ok := n.nodes[k]; ok {
		return n.nodes[k].insert(keys[1:], verb, handlers...)
	}

	if k[0] == ':' {
		for key := range n.nodes {
			if key[0] == ':' {
				if k != key {
					return fmt.Errorf("Error, cannot add %s, another parameter %s already exists.", k, key)
				}
			}
		}
	}

	n.nodes[k] = newNode()
	return n.nodes[k].insert(keys[1:], verb, handlers...)
}

func (n *node) search(keys []string, verb string, set setFunc) ([]http.HandlerFunc, bool) {
	if len(keys) == 0 {
		return n.handlers[verb], true
	}
	k := keys[0]
	if _, ok := n.nodes[k]; ok {
		return n.nodes[k].search(keys[1:], verb, set)
	}
	for key := range n.nodes {
		if key[0] == ':' {
			set(key, k)
			return n.nodes[key].search(keys[1:], verb, set)
		}
	}
	return nil, false
}
