package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
)

const (
	OP_LIMIT = 500
)

func init() {
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		client := urlfetch.Client(c)
		stats := struct {
			Files, EmptyFiles, Tokens int
		}{}

		repo := &Repo{
			User:    r.FormValue("user"),
			Repo:    r.FormValue("repo"),
			Updated: time.Now(),
		}
		if repo.User == "" || repo.Repo == "" {
			c.Debugf("needs user/repo set")
			w.WriteHeader(400)
			return
		}

		ffs, err := Fetch(client, repo.User, repo.Repo)
		if err != nil {
			c.Warningf("couldn't fetch %s: %s", repo.ID(), err)
			w.WriteHeader(500)
			return
		}

		var ents []interface{}
		var keys []*datastore.Key
		for _, ff := range ffs {
			file := &File{Path: ff.Path, When: ff.When}
			fkey := file.Key(c, repo)

			// TODO: tokens
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

			ents = append(ents, file)
			keys = append(keys, fkey)
			stats.Files++
		}
		ents = append(ents, repo)
		keys = append(keys, repo.Key(c))

		err = datastore.RunInTransaction(c, func(c appengine.Context) error {
			// 1: delete everything
			rkey := repo.Key(c)
			q := datastore.NewQuery("").KeysOnly().Ancestor(rkey)
			oldKeys, err := q.GetAll(c, nil)
			if err != nil {
				return err
			}
			for len(oldKeys) > 0 {
				l := OP_LIMIT
				if l > len(oldKeys) {
					l = len(oldKeys)
				}
				if err = datastore.DeleteMulti(c, oldKeys[:l]); err != nil {
					return err
				}
				oldKeys = oldKeys[l:]
			}

			// 2: insert everything
			for len(keys) > 0 {
				l := OP_LIMIT
				if l > len(keys) {
					l = len(keys)
				}
				_, err = datastore.PutMulti(c, keys[:l], ents[:l])
				if err != nil {
					return err
				}
				ents = ents[l:]
				keys = keys[l:]
			}

			// success!
			return nil
		}, nil)
		if err != nil {
			c.Warningf("couldn't do update for %s: %s", repo.ID(), err)
			w.WriteHeader(500)
			return
		}
		c.Infof("updated %s (%v)", repo.ID(), stats)
	})
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)

		query := r.FormValue("q")
		if query == "" {
			c.Debugf("no search query")
			w.WriteHeader(400)
			return
		}

		results, err := Search(c, query)
		if err != nil {
			c.Warningf("search failed: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		json, err := json.Marshal(results)
		fmt.Fprintf(w, string(json))
	})
}
