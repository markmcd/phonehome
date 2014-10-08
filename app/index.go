package main

import (
	"appengine"
)

type Result struct {
	User string // e.g., google
	Repo string // e.g., go-github

	Path string // file path, e.g. "examples/markdown/main.go"
	Line int    // line number hit

	Offset int     // offset within Lines where result is
	Lines []string
}

func Search(c appengine.Context, query string) (results []*Result, err error) {
	// TODO: everything

	if query != "test" {
		// dummy return nothing
		return nil, nil
	}

	// with "test", return...
	result := &Result{
		User: "google",
		Repo: "go-github",
		Path: "examples/markdown/main.go",
		Line: 100,
		Offset: 1,
		Lines: []string{
			"const x = \"whatever\"",
			"# your result test is awesome",
			"const y = \"other stuff\"",
		},
	}

	results = append(results, result)
	return
}