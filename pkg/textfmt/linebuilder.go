package textfmt

import (
	"strings"
)

type lineBuilder struct {
	line     strings.Builder
	width    int
	sgr      *sgrState
	wrapSpec *WrapSpec
}

func newLineBuilder(sgr *sgrState, wrapSpec *WrapSpec) *lineBuilder {
	return &lineBuilder{
		sgr:      sgr.copy(),
		wrapSpec: wrapSpec,
	}
}

func (b *lineBuilder) reset() {
	b.line.Reset()
	b.width = 0
	if b.wrapSpec.Color == AllowColor && len(b.sgr.attrs) > 0 {
		b.line.WriteString(b.sgr.string())
	}
}

func (b *lineBuilder) flushAsString() string {
	text := b.line.String()
	if b.wrapSpec.Color == AllowColor && len(b.sgr.attrs) > 0 {
		text += "\x1b[0m"
	}
	result := applyAlignment(text, b.wrapSpec.Width, b.wrapSpec.Align, b.wrapSpec.PadChar)
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

func (b *lineBuilder) writeTab() {
	spaces := ((b.width/b.wrapSpec.TabStop)+1)*b.wrapSpec.TabStop - b.width
	b.line.WriteString(strings.Repeat(" ", spaces))
	b.width += spaces
}

func (b *lineBuilder) canFit(w int) bool {
	return b.width+w <= b.wrapSpec.Width
}

func (b *lineBuilder) remaining() int {
	return b.wrapSpec.Width - b.width
}

func (b *lineBuilder) len() int {
	return b.line.Len()
}
