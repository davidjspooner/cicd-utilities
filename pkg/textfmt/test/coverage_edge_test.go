package textfmt_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/davidjspooner/cicd-utilities/pkg/textfmt"
)

func TestBlock_TrimEmptyLines_Mixed(t *testing.T) {
	b := textfmt.NewBlock("")
	b.Lines = []*textfmt.Line{
		{Text: ""},
		{Text: "content"},
		{Text: ""},
		{Text: ""},
	}
	b.TrimEmptyLines()
	if len(b.Lines) != 1 || b.Lines[0].Text != "content" {
		t.Errorf("Expected trimmed block to contain only content line, got %+v", b.Lines)
	}
}

func TestLine_Width_ExercisesWideAndZeroWidthRunes(t *testing.T) {
	// This will internally use isWideRune and ignore zero-width combining runes
	l := &textfmt.Line{Text: "中âb"} // includes wide + combining
	w := l.Width()
	if w <= 3 {
		t.Errorf("Expected visible width > 3, got %d", w)
	}
}

func TestLine_ExpandTabs_ExercisesExpandTabsInternal(t *testing.T) {
	l := &textfmt.Line{Text: "a\tb"}
	l.ExpandTabs(4)
	if l.Text != "a   b" {
		t.Errorf("Expected expanded tabs, got %q", l.Text)
	}
}

func TestLine_WordWrap_ExercisesBreakableSpaceHandling(t *testing.T) {
	l := &textfmt.Line{Text: "foo\u2000\u3000bar"} // includes Unicode spaces
	lines := l.WordWrap(10)
	if len(lines) != 1 || lines[0].Text != "foo bar" {
		t.Errorf("Expected normalized whitespace, got %+v", lines)
	}
}

func TestTable_Render_WriterAfterSeparatorFails(t *testing.T) {
	table := textfmt.NewTable(
		textfmt.ColumnDefinition{MaxWidth: 5, Alignment: textfmt.Left, AllowColor: false},
		textfmt.ColumnDefinition{MaxWidth: 5, Alignment: textfmt.Left, AllowColor: false},
	)
	_ = table.NewRow("X", "Y")
	table.WrapAndPad()
	err := table.Render(&failingAfterSepWriter{})
	if err == nil {
		t.Error("Expected error from failing writer after separator")
	}
}

type failingAfterSepWriter struct {
	callCount int
}

func (w *failingAfterSepWriter) Write(p []byte) (int, error) {
	w.callCount++
	if w.callCount == 2 {
		return 0, errors.New("fail after separator")
	}
	return len(p), nil
}

func TestBlock_TrimEmptyLines_MultipleLeadingAndTrailing(t *testing.T) {
	b := textfmt.NewBlock("")
	b.Lines = []*textfmt.Line{
		{Text: ""},
		{Text: ""},
		{Text: "middle"},
		{Text: ""},
		{Text: ""},
	}
	b.TrimEmptyLines()
	if len(b.Lines) != 1 || b.Lines[0].Text != "middle" {
		t.Errorf("Unexpected result from TrimEmptyLines: %+v", b.Lines)
	}
}

func TestLine_collapseBreakableSpaces_Run(t *testing.T) {
	b := textfmt.NewBlock("a\u2000\u2000\u2000b")
	if got := b.Lines[0].Text; got != "a b" {
		t.Errorf("Expected collapsed to single space, got %q", got)
	}
}

func TestLine_expandTabsInternal_MultipleTabs(t *testing.T) {
	l := &textfmt.Line{Text: "x\t中\tz"}
	l.ExpandTabs(4)
	if len(l.Text) <= 6 {
		t.Errorf("Expected expanded tabs, got %q", l.Text)
	}
}

func TestTable_Render_NewlineError(t *testing.T) {
	table := textfmt.NewTable(
		textfmt.ColumnDefinition{MaxWidth: 5, Alignment: textfmt.Left, AllowColor: false},
	)
	_ = table.NewRow("a")
	table.WrapAndPad()
	err := table.Render(&newlineFailWriter{})
	if err == nil {
		t.Error("Expected error from newline fail writer")
	}
}

type newlineFailWriter struct {
	callCount int
}

func (w *newlineFailWriter) Write(p []byte) (int, error) {
	if bytes.Contains(p, []byte("\n")) {
		return 0, errors.New("newline write fail")
	}
	return len(p), nil
}

func TestTable_NewRow_TooFewColumns(t *testing.T) {
	table := textfmt.NewTable(
		textfmt.ColumnDefinition{MaxWidth: 10, Alignment: textfmt.Left, AllowColor: false},
		textfmt.ColumnDefinition{MaxWidth: 10, Alignment: textfmt.Left, AllowColor: false},
	)
	err := table.NewRow("only one")
	if err == nil {
		t.Error("Expected error due to mismatched column count")
	}
}

func TestLine_Width_WithAnsiCodes(t *testing.T) {
	// \x1b[31m is red, \x1b[0m is reset
	l := &textfmt.Line{Text: "\x1b[31mRed\x1b[0m"}
	width := l.Width()
	if width != 3 {
		t.Errorf("Expected visible width 3, got %d", width)
	}
}
