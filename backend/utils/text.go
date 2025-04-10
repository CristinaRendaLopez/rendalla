package utils

import (
	"strings"
)

// Normalize returns a lowercase, accent-free version of the input string.
// Useful for performing normalized search or comparisons.
func Normalize(input string) string {
	output := strings.ToLower(removeAccents(input))
	return output
}

// removeAccents replaces common accented characters with their non-accented counterparts.
// Supports Spanish vowels and the 'ñ' character.
func removeAccents(s string) string {
	replacer := strings.NewReplacer(
		"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u",
		"Á", "a", "É", "e", "Í", "i", "Ó", "o", "Ú", "u",
		"ñ", "n", "Ñ", "n",
	)
	return replacer.Replace(s)
}
