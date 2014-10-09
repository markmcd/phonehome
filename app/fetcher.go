package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	pathutil "path"
	"strings"
	"time"
)

type FetchFile struct {
	Path  string
	When  time.Time
	Bytes []byte
}

// Fetch pulls a zip file from github.
func Fetch(client *http.Client, user, repo, branch string) ([]*FetchFile, error) {
	if branch == "" {
		branch = "master"
	}
	path := fmt.Sprintf("%s/%s", user, repo)
	url := fmt.Sprintf("https://github.com/%s/archive/%s.zip", path, branch)

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
		var bytes []byte

		// Remove expected prefix: "<repo>-<branch>/".
		path := strings.TrimPrefix(file.Name, repo+"-"+branch+"/")

		if shouldReadFile(path, &file.FileHeader) {
			reader, err := file.Open()
			if err != nil {
				// ignore for now, continue
				log.Printf("can't open file from zip: %s", file.Name)
			} else {
				bytes, _ = ioutil.ReadAll(reader)
				reader.Close()

				if !probablyText(bytes) {
					bytes = nil
				}
			}
		}

		ff := &FetchFile{
			Path:  path,
			When:  file.ModTime(),
			Bytes: bytes,
		}
		out = append(out, ff)
	}
	return out, nil
}

// shouldReadFile guesses whether the file should be read/indexed.
func shouldReadFile(path string, fh *zip.FileHeader) bool {
	basename := pathutil.Base(path)
	if len(basename) == 0 {
		return false
	}
	if basename[0] == '.' {
		// e.g., ".gitignore"
		return false
	}

	ext := pathutil.Ext(basename)
	mimetype := mime.TypeByExtension(ext)
	if strings.HasPrefix(mimetype, "image/") ||
		strings.HasPrefix(mimetype, "audio/") ||
		strings.HasPrefix(mimetype, "video/") {
		return false
	}
	return true
}

// probablyText returns whether the raw input is probably ASCII text.
func probablyText(raw []byte) bool {
	var valid int
	for _, char := range raw {
		if char >= 32 && char < 128 {
			valid++
		}
		if valid >= 256 {
			return true
		}
	}

	if valid >= len(raw)/4 {
		return true
	}
	return false
}
