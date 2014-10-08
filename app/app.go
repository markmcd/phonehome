package main

import (
	"fmt"
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

func init() {
	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		client := urlfetch.Client(c)

		user := r.FormValue("user")
		repo := r.FormValue("repo")
		if user == "" || repo == "" {
			c.Debugf("needs user/repo set")
			w.WriteHeader(400)
			return
		}

		names, err := Fetch(client, user, repo)
		if err != nil {
			c.Warningf("couldn't fetch %s/%s: %s", user, repo, err)
			w.WriteHeader(500)
			return
		}

		// success-ish
		fmt.Fprintf(w, "found %d files in %s/%s", len(names), user, repo)
		for _, name := range names {
			c.Infof("got file: %s", name)
		}
	})
}
