package textfmt

import (
	"fmt"
	"io"
)

type Table struct {
	Rows            []*Row
	WrapSpec        []*WrapSpec
	BannerSpec      *WrapSpec
	ColumnSeparator string
}

func NewTable(defs ...*WrapSpec) *Table {
	t := &Table{
		WrapSpec:        defs,
		ColumnSeparator: " ",
		BannerSpec:      &WrapSpec{},
	}
	for i := range t.WrapSpec {
		t.WrapSpec[i].normalizeSpec()
	}
	t.BannerSpec.normalizeSpec()
	return t
}

func (t *Table) AddBanner(text string) error {
	row := &Row{
		RowType: RowTypeBanner,
		Columns: []*Block{NewBlock(text)},
	}
	t.Rows = append(t.Rows, row)
	return nil
}

func (t *Table) AddRow(cells ...string) error {
	r := &Row{
		RowType: RowTypeColumns,
		Columns: make([]*Block, len(cells)),
	}
	for i, cell := range cells {
		r.Columns[i] = NewBlock(cell)
	}
	t.Rows = append(t.Rows, r)
	return nil
}

func (t *Table) RenderTo(w io.Writer) error {
	// Helper function to calculate the maximum of two integers
	max := func(a, b int) int {
		if a > b {
			return a
		}
		return b
	}

	// Create a copy of WrapSpecs for columns
	tempWrapSpecs := make([]*WrapSpec, len(t.WrapSpec))
	for i, spec := range t.WrapSpec {
		temp := *spec // Copy the WrapSpec
		tempWrapSpecs[i] = &temp
	}

	// Update the width of the temporary WrapSpecs for columns
	for _, row := range t.Rows {
		if row.RowType == RowTypeColumns {
			for i, block := range row.Columns {
				tempWrapSpecs[i].Width = max(tempWrapSpecs[i].Width, block.Width(0))
			}
		}
	}

	for i, spec := range tempWrapSpecs {
		spec.Width = min(spec.Width, t.WrapSpec[i].Width)
	}

	totalWidth := 0
	for _, spec := range tempWrapSpecs {
		totalWidth += spec.Width
	}
	// Add space for column separators
	totalWidth += (len(tempWrapSpecs) - 1) * len(t.ColumnSeparator)

	// Create a copy of the WrapSpec for banners
	tempBannerSpec := *t.BannerSpec // Copy the WrapSpec
	tempBannerSpec.Width = totalWidth

	// Render each row
	for _, row := range t.Rows {
		var wrapSpecs []*WrapSpec
		if row.RowType == RowTypeColumns {
			wrapSpecs = tempWrapSpecs
		} else if row.RowType == RowTypeBanner {
			wrapSpecs = []*WrapSpec{&tempBannerSpec}
		}
		if err := row.RenderTo(w, wrapSpecs); err != nil {
			return fmt.Errorf("error rendering row: %v", err)
		}
	}

	// Ensure the function always returns a value
	return nil
}
