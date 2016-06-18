package router

import (
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func Testinsert(t *testing.T) {
	tree := newNode()
	path := "/test-add"
	verb := "GET"
	tree.insert(verb, path, handlerA)
}

func Testsearch(t *testing.T) {
	tree := newNode()
	verb := "GET"
	params := url.Values{}
	path := "/test-find"
	tree.insert(verb, path, handlerB)
	leaf, _ := tree.search(path[1:], params)

	handlers := leaf.handlers[verb]
	if len(handlers) != 1 {
		t.Fatalf("Expected 1 handler, %d found", len(handlers))
		return
	}

	handlerName := FuncName(handlers[0])
	if strings.Compare("handlerB", handlerName) != 0 {
		t.Fatalf("Expected to find handlerB in the handlers list")
	}
}

func handlerA(w http.ResponseWriter, req *http.Request) {
}

func handlerB(w http.ResponseWriter, req *http.Request) {
}

func FuncName(h Handler) string {
	p := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	handlers := strings.Split(p, "/")
	longName := handlers[len(handlers)-1]
	return strings.Split(longName, ".")[1]
}
