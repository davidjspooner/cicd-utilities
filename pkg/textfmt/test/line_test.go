package textfmt_test

import (
    "testing"

    "github.com/davidjspooner/cicd-utilities/pkg/textfmt"
)

func TestLinePaddingAndTabs(t *testing.T) {
    l := &textfmt.Line{Text: "x\t中"}
    l.ExpandTabs(4)
    l.Pad(10, textfmt.Center, ' ')
    if l.Width() == 0 {
        t.Errorf("Expected non-zero width after padding and tab expansion")
    }
}