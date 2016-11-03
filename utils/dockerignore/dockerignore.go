package dockerignore

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ReadAll reads a .dockerignore file and returns the list of file patterns
// to ignore. Note this will trim whitespace from each line as well
// as use GO's "clean" func to get the shortest/cleanest path for each.
// A little copy from the analogous package at github.com/docker/docker/builder/dockerignore
func ReadAll() ([]string, error) {
	bytesrece, err := ioutil.ReadFile(".dockerignore")
	if err != nil {
		return []string{}, nil
	}
	r := bytes.NewReader(bytesrece)
	if r == nil {
		return nil, nil
	}
	scanner := bufio.NewScanner(r)
	var excludes []string
	currentLine := 0

	utf8bom := []byte{0xEF, 0xBB, 0xBF}
	for scanner.Scan() {
		scannedBytes := scanner.Bytes()
		// We trim UTF8 BOM
		if currentLine == 0 {
			scannedBytes = bytes.TrimPrefix(scannedBytes, utf8bom)
		}
		pattern := string(scannedBytes)
		currentLine++
		// Lines starting with # (comments) are ignored before processing
		if strings.HasPrefix(pattern, "#") {
			continue
		}
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		pattern = filepath.Clean(pattern)
		pattern = filepath.ToSlash(pattern)
		excludes = append(excludes, pattern)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading .dockerignore: %v", err)
	}
	return excludes, nil
}

// Difference subtracts a slice of strings from another
func Difference(context []string, ignored []string) []string {
	for _, singleIgnoredFile := range ignored {
		if index, ok := contains(context, singleIgnoredFile); ok {
			context = append(context[:index], context[index+1:]...)
		}
	}
	return context
}

func contains(s []string, e string) (int, bool) {
	for i, a := range s {
		if a == e {
			return i, true
		}
	}
	return 0, false
}
