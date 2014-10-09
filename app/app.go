package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"appengine"
	"appengine/urlfetch"
)

func init() {
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		client := urlfetch.Client(c)

		repo := &Repo{
			User:    r.FormValue("user"),
			Repo:    r.FormValue("repo"),
			Updated: time.Now(),
		}
		if r.FormValue("branch") != "" {
			repo.Branch = r.FormValue("branch")
		}
		if repo.User == "" || repo.Repo == "" {
			c.Debugf("needs user/repo set")
			w.WriteHeader(400)
			return
		}

		ffs, err := Fetch(client, repo.User, repo.Repo, repo.Branch)
		if err != nil {
			c.Warningf("couldn't fetch %s: %s", repo.ID(), err)
			w.WriteHeader(500)
			return
		}

		err = Index(c, repo, ffs)
		if err != nil {
			c.Warningf("couldn't index %s: %s", repo.ID(), err)
			w.WriteHeader(500)
			return
		}

		fmt.Fprintf(w, "updated %s", repo.ID())
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
