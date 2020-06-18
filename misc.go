package crud

import (
	"strconv"
	"strings"
)

// markers returns a string of n bind markers, comma separated
func markers(n int) string {
	sb := strings.Builder{}
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("$")
		sb.WriteString(strconv.Itoa(i + 1))
	}

	return sb.String()
}
