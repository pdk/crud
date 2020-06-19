package crud

import (
	"fmt"
	"strconv"
	"strings"
)

// markers returns a string of n bind markers, comma separated
func markers(bindStyle BindStyle, first, count int) string {

	sb := strings.Builder{}

	for i := 0; i < count; i++ {

		if i > 0 {
			sb.WriteString(", ")
		}

		sb.WriteString(marker(bindStyle, i+first))
	}

	return sb.String()
}

// marker returns a single bind marker of the desired style
func marker(bindStyle BindStyle, pos int) string {

	switch bindStyle {
	case QuestionMark:
		return "?"
	case DollarNumber:
		return "$" + strconv.Itoa(pos)
	case ColonName:
		return ":p" + strconv.Itoa(pos)
	}

	// this should never happen, obviously
	return fmt.Sprintf("**unknown bindStyle %v**", bindStyle)
}
