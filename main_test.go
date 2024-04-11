package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var content []byte

func init() {
	b, err := os.ReadFile("testdata/mock.yaml")
	if err != nil {
		panic("oops")
	}
	content = b
}

func BenchmarkYAMLGrep(b *testing.B) {
	r := bytes.NewReader(content)
	f := &filterer{
		writer:   io.Discard,
		grepArgs: []string{"Male"},
	}
	b.ResetTimer()
	if err := lex(r, f.filter); err != nil {
		b.Fatal(err)
	}
}

func BenchmarkYAMLGrep_Bash(b *testing.B) {
	r := bytes.NewReader(content)
	b.ResetTimer()
	cmd := exec.Command("yamlgrep", "Male")
	cmd.Stdin = r
	cmd.Stdout = io.Discard
	b.ResetTimer()
	cmd.Run()
}

func TestYAMLGrep(t *testing.T) {
	r := bytes.NewReader(content)
	buf := &bytes.Buffer{}
	f := &filterer{
		writer:   buf,
		grepArgs: []string{"Male"},
	}
	if err := lex(r, f.filter); err != nil {
		t.Fatal(err)
	}
	expectedLines := 456 * 7
	lines := bytes.Count(buf.Bytes(), []byte("\n"))
	if lines != expectedLines {
		t.Fatalf("incorrect line count, got: %d, expected %d", lines, expectedLines)
	}
}

func TestYAMLGrep_CaseInsensitive(t *testing.T) {
	r := bytes.NewReader(content)
	buf := &bytes.Buffer{}
	f := &filterer{
		writer:   buf,
		grepArgs: []string{"-i", "Male"},
	}
	if err := lex(r, f.filter); err != nil {
		t.Fatal(err)
	}
	expectedLines := 887 * 7
	lines := bytes.Count(buf.Bytes(), []byte("\n"))
	if lines != expectedLines {
		t.Fatalf("incorrect line count, got: %d, expected %d", lines, expectedLines)
	}
}

func TestLex_DocCounts(t *testing.T) {
	tests := []string{
		// 0 Empty
		``,

		// 1 Empty with leading delim
		`---`,

		// 2 Single with leading delim
		`---
var: 1
`,

		// 3 Single without delim
		`var: 1`,

		// 4 Single with trailing delim
		`var: 1
---
`,

		// 5 Multi with leading delim
		`---
var: 1
---
var: 2
---
var: 3
`,

		// 6 Multi without leading delim
		`var: 1
---
var: 2
---
var: 3
`,

		// 7 Multi with leading and trailing delim
		`---
var: 1
---
var: 2
---
var: 3
---
`,

		// 8 Multi without leading delim
		`---
var: 1
---
var: 2
---
---
var: 3
`,

		// 9 mock.yaml
		string(content),
	}

	expected := []int{1, 1, 1, 1, 2, 3, 3, 4, 4, 1000}

	for i, f := range tests {
		c := 0
		r := strings.NewReader(f)
		lex(r, func(b []byte) error { c++; return nil })
		e := expected[i]
		if c != e {
			t.Fatalf("incorrect doc count from lex at index %d, got: %d, expected %d", i, c, e)
		}
	}
}
