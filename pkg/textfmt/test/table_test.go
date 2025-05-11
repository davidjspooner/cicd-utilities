package textfmt_test

import (
	"strings"
	"testing"

	"github.com/davidjspooner/cicd-utilities/pkg/textfmt"
)

type testCase struct {
	name            string
	wrapspec        []*textfmt.WrapSpec
	columnSeperator string
	inputCells      []string
	expectedLines   []string
	expectedError   string
}

func (tc *testCase) runTest(t *testing.T) {
	//make a row of cells
	row := textfmt.NewRow(textfmt.RowTypeColumns, tc.inputCells...)
	b := strings.Builder{}
	err := row.RenderTo(&b, tc.wrapspec)
	if err != nil {
		if tc.expectedError == "" {
			t.Errorf("%s: unexpected error:\n  %v", tc.name, err)
		} else if err.Error() != tc.expectedError {
			t.Errorf("%s: expected error\n   %q\ngot\n   %q", tc.name, tc.expectedError, err.Error())
		}
		return
	} else if tc.expectedError != "" {
		t.Errorf("%s: expected error\n  %q\n  got nil", tc.name, tc.expectedError)
		return
	}
	gotLines := strings.Split(b.String(), "\n")
	overlap := min(len(gotLines), len(tc.expectedLines))
	for i := 0; i < overlap; i++ {
		if gotLines[i] != tc.expectedLines[i] {
			t.Errorf("%s: line #%d\n  expected %q\n  got      %q", tc.name, i, tc.expectedLines[i], gotLines[i])
		}
	}
	if overlap < len(tc.expectedLines) {
		//show extra lines in expected output
		for i := overlap; i < len(tc.expectedLines); i++ {
			t.Errorf("%s: expected extra line #%d:\n  %q", tc.name, i, tc.expectedLines[i])
		}
	}
	if overlap < len(gotLines) {
		//show extra lines in actual output
		for i := overlap; i < len(gotLines); i++ {
			t.Errorf("%s: got extra line #%d\n   %q", tc.name, i, gotLines[i])
		}
	}
}

func TestTableRendering(t *testing.T) {
	tc := &testCase{
		name: "Basic table rendering",
		wrapspec: []*textfmt.WrapSpec{
			{Width: 10, Align: textfmt.Left, PadChar: ' '},
			{Width: 10, Align: textfmt.Left, PadChar: ' '},
			{Width: 10, Align: textfmt.Left, PadChar: ' '},
		},
		columnSeperator: " | ",
		inputCells:      []string{"Column1", "Column2", "Column3"},
		expectedLines: []string{
			"Column1   | Column2   | Column3   ",
		},
		expectedError: "",
	}
	tc.runTest(t)
}
