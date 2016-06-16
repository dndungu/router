package router

import (
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	tree := NewNode()
	path := "/test-add"
	verb := "GET"
	tree.Add(verb, path, handlerA)
}

func TestFind(t *testing.T) {
	tree := NewNode()
	verb := "GET"
	params := url.Values{}
	path := "/test-find"
	tree.Add(verb, path, handlerB)
	leaf, _ := tree.Find(path[1:], params)

	handlers := leaf.handlers[verb]
	if len(handlers) != 1 {
		t.Fatalf("Expected 1 handler, %d found", len(handlers))
		return
	}

	handler_name := FuncName(handlers[0])
	if strings.Compare("handlerB", handler_name) != 0 {
		t.Error("Test")
		t.Error(handler_name)
		t.Errorf("Expected to find handlerB in the handlers list")
	}
}

func handlerA(w http.ResponseWriter, req *http.Request) {
}

func handlerB(w http.ResponseWriter, req *http.Request) {
}

func FuncName(h Handler) string {
	p := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	handlers := strings.Split(p, "/")
	long_name := handlers[len(handlers)-1]
	return strings.Split(long_name, ".")[1]
}
