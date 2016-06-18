package router

import (
	"bytes"
	"fmt"
	"github.com/gorilla/context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	request(router, "a/b")
}

func TestParams(t *testing.T) {
	methods := []string{"GET"}
	routes := []string{"/test/:name"}
	handlers := []Handler{paramHandler}
	router := initRouter(methods, routes, handlers)
	response := request(router, "/test/golang")
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

func request(router *Router, path string) string {
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
	u.WriteString("/test/golang")
	return u.String()
}

func paramHandler(w http.ResponseWriter, req *http.Request) {
	params := context.Get(req, "params")
	fmt.Fprint(w, params.(url.Values).Get("name"))
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
