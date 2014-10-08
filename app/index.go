package main

import (
	"strings"

	"appengine"
	"appengine/datastore"
)

type Result struct {
	User string // e.g., google
	Repo string // e.g., go-github
	Branch string // e.g., 'gh-pages', blank is assumed master

	Path string // file path, e.g. "examples/markdown/main.go"
	Line int    // line number hit

	Offset int // offset within Lines where result is
	Lines  []string
}

func Index(c appengine.Context, repo *Repo, ffs []*FetchFile) error {
	stats := struct {
		Files, EmptyFiles, Tokens int
	}{}
	c.Infof("indexing %s", repo.ID())

	rkey := datastore.NewKey(c, "Repo", repo.ID(), 0, nil)

	var ents []interface{}
	var keys []*datastore.Key
	for _, ff := range ffs {
		fkey := datastore.NewKey(c, "File", ff.Path, 0, rkey)
		tkey := datastore.NewIncompleteKey(c, "Token", fkey)
		tokens := Tags(ff.Bytes)
		if len(tokens) == 0 {
			stats.EmptyFiles++
			continue // don't store this file
		}

		for _, token := range tokens {
			ents = append(ents, token)
			keys = append(keys, tkey)
		}
		stats.Tokens += len(tokens)

		ents = append(ents, &File{When: ff.When})
		keys = append(keys, fkey)
		stats.Files++
	}
	ents = append(ents, repo)
	keys = append(keys, rkey)
	c.Infof("generated %d ents for %s", len(keys), repo.ID())

	return datastore.RunInTransaction(c, func(c appengine.Context) error {
		// 1: delete everything
		q := datastore.NewQuery("").KeysOnly().Ancestor(rkey)
		oldKeys, err := q.GetAll(c, nil)
		if err != nil {
			return err
		}
		if err = PartDeleteMulti(c, oldKeys); err != nil {
			return err
		}

		// 2: insert everything
		if err = PartPutMulti(c, keys, ents); err != nil {
			return err
		}

		// success!
		c.Infof("updated %s (%v)", repo.ID(), stats)
		return nil
	}, nil)
}

func Search(c appengine.Context, query string) (results []*Result, err error) {
	fields := Fields([]byte(query))
	if len(fields) == 0 {
		return nil, nil
	}

	for _, field := range fields {
		s := string(field)
		q := datastore.NewQuery("Token").
			Filter("Term >=", s).
			Filter("Term <", s+"\uffffd")
		var out []*Token
		keys, err := q.GetAll(c, &out)
		if err != nil {
			return nil, err
		}

		// For now, just add everything.

		for i, tkey := range keys {
			token := out[i]
			fkey := tkey.Parent() // file key
			rkey := fkey.Parent() // repo key

			parts := strings.Split(rkey.StringID(), "/")
			user, repo := parts[0], parts[1]
			var branch string
			if len(parts) >= 3 {
				branch = parts[2]
			}

			for _, line := range token.Line {
				result := &Result{
					User: user,
					Repo: repo,
					Branch: branch,
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
