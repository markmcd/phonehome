package main

import (
	"strings"

	"appengine"
	"appengine/datastore"
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

	fields := Fields([]byte(query))
	if len(fields) == 0 {
		return  nil, nil
	}

	if query != "test" {
		// dummy return nothing
		return nil, nil
	}

	for _, field := range fields {
		s := string(field)
		q := datastore.NewQuery("Token").
				Filter("Term >", s).
				Filter("Term <", s + "\uffffd")
		var out []*Token
		keys, _ := q.GetAll(c, out)

		// For now, just add everything.

		for i, tkey := range keys {
			token := out[i]
			fkey := tkey.Parent() // file key
			rkey := fkey.Parent() // repo key

			parts := strings.Split(rkey.StringID(), "/")

			for _, line := range token.Line {
				result := &Result{
					User: parts[0],
					Repo: parts[1],
					Path: fkey.StringID(),
					Line: line,
					// TODO: offset/lines
				}
				results = append(results, result)
			}
		}
	}

	return
}