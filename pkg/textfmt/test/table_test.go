package textfmt_test

import (
	"bytes"
	"testing"

	"github.com/davidjspooner/cicd-utilities/pkg/textfmt"
)

func TestTableWithColorCodes(t *testing.T) {
	wrapSpec := textfmt.NewWrapSpec(20, 4, textfmt.Left, textfmt.AllowColor, ' ')
	table := textfmt.NewTable(wrapSpec)
	table.Rows = append(table.Rows, textfmt.NewRow(textfmt.RowTypeColumns, "\x1b[31mRed\x1b[0m", "\x1b[32mGreen\x1b[0m"))
	table.Rows = append(table.Rows, textfmt.NewRow(textfmt.RowTypeColumns, "\x1b[34mBlue\x1b[0m", "\x1b[33mYellow\x1b[0m"))

	// Render the table to a buffer
	var buffer bytes.Buffer
	err := table.Rows[0].RenderTo(&buffer, []*textfmt.WrapSpec{wrapSpec, wrapSpec})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTableWithUnicodeChars(t *testing.T) {
	wrapSpec := textfmt.NewWrapSpec(20, 4, textfmt.Left, textfmt.NoColor, ' ')
	table := textfmt.NewTable(wrapSpec)
	table.Rows = append(table.Rows, textfmt.NewRow(textfmt.RowTypeColumns, "こんにちは", "世界"))
	table.Rows = append(table.Rows, textfmt.NewRow(textfmt.RowTypeColumns, "😀", "🌍"))

	// Render the table to a buffer
	var buffer bytes.Buffer
	err := table.Rows[0].RenderTo(&buffer, []*textfmt.WrapSpec{wrapSpec, wrapSpec})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
