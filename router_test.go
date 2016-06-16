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
	r := NewRouter(DefaultHandler)
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

func TestPath(t *testing.T) {
	r := NewRouter(DefaultHandler)
	r.Route("GET", "/test/:name", SomeHandler)
	r.Route("GET", "/test/dosimilar", SomeHandler)
	r.Route("GET", "/test/dosimilar/still", SomeHandler)
	ts := httptest.NewServer(r)
	defer ts.Close()
	var u bytes.Buffer
	u.WriteString(ts.URL)
	u.WriteString("/test/golang")
	res, _ := http.Get(u.String())
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if strings.Compare(string(body), "golang") != 0 {
		t.Fatalf("Expected a body with 'golang' got '%s'", body)
	}
}

func SomeHandler(w http.ResponseWriter, r *http.Request) {
	//	params := context.Get(r, "path_params")
	fmt.Fprint(w, "golang")
}
