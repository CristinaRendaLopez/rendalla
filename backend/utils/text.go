package utils

import (
	"strings"
)

func Normalize(input string) string {
	output := strings.ToLower(removeAccents(input))
	return output
}

func removeAccents(s string) string {
	replacer := strings.NewReplacer(
		"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u",
		"Á", "a", "É", "e", "Í", "i", "Ó", "o", "Ú", "u",
		"ñ", "n", "Ñ", "n",
	)
	return replacer.Replace(s)
}
