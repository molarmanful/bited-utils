package bitedutils

import (
	"unicode"

	"github.com/mattn/go-runewidth"
)

// Check panics if error is non-nil.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func WcWidth(r rune, nerd bool) int {
	if nerd && unicode.Is(NerdFont, r) {
		return 2
	}
	return runewidth.RuneWidth(r)
}
