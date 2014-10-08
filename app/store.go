package main

import (
	"fmt"
	"time"
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

// File is a file within a Repo. Its key is its path.
type File struct {
	When time.Time
}

// Token is an indexed term. It has no obvious key, but its parent should be a
// file Key and contain its path.
type Token struct {
	Line []int
	Term string
}
