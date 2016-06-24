[![GoDoc](https://godoc.org/github.com/zatiti/router?status.svg)](https://godoc.org/github.com/zatiti/router)
[![Build Status](https://travis-ci.org/zatiti/router.svg?branch=master)](https://travis-ci.org/zatiti/router)
[![codecov](https://codecov.io/gh/zatiti/router/branch/master/graph/badge.svg)](https://codecov.io/gh/zatiti/router)
[![GoDoc](https://godoc.org/github.com/zatiti/router?status.svg)](https://godoc.org/github.com/zatiti/router)
[![Go Report Card](https://goreportcard.com/badge/github.com/zatiti/router)](https://goreportcard.com/report/github.com/zatiti/router)
[![Issue Count](https://codeclimate.com/github/zatiti/router/badges/issue_count.svg)](https://codeclimate.com/github/zatiti/router)

# router
This is a golang HTTP router that implements a trie data structure for optimal performance.

## Example

```go
package main

import (
	"fmt"
	"github.com/zatiti/router"
	"net/http"
)

func main() {
	r := router.New(DefaultHandler)
	r.Get("/test/:test_id", HandlerA, HandlerB, HandlerC)
	http.ListenAndServe(":8080", r)
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}

func HandlerA(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "A")
}

func HandlerB(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "B")
}

func HandlerC(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "C")
}
```
