package textfmt

import (
	"strings"
	"unicode"
)

type Block struct {
	Lines []*Line
}

func NewBlock(input string) *Block {
	trimmed := strings.TrimSpace(input)
	normalized := collapseBreakableSpaces(trimmed)
	lines := []*Line{&Line{Text: normalized}}
	return &Block{Lines: lines}
}

func (b *Block) Width(tabStop int) int {
	w := 0
	for _, l := range b.Lines {
		w = max(w, l.Width(tabStop))
	}
	return w
}

func (b *Block) WordWrap(width int, tabStop int, align Align, color ColorControl, padChar rune) []string {
	if padChar == 0 {
		padChar = ' '
	}
	var wrapped []string
	for _, l := range b.Lines {
		wrapped = append(wrapped, l.WordWrap(width, tabStop, align, color, padChar)...)
	}
	return wrapped
}

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
