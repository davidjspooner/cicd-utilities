package main

import (
	"path/filepath"
	"strings"
)

func splitLines(output string) []string {
	return strings.Split(strings.TrimSpace(output), "\n")
}

func globFiles(patterns []string) ([]string, error) {
	var files []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}
		files = append(files, matches...)
	}
	return files, nil
}
