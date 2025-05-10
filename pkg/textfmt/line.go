package textfmt

import (
	"regexp"
	"strings"
	"unicode"
)

var ansiEscapeRe = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)

type Align int

const (
	Left Align = iota
	Center
	Right
)

type Line struct {
	Text  string
	width int
}

// Width returns the display width of the line, excluding ANSI sequences.
func (f *Line) Width() int {
	if f.width > 0 {
		return f.width
	}
	parts := ansiEscapeRe.Split(f.Text, -1)
	width := 0
	for _, segment := range parts {
		for _, r := range segment {
			if r == 0 || unicode.Is(unicode.Mn, r) {
				continue
			}
			if isWideRune(r) {
				width += 2
			} else {
				width += 1
			}
		}
	}
	f.width = width
	return width
}

func (f *Line) StripColors() {
	f.Text = ansiEscapeRe.ReplaceAllString(f.Text, "")
	f.width = 0
}

func (f *Line) Pad(width int, align Align, padChar rune) {
	actual := f.Width()
	if actual >= width {
		return
	}
	padding := width - actual
	var left, right int
	switch align {
	case Center:
		left = padding / 2
		right = padding - left
	case Right:
		left = padding
	default:
		right = padding
	}
	f.Text = strings.Repeat(string(padChar), left) + f.Text + strings.Repeat(string(padChar), right)
	f.width = width
}

func (f *Line) ExpandTabs(tabStop int) {
	f.Text = expandTabsInternal(f.Text, tabStop)
	f.width = 0
}

func (f *Line) WordWrap(limit int) []*Line {
	normalized := collapseBreakableSpaces(strings.TrimSpace(f.Text))
	words := strings.Fields(normalized)

	var lines []*Line
	var current string
	for _, word := range words {
		if len(current)+len(word)+1 > limit && current != "" {
			lines = append(lines, &Line{Text: current})
			current = word
		} else {
			if current != "" {
				current += " "
			}
			current += word
		}
	}
	if current != "" {
		lines = append(lines, &Line{Text: current})
	}
	return lines
}

// collapseBreakableSpaces replaces any run of Unicode-defined whitespace with a single ASCII space.
func collapseBreakableSpaces(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	spaceRun := false
	for _, r := range s {
		if unicode.Is(unicode.White_Space, r) {
			if !spaceRun {
				b.WriteRune(' ')
				spaceRun = true
			}
		} else {
			b.WriteRune(r)
			spaceRun = false
		}
	}

	return b.String()
}

// isWideRune returns true for runes that should be treated as width 2 when displayed.
func isWideRune(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hangul, unicode.Hiragana, unicode.Katakana)
}

func expandTabsInternal(s string, tabStop int) string {
	var b strings.Builder
	col := 0
	for _, r := range s {
		if r == '\t' {
			spaces := tabStop - (col % tabStop)
			b.WriteString(strings.Repeat(" ", spaces))
			col += spaces
		} else {
			b.WriteRune(r)
			if unicode.Is(unicode.Mn, r) {
				continue
			} else if isWideRune(r) {
				col += 2
			} else {
				col += 1
			}
		}
	}
	return b.String()
}
