package textfmt

import "strings"

type Block struct {
	Lines []*Line
}

func NewBlock(input string) Block {
	trimmed := strings.TrimSpace(input)
	normalized := collapseBreakableSpaces(trimmed)
	lines := []*Line{&Line{Text: normalized}}
	return Block{Lines: lines}
}

func (b *Block) Width() int {
	max := 0
	for _, l := range b.Lines {
		if w := l.Width(); w > max {
			max = w
		}
	}
	return max
}

func (b *Block) LineCount() int {
	return len(b.Lines)
}

func (b *Block) SetLineCount(n int) {
	for len(b.Lines) < n {
		b.Lines = append(b.Lines, &Line{Text: ""})
	}
	b.Lines = b.Lines[:n]
}

func (b *Block) StripColors() {
	for _, l := range b.Lines {
		l.StripColors()
	}
}

func (b *Block) Pad(width int, align Align, padChar rune) {
	for _, l := range b.Lines {
		l.Pad(width, align, padChar)
	}
}

func (b *Block) ExpandTabs(tabStop int) {
	for _, l := range b.Lines {
		l.ExpandTabs(tabStop)
	}
}

func (b *Block) WordWrap(limit int) {
	var out []*Line
	for _, l := range b.Lines {
		out = append(out, l.WordWrap(limit)...)
	}
	b.Lines = out
}

func (b *Block) TrimEmptyLines() {
	for len(b.Lines) > 0 && b.Lines[0].Width() == 0 {
		b.Lines = b.Lines[1:]
	}
	for len(b.Lines) > 0 && b.Lines[len(b.Lines)-1].Width() == 0 {
		b.Lines = b.Lines[:len(b.Lines)-1]
	}
}
