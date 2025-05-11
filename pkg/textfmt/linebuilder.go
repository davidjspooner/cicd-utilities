package textfmt

import (
	"strings"
)

type lineBuilder struct {
	line   strings.Builder
	width  int
	sgr    *sgrState
	color  ColorControl
	align  Align
	pad    rune
	target int
}

func newLineBuilder(sgr *sgrState, color ColorControl, align Align, pad rune, targetWidth int) *lineBuilder {
	return &lineBuilder{
		sgr:    sgr.copy(),
		color:  color,
		align:  align,
		pad:    pad,
		target: targetWidth,
	}
}

func (b *lineBuilder) reset() {
	b.line.Reset()
	b.width = 0
	if b.color == AllowColor && len(b.sgr.attrs) > 0 {
		b.line.WriteString(b.sgr.string())
	}
}

func (b *lineBuilder) flushTo(lines *[]*Line) {
	text := b.line.String()
	if b.color == AllowColor && len(b.sgr.attrs) > 0 {
		text += "\x1b[0m"
	}
	*lines = append(*lines, &Line{
		Text:  applyAlignment(text, b.target, b.align, b.pad),
		width: 0,
	})
	b.reset()
}

func (b *lineBuilder) flushAsString() string {
	text := b.line.String()
	if b.color == AllowColor && len(b.sgr.attrs) > 0 {
		text += "\x1b[0m"
	}
	result := applyAlignment(text, b.target, b.align, b.pad)
	b.reset()
	return result
}

func (b *lineBuilder) writeString(s string, w int) {
	b.line.WriteString(s)
	b.width += w
}

func (b *lineBuilder) writeRune(r rune, w int) {
	b.line.WriteRune(r)
	b.width += w
}

func (b *lineBuilder) writeSpace() {
	b.line.WriteRune(' ')
	b.width++
}

func (b *lineBuilder) writeTab(tabStop int) {
	spaces := ((b.width/tabStop)+1)*tabStop - b.width
	b.line.WriteString(strings.Repeat(" ", spaces))
	b.width += spaces
}

func (b *lineBuilder) canFit(w int) bool {
	return b.width+w <= b.target
}

func (b *lineBuilder) remaining() int {
	return b.target - b.width
}

func (b *lineBuilder) len() int {
	return b.line.Len()
}
