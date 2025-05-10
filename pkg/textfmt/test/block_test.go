package textfmt_test

import (
	"testing"

	"github.com/davidjspooner/cicd-utilities/pkg/textfmt"
)

func TestBlockBehavior(t *testing.T) {
	b := textfmt.NewBlock("one\ntwo\n")
	b.WordWrap(3)
	b.SetLineCount(5)

	// Avoid trimming if you're testing SetLineCount behavior
	if b.LineCount() != 5 {
		t.Errorf("Expected 5 lines after SetLineCount, got %d", b.LineCount())
	}

	b.ExpandTabs(4)
	b.Pad(4, textfmt.Right, '_')
	b.StripColors()

	if b.Width() == 0 {
		t.Errorf("Expected block to have non-zero width after Pad")
	}
}
