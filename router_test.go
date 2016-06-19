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

func TestDefaultHandler(t *testing.T) {
	r := New(DefaultHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()
	res, _ := http.Get(ts.URL)
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if strings.Compare(string(body), "DefaultHandler") != 0 {
		t.Fatalf("Expected a body with 'DefaultHandler' got '%s'", body)
	}
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "DefaultHandler")
}

func TestPaths(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	routes := []string{"/a", "/a/b", "/a/b/c", "/a/b/c/d"}
	handlers := []Handler{handlerZero, handlerOne, handlerTwo, handlerThree}
	router := initRouter(methods, routes, handlers)
	get(router, "a/b")
}

func TestVerbs(t *testing.T) {
	router := New(DefaultHandler)
  router.Get("/t/b", handlerZero)
  router.Post("t/b", handlerOne)
  router.Put("t/b", handlerTwo)
  router.Delete("t/b", handlerThree)
  router.Connect("t/b", handlerFour)
  router.Patch("t/b", handlerFive)
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
	router := New(DefaultHandler)
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
	params := Params(req)
	fmt.Fprint(w, params.Get("name"))
}

func handlerZero(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "handlerZero")
}

func handlerOne(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "handlerOne")
}

func handlerTwo(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "handlerTwo")
}

func handlerThree(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "handlerThree")
}

func handlerFour(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "handlerFour")
}

func handlerFive(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "handlerFive")
}

func handlerSix(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "handlerSix")
}
