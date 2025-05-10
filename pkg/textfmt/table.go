package textfmt

import (
	"fmt"
	"io"
)

type ColumnDefinition struct {
	MaxWidth   int
	Alignment  Align
	AllowColor bool
}

type Table struct {
	Rows              []*row
	ColumnDefinitions []ColumnDefinition
	ColumnSeparator   string
}

func NewTable(defs ...ColumnDefinition) *Table {
	return &Table{
		ColumnDefinitions: defs,
		ColumnSeparator:   " ",
	}
}

func (t *Table) NewRow(cells ...string) error {
	if len(cells) != len(t.ColumnDefinitions) {
		return fmt.Errorf("expected %d columns, got %d", len(t.ColumnDefinitions), len(cells))
	}
	r := &row{}
	for i, text := range cells {
		block := NewBlock(text)
		if !t.ColumnDefinitions[i].AllowColor {
			block.StripColors()
		}
		r.Blocks = append(r.Blocks, &block)
	}
	t.Rows = append(t.Rows, r)
	return nil
}

func (t *Table) getColumnWidths() []int {
	widths := make([]int, len(t.ColumnDefinitions))
	for _, r := range t.Rows {
		w := r.getColumnWidths()
		for i := range w {
			widths[i] = max(widths[i], w[i])
		}
	}
	for i, def := range t.ColumnDefinitions {
		if def.MaxWidth > 0 {
			widths[i] = min(widths[i], def.MaxWidth)
		}
	}
	return widths
}

func (t *Table) WrapAndPad() {
	widths := t.getColumnWidths()
	for _, r := range t.Rows {
		for i, b := range r.Blocks {
			b.WordWrap(widths[i])
			for _, l := range b.Lines {
				l.Pad(widths[i], t.ColumnDefinitions[i].Alignment, ' ')
			}
		}
	}
	for _, r := range t.Rows {
		h := r.MaxHeight()
		r.SetLineCount(h)
	}
}

func (t *Table) Render(w io.Writer) error {
	for _, r := range t.Rows {
		h := r.MaxHeight()
		for i := 0; i < h; i++ {
			for j, b := range r.Blocks {
				if _, err := io.WriteString(w, b.Lines[i].Text); err != nil {
					return err
				}
				if j < len(r.Blocks)-1 {
					if _, err := io.WriteString(w, t.ColumnSeparator); err != nil {
						return err
					}
				}
			}
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}
		}
	}
	return nil
}
