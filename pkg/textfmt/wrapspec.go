package textfmt

import (
	"math"
	"strings"
	"unicode"
)

type WrapSpec struct {
	Width   int
	TabStop int
	Align   Align
	Color   ColorControl
	PadChar rune
}

func NewWrapSpec(width, tabStop int, align Align, color ColorControl, padChar rune) *WrapSpec {
	return &WrapSpec{
		Width:   width,
		TabStop: tabStop,
		Align:   align,
		Color:   color,
		PadChar: padChar,
	}
}

// WordWrap wraps raw input into aligned, color-aware, padded lines.
func (w *WrapSpec) WordWrap(text string) []string {
	scanner := newScanner(text)
	var lines []string
	var word strings.Builder
	var wordWidth int
	activeSGR := newSGRState()
	lb := newLineBuilder(activeSGR, w.Color, w.Align, w.PadChar, w.Width)

	for {
		tok := scanner.nextToken()
		if tok.Type == tokenEOF {
			break
		}

		switch tok.Type {
		case tokenColor:
			if w.Color == AllowColor {
				codes := parseSGRCodes(tok.Value)
				for _, code := range codes {
					activeSGR.apply(code)
				}
				word.WriteString(tok.Value)
			}
		case tokenTab:
			if word.Len() > 0 {
				if !lb.canFit(wordWidth) {
					lines = append(lines, lb.flushAsString())
				}
				lb.writeString(word.String(), wordWidth)
				word.Reset()
				wordWidth = 0
			}
			lb.writeTab(w.TabStop)
		case tokenWhitespace:
			if word.Len() > 0 {
				if !lb.canFit(wordWidth) {
					lines = append(lines, lb.flushAsString())
				}
				lb.writeString(word.String(), wordWidth)
				word.Reset()
				wordWidth = 0
			}
			if !lb.canFit(1) {
				lines = append(lines, lb.flushAsString())
			}
			lb.writeSpace()
		case tokenOther:
			if tok.Width > w.Width {
				if word.Len() > 0 {
					if !lb.canFit(wordWidth) {
						lines = append(lines, lb.flushAsString())
					}
					lb.writeString(word.String(), wordWidth)
					word.Reset()
					wordWidth = 0
				}
				for _, r := range tok.Value {
					if !lb.canFit(1) {
						lines = append(lines, lb.flushAsString())
					}
					lb.writeRune(r, 1)
				}
			} else {
				word.WriteString(tok.Value)
				wordWidth += tok.Width
			}
		}
	}

	if word.Len() > 0 {
		if !lb.canFit(wordWidth) {
			lines = append(lines, lb.flushAsString())
		}
		lb.writeString(word.String(), wordWidth)
	}

	if lb.len() > 0 {
		lines = append(lines, lb.flushAsString())
	}

	return lines
}

func (s *WrapSpec) normalizeSpec() {
	if s.Width == 0 {
		s.Width = math.MaxInt
	}
	if s.PadChar == 0 {
		s.PadChar = ' '
	}
	if s.TabStop == 0 {
		s.TabStop = 4
	}
}

func applyAlignment(s string, width int, align Align, padChar rune) string {
	plain := ansiEscapeRe.ReplaceAllString(s, "")
	actualWidth := 0
	for _, r := range plain {
		if r == 0 || unicode.Is(unicode.Mn, r) {
			continue
		}
		if isWideRune(r) {
			actualWidth += 2
		} else {
			actualWidth++
		}
	}
	padding := width - actualWidth
	if padding <= 0 {
		return s
	}

	switch align {
	case Left:
		return s + strings.Repeat(string(padChar), padding)
	case Right:
		return strings.Repeat(string(padChar), padding) + s
	case Center:
		left := padding / 2
		right := padding - left
		return strings.Repeat(string(padChar), left) + s + strings.Repeat(string(padChar), right)
	default:
		return s
	}
}
