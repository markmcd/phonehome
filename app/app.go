package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"net/url"

	"appengine"
	"appengine/taskqueue"
	"appengine/urlfetch"
)

func init() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		request := struct{
			Repo struct{
				Name string `json:"name"`
			} `json:"repository"`
		}{}

		d := json.NewDecoder(r.Body)
		err := d.Decode(&request)
		if err != nil {
			c.Warningf("couldn't decode webhook: %v", err)
			w.WriteHeader(500)
			return
		}

		repo := BuildRepo(request.Repo.Name)
		if repo == nil {
			c.Warningf("couldn't build repo from name: %s", request.Repo.Name)
			w.WriteHeader(400)
			return
		}
		t := taskqueue.NewPOSTTask("/index", url.Values{
			"user": {repo.User},
			"repo": {repo.Repo},
			// TODO: branch
		})
		_, err = taskqueue.Add(c, t, "")
		if err != nil {
			c.Warningf("couldn't enqueue task: %v", err)
			w.WriteHeader(500)
			return
		}
	})
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
