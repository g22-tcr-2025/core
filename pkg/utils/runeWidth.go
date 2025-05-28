package utils

import (
	"github.com/mattn/go-runewidth"
)

func StringDisplayWidth(s string) int {
	return runewidth.StringWidth(s)
}
