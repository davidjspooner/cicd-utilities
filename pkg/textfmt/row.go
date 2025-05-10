package textfmt

type row struct {
	Blocks []*Block
}

func (r *row) MaxHeight() int {
	max := 0
	for _, b := range r.Blocks {
		if h := b.LineCount(); h > max {
			max = h
		}
	}
	return max
}

func (r *row) SetLineCount(n int) {
	for _, b := range r.Blocks {
		b.SetLineCount(n)
	}
}

func (r *row) ExpandTabs(tabStop int) {
	for _, b := range r.Blocks {
		b.ExpandTabs(tabStop)
	}
}

func (r *row) StripColors() {
	for _, b := range r.Blocks {
		b.StripColors()
	}
}

func (r *row) getColumnWidths() []int {
	w := make([]int, len(r.Blocks))
	for i, b := range r.Blocks {
		w[i] = b.Width()
	}
	return w
}
