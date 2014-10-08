package main

import (
	"fmt"
	"time"

	"appengine"
	"appengine/datastore"
)

// Repo is an indexed repository.
type Repo struct {
	User    string
	Repo    string
	Branch  string
	Updated time.Time
}

func (r *Repo) ID() string {
	if r.Branch != "master" && r.Branch != "" {
		return fmt.Sprintf("%s/%s/%s", r.User, r.Repo, r.Branch)
	}
	return fmt.Sprintf("%s/%s", r.User, r.Repo)
}

// Key gets the top-level (no parent) key for Repo.
func (r *Repo) Key(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Repo", r.ID(), 0, nil)
}

// File is a file within a Repo.
type File struct {
	Path string
	When time.Time
}

func (f *File) Key(c appengine.Context, r *Repo) *datastore.Key {
	return datastore.NewKey(c, "File", f.Path, 0, r.Key(c))
}

// Token is an indexed term. It has no obvious key, but its parent should be a
// file Key and contain its path.
type Token struct {
	Line []int
	Term string
}
