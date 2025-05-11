package textfmt

import (
	"regexp"
	"unicode"
)

var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)

type Align int

const (
	Left Align = iota
	Center
	Right
)

type ColorControl int

const (
	NoColor ColorControl = iota
	AllowColor
)

type Line struct {
	Text  string
	width int
}

func (f *Line) Width(tabStop int) int {
	if f.width > 0 {
		return f.width
	}

	scanner := newScanner(f.Text)
	width := 0

	for {
		tok := scanner.nextToken()
		if tok.Type == tokenEOF {
			break
		}
		switch tok.Type {
		case tokenTab:
			nextTab := ((width / tabStop) + 1) * tabStop
			width += nextTab - width
		default:
			width += tok.Width
		}
	}

	f.width = width
	return width
}

func (f *Line) StripColors() {
	f.Text = ansiEscapeRe.ReplaceAllString(f.Text, "")
	f.width = 0
}

func (f *Line) WordWrap(width int, tabStop int, align Align, color ColorControl, padChar rune) []string {
	spec := NewWrapSpec(width, tabStop, align, color, padChar)
	return spec.WordWrap(f.Text)
}

func isWideRune(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hangul, unicode.Hiragana, unicode.Katakana)
}
