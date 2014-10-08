package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type FetchFile struct {
	Path  string
	When  time.Time
	Bytes []byte
}

// Fetch pulls a zip file from github.
func Fetch(client *http.Client, user, repo string) ([]*FetchFile, error) {
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

	var out []*FetchFile
	for _, file := range r.File {
		// TODO: we probably don't want to read large-ish binary files
		reader, err := file.Open()
		if err != nil {
			// ignore for now, continue
			log.Printf("can't open file from zip: %s", file.Name)
			continue
		}
		bytes, _ := ioutil.ReadAll(reader)
		reader.Close()

		ff := &FetchFile{
			Path:  file.Name,
			When:  file.ModTime(),
			Bytes: bytes,
		}
		out = append(out, ff)
	}
	return out, nil
}
