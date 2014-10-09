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

	// drop leading '.' from any result (e.g. .foobar => foobar)...
	for i, f := range fields {
		if f[0] == '.' {
			f = f[1:]
			fields[i] = f
		}
	}

	return fields
}

func Tags(raw []byte) (out []*Token) {
	mapping := make(map[string]map[int]bool)

	lines := bytes.Split(raw, []byte("\n"))
	for no, line := range lines {

		// TODO: language-specific stuff

		for _, f := range Fields(line) {
			s := string(f)
			lines, ok := mapping[s]
			if !ok {
				lines = make(map[int]bool)
				mapping[s] = lines
			}
			lines[1+no] = true
		}
	}

	for term, linemap := range mapping {
		lines := make([]int, 0, len(linemap))
		for no, _ := range linemap {
			lines = append(lines, no)
		}
		t := &Token{Line: lines, Term: term}
		out = append(out, t)
	}

	return
}
