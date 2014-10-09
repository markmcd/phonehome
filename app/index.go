package main

import (
	"bytes"

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

	fkey *datastore.Key
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
		lines := bytes.Split(ff.Bytes, []byte("\n"))

		fkey := datastore.NewKey(c, "File", ff.Path, 0, rkey)
		tkey := datastore.NewIncompleteKey(c, "Token", fkey)
		tokens := Tags(lines)
		if len(tokens) == 0 {
			stats.EmptyFiles++
			continue // don't store this file
		}

		for _, token := range tokens {
			ents = append(ents, token)
			keys = append(keys, tkey)
		}
		stats.Tokens += len(tokens)

		ents = append(ents, &File{When: ff.When, Lines: lines})
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

	files := make(map[string]*File)
	var fkeys []*datastore.Key

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

			path := fkey.StringID()
			if _, ok := files[fkey.String()]; !ok {
				files[fkey.String()] = nil
				fkeys = append(fkeys, fkey)
			}

			repo := BuildRepo(rkey.StringID())
			for _, line := range token.Lines {
				result := &Result{
					fkey: fkey,
					User: repo.User,
					Repo: repo.Repo,
					Branch: repo.Branch,
					Path: path,
					Line: 1+line,
				}
				results = append(results, result)
			}
		}
	}

	// Load context for files.
	out := make([]File, len(fkeys))
	err = datastore.GetMulti(c, fkeys, out)
	if err != nil {
		c.Infof("can't get files: keys=%v", fkeys)
		return nil, err
	}
	for i, _ := range out {
		fkey := fkeys[i]
		files[fkey.String()] = &out[i]
	}
	for _, result := range results {
		file := files[result.fkey.String()]
		if file == nil {
			c.Warningf("couldn't load file key: %s", result.fkey.String())
			continue
		}
		lineno := result.Line - 1  // actual line

		if len(file.Lines) < lineno {
			c.Warningf("file has %d lines, wanted %d: %s", len(file.Lines), lineno, result.fkey.String())
			continue
		}

		// TODO: line context
		line := string(file.Lines[lineno])
		result.Lines = make([]string, 1)
		result.Lines[0] = line
	}

	return
}
