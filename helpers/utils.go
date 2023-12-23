package helpers

import (
	"unicode"
)

func RuneIsAlpha(r rune) bool {
	return unicode.In(r, unicode.Latin)
}

func RuneToAlpha(r rune) string {
	return string([]rune{r, -1})
}
