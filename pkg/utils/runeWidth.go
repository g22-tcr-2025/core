package utils

import "unicode"

func isWideRune(r rune) bool {
	if unicode.In(r,
		unicode.Han,
		unicode.Hangul,
		unicode.Hiragana,
		unicode.Katakana,
	) || (r >= 0x1F300 && r <= 0x1FAFF) {
		return true
	}
	return false
}

func StringDisplayWidth(s string) int {
	width := 0
	for _, r := range s {
		if isWideRune(r) {
			width += 2
		} else {
			width += 1
		}
	}
	return width
}
