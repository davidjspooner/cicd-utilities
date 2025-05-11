package textfmt

import (
	"io"
)

type RowType int

const (
	RowTypeColumns RowType = iota
	RowTypeBanner
)

type Row struct {
	Blocks  []*Block
	RowType RowType
}

func NewRow(rowType RowType, cells ...string) *Row {
	r := &Row{
		RowType: rowType,
		Blocks:  make([]*Block, len(cells)),
	}
	for i := range cells {
		r.Blocks[i] = NewBlock(cells[i])
	}
	return r
}

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
	//TODO
	// replce runs of whitespaces with a single space
	// remove leading and trailing spaces
	// split the lines on the two rune text sequence \n
	// make a row with a singe block
	// and add it to the table
	t.Rows = append(t.Rows, NewRow(RowTypeBanner, text))

	return nil
}

func (t *Table) AddRow(cells ...string) error {
	r := &Row{
		RowType: RowTypeColumns,
		Blocks:  make([]*Block, len(cells)),
	}
	for i := range cells {
		_ = i
		//TODO
		// replce runs of whitespaces with a single space
		// remove leading and trailing spaces
		// split the lines on the two rune text sequence \n
		// set as a column in the row
	}
	t.Rows = append(t.Rows, r)
	return nil
}

func (t *Table) RenderTo(w io.Writer) {
	// Copy WrapSpecs and adjust widths
}
