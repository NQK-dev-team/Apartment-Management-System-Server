package utils

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func RemoveDiacritics(str string) string {
	// Create a transformer that converts to NFD and removes combining characters.
	// runes.Remove is used to remove all runes in the unicode.Mn category.
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)))

	// Transform the input string.
	normalizedStr, _, _ := transform.String(t, str)

	return normalizedStr
}

func CompareStringRaw(str1, str2 string) bool {
	return strings.ToLower(RemoveDiacritics(str1)) == strings.ToLower(RemoveDiacritics(str2))
}
