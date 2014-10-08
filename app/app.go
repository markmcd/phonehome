package main

import (
	"fmt"
	"html"
	"net/http"
	// "strings"
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		//
		// q := r.FormValue("q")
		// re, err := regexp.Compile(q)
		// if err != nil {
		// 	w.WriteHeader(500)
		// 	return
		// }
		//
		// r := strings.NewReader("import test\nx = 1\nprint(\"hello %d\" % x)")
		//
		// var g regexp.Grep
		// g.Regexp = re
		// g.Reader(r, "test.py")
	})
}