package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Fetch pulls a zip file from github.
func Fetch(client *http.Client, user, repo string) ([]string, error) {
	path := fmt.Sprintf("%s/%s", user, repo)
	url := fmt.Sprintf("https://github.com/%s/archive/master.zip", path)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewReader(raw)
	r, err := zip.NewReader(buf, int64(len(raw)))
	if err != nil {
		return nil, err
	}

	var out []string
	for _, file := range r.File {
		out = append(out, file.Name)
	}
	return out, nil
}
