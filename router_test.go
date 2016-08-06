package router

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var defaultHandler = NewHandler("defaultHandler")
var handlerZero = NewHandler("handlerZero")
var handlerOne = NewHandler("handlerOne")
var handlerTwo = NewHandler("handlerTwo")
var handlerThree = NewHandler("handlerThree")
var handlerFour = NewHandler("handlerFour")
var handlerFive = NewHandler("handlerFive")
var handlerSix = NewHandler("handlerSix")

func NewHandler(body string) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	}
}

func TestdefaultHandler(t *testing.T) {
	expected := "default"
	r := New(NewHandler(expected))
	actual := get(r, "/")
	if strings.Compare(actual, expected) != 0 {
		t.Fatalf("Expected a body with '%s' got '%s'", expected, actual)
	}
}

func TestPaths(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	routes := []string{"/a", "/a/b", "/a/b/c", "/a/b/c/d"}
	handlers := []Handler{handlerZero, handlerOne, handlerTwo, handlerThree}
	router := initRouter(methods, routes, handlers)
	get(router, "a/b")
}

func TestVerbs(t *testing.T) {
	router := New(defaultHandler)
	router.Get("/t/b", handlerZero)
	router.Post("t/b", handlerOne)
	router.Put("t/b", handlerTwo)
	router.Delete("t/b", handlerThree)
	router.Connect("t/b", handlerFour)
	router.Patch("t/b", handlerFive)
	router.Options("/t/b", handlerFive)
	router.Trace("t/b", handlerSix)
	response := get(router, "/t/b")
	if strings.Compare(response, "handlerZero") != 0 {
		t.Fatalf("Expected a body with 'handlerZero', got %s", response)
	}
}

func TestParams(t *testing.T) {
	methods := []string{"GET"}
	routes := []string{"/test/:name"}
	handlers := []Handler{paramHandler}
	router := initRouter(methods, routes, handlers)
	response := get(router, "/test/golang")
	if strings.Compare(response, "golang") != 0 {
		t.Fatalf("Expected a body with 'golang' got '%s'", response)
	}
}

func initRouter(methods []string, routes []string, handlers []Handler) *Router {
	router := New(defaultHandler)
	for key, path := range routes {
		router.Add(methods[key], path, handlers[key])
	}
	return router
}

func get(router *Router, path string) string {
	ts := httptest.NewServer(router)
	defer ts.Close()
	res, _ := http.Get(testURL(ts.URL, path))
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return string(body)
}

func testURL(baseURL, path string) string {
	var u bytes.Buffer
	u.WriteString(baseURL)
	u.WriteString(path)
	return u.String()
}

func paramHandler(w http.ResponseWriter, req *http.Request) {
	router := New(defaultHandler)
	body1 := router.Params(req).Get("name")
	body2 := Params(req).Get("name")
	if strings.Compare(body1, body2) != 0 {
		panic(fmt.Sprintf("Expected %s to equal %s", body1, body2))
	}
	fmt.Fprint(w, body2)
}
