package main

import (
	"strings"
)

func splitLines(output string) []string {
	return strings.Split(strings.TrimSpace(output), "\n")
}
