package pkg

import (
	"net/http"
	"testing"
)

func TestInsert(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {}
	testCases := []struct {
		name        string
		keys        []string
		verb        string
		expectError bool
	}{
		{"/", []string{}, "GET", false},
		{"/a", []string{"a"}, "GET", false},
		{"/a", []string{"a"}, "GET", false},
		{"/a", []string{"a"}, "POST", false},
		{"/a", []string{"a"}, "PUT", false},
		{"/a", []string{"a"}, "DELETE", false},
		{"/a", []string{"a"}, "PATCH", false},
		{"/:a", []string{":a"}, "GET", false},
		{"/:b", []string{":b"}, "GET", true},
		{"/:a/f", []string{":a", "f"}, "GET", false},
		{"/a/b", []string{"a", "b"}, "GET", false},
		{"/d/e", []string{"d", "e"}, "GET", false},
	}
	n := newNode()
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := n.insert(testCase.keys, testCase.verb, handler)
			if _, actualError := err.(error); actualError != testCase.expectError {
				t.Errorf("Error, expected error to be %t, got %t", testCase.expectError, actualError)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {}
	testCases := []struct {
		name             string
		insertVerb       string
		insertKeys       []string
		searchVerb       string
		searchKeys       []string
		expectedHandlers int
	}{
		{"/", "GET", []string{}, "GET", []string{}, 1},
		{"/", "GET", []string{"a"}, "GET", []string{"a"}, 1},
		{"/", "GET", []string{"a"}, "GET", []string{"b"}, 0},
		{"/", "GET", []string{":a/c"}, "GET", []string{"1/c"}, 1},
	}
	n := newNode()
	m := make(map[string]string)
	f := func(k, v string) {
		m[k] = v
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			n.insert(testCase.insertKeys, testCase.insertVerb, handler)
			handlers, ok := n.search(testCase.searchKeys, testCase.searchVerb, f)
			actualHandlers := 0
			if ok {
				actualHandlers = len(handlers)
			}
			if actualHandlers != testCase.expectedHandlers {
				t.Errorf(
					"Error, expected to get %d handlers got %d",
					testCase.expectedHandlers,
					actualHandlers,
				)
			}
		})
	}
}
