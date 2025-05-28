package utils

import (
	"regexp"

	"github.com/mattn/go-runewidth"
)

var ansiRegax = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiRegax.ReplaceAllString(s, "")
}

func StringDisplayWidth(s string) int {
	return runewidth.StringWidth(stripANSI(s))
}
