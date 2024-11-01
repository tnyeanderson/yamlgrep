package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type filterer struct {
	writer   io.Writer
	grepArgs []string
}

// filter provides doc to the system's grep program, and writes doc to f.writer
// if the grep is successful.
func (f *filterer) filter(doc []byte) error {
	if len(doc) == 0 {
		return nil
	}
	cmd := exec.Command("grep", f.grepArgs...)
	cmd.Stdin = bytes.NewReader(doc)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err == nil {
		// Found
		if _, err := f.writer.Write(doc); err != nil {
			return err
		}
	}
	return nil
}

// lex reads from an io.Reader which contains a YAML multidoc, and runs
// callback with each document. The multidoc separator is included in the next
// document, if present.
func lex(r io.Reader, callback func([]byte) error) error {
	sep := []byte("\n---")
	prefix := sep[0]
	suffix := sep[1:]
	doc := []byte{}
	br := bufio.NewReader(r)
	for {
		// Read until the first character in sep
		b, err := br.ReadBytes(prefix)
		if len(b) >= 3 && len(doc) > 0 {
			// Check if the rest of the sep matches
			if string(b[0:3]) == string(suffix) {
				if err := callback(doc); err != nil {
					return err
				}
				doc = []byte{}
			}
		}
		doc = append(doc, b...)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return callback(doc)
}

func main() {
	f := &filterer{
		writer:   os.Stdout,
		grepArgs: os.Args[1:],
	}
	if err := lex(os.Stdin, f.filter); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
