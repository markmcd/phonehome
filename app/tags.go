package main

import (
	"bytes"
	"unicode"
)

func Fields(raw []byte) [][]byte {
	raw = bytes.TrimSpace(bytes.ToLower(raw))
	fields := bytes.FieldsFunc(raw, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '.'
	})

	for i, f := range fields {
		// drop leading '.' from any result (e.g. .foobar => foobar)...
		if f[0] == '.' {
			f = f[1:]
			fields[i] = f
		}
		// limit length (appengine max is 500)
		if len(f) > 128 {
			f = f[:128]
			fields[i] = f
		}
	}

	return fields
}

func Tags(lines [][]byte) (out []*Token) {
	mapping := make(map[string]map[int]bool)

	for no, line := range lines {

		// TODO: language-specific stuff

		for _, f := range Fields(line) {
			s := string(f)
			lines, ok := mapping[s]
			if !ok {
				lines = make(map[int]bool)
				mapping[s] = lines
			}
			lines[no] = true
		}
	}

	for term, linemap := range mapping {
		lines := make([]int, 0, len(linemap))
		for no, _ := range linemap {
			lines = append(lines, no)
		}
		t := &Token{Lines: lines, Term: term}
		out = append(out, t)
	}

	return
}
