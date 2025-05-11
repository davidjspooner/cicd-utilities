package textfmt_test

import (
	"testing"

	"github.com/davidjspooner/cicd-utilities/pkg/textfmt"
)

func TestLineWidth(t *testing.T) {
	line := &textfmt.Line{Text: "Hello\\tWorld"}
	width := line.Width(4)
	if width != 13 {
		t.Errorf("expected width 13, got %d", width)
	}
}
