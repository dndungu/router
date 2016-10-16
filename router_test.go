package router

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter(t *testing.T) {
	router := NewTestRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()
	tests := []struct {
		method   string
		path     string
		response string
	}{
		{"GET", "/abcd", "default:GET:"},
		{"POST", "/abcd", "default:POST:"},
		{"PUT", "/abcd", "default:PUT:"},
		{"DELETE", "/abcd", "default:DELETE:"},
		{"CONNECT", "/abcd", "default:CONNECT:"},
		{"PATCH", "/abcd", "default:PATCH:"},
		{"OPTIONS", "/abcd", "default:OPTIONS:"},
		{"TRACE", "/abcd", "default:TRACE:"},
		{"GET", "/abcd", "default:GET:"},

		{"GET", "abcd", "default:GET:"},
		{"POST", "abcd", "default:POST:"},
		{"PUT", "abcd", "default:PUT:"},
		{"DELETE", "abcd", "default:DELETE:"},
		{"CONNECT", "abcd", "default:CONNECT:"},
		{"PATCH", "abcd", "default:PATCH:"},
		{"OPTIONS", "abcd", "default:OPTIONS:"},
		{"TRACE", "abcd", "default:TRACE:"},
		{"GET", "abcd", "default:GET:"},

		{"GET", "/ab/cd", "default:GET:"},
		{"POST", "/ab/cd", "default:POST:"},
		{"PUT", "/ab/cd", "default:PUT:"},
		{"DELETE", "/ab/cd", "default:DELETE:"},
		{"CONNECT", "/ab/cd", "default:CONNECT:"},
		{"PATCH", "/ab/cd", "default:PATCH:"},
		{"OPTIONS", "/ab/cd", "default:OPTIONS:"},
		{"TRACE", "/ab/cd", "default:TRACE:"},
		{"GET", "/ab/cd/:id", "default:GET:"},

		{"GET", "/test", "root:GET:"},
		{"POST", "/test", "root:POST:"},
		{"PUT", "/test", "root:PUT:"},
		{"DELETE", "/test", "root:DELETE:"},
		{"CONNECT", "/test", "root:CONNECT:"},
		{"PATCH", "/test", "root:PATCH:"},
		{"OPTIONS", "/test", "root:OPTIONS:"},
		{"TRACE", "/test", "root:TRACE:"},
		{"GET", "/test/123", "root:GET:123"},
	}
	transport := &http.Transport{}
	client := &http.Client{Transport: transport}
	for _, test := range tests {
		url := NewTestURL(ts.URL, test.path)
		req, err := http.NewRequest(test.method, url, nil)
		if err != nil {
			t.Fatal(err)
		}
		w, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		actual, err := ioutil.ReadAll(w.Body)
		w.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if strings.Compare(string(actual), test.response) != 0 {
			t.Fatalf("Error, expected a body with '%s', got '%s'. path: %s, verb: %s", test.response, actual, test.path, test.method)
		}
	}
}

func NewTestRouter() *Router {
	router := New(NewTestHandler("default"))
	router.Get("/", NewTestHandler("root"))
	router.Get("/test", NewTestHandler("root"))
	router.Post("/test", NewTestHandler("root"))
	router.Put("/test", NewTestHandler("root"))
	router.Delete("/test", NewTestHandler("root"))
	router.Connect("/test", NewTestHandler("root"))
	router.Patch("/test", NewTestHandler("root"))
	router.Options("/test", NewTestHandler("root"))
	router.Trace("/test", NewTestHandler("root"))
	router.Get("/test/:id", NewTestHandler("root"))
	router.Add("GET", "/te/st", NewTestHandler("root"))
	return router
}

func NewTestHandler(body string) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		root := newNode()
		router := &Router{tree: root}
		id := router.Params(r).Get("id")
		response := r.Method + ":" + r.URL.Path + ":" + id + ":" + body
		response = body + ":" + r.Method + ":" + id
		w.Write([]byte(response))
	}
}

func NewTestURL(base, path string) string {
	var u bytes.Buffer
	u.WriteString(base)
	if path[0] != '/' {
		u.WriteString("/")
	}
	u.WriteString(path)
	return u.String()
}
