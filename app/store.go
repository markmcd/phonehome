package main

import (
	"fmt"
	"strings"
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

func BuildRepo(raw string) *Repo {
	parts := strings.Split(raw, "/")
	if len(parts) < 2 {
		return nil
	}
	r := &Repo{}
	r.User, r.Repo = parts[0], parts[1]
	if len(parts) >= 3 {
		r.Branch = parts[2]
	}
	return r
}

// File is a file within a Repo. Its key is its path.
type File struct {
	When time.Time
	Lines [][]byte
}

// Token is an indexed term. It has no obvious key, but its parent should be a
// file Key and contain its path.
type Token struct {
	Lines []int // stores lines starting with 0
	Term string
}
