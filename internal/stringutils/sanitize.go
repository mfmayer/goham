package stringutils

import (
	"strings"
	"unicode"
)

// Sanitize replaces all characters in a string that are not A-Z, a-z, or 0-9 with underscores
func Sanitize(s string) string {
	var result strings.Builder
	for _, c := range s {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			result.WriteRune(c)
		} else {
			result.WriteRune('_')
		}
	}
	return result.String()
}
