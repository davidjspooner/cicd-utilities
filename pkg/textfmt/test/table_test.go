package textfmt_test

import (
    "bytes"
    "errors"
    "testing"

    "github.com/davidjspooner/cicd-utilities/pkg/textfmt"
)

func TestTable_RenderAll(t *testing.T) {
    table := textfmt.NewTable(
        textfmt.ColumnDefinition{MaxWidth: 8, Alignment: textfmt.Left, AllowColor: false},
        textfmt.ColumnDefinition{MaxWidth: 5, Alignment: textfmt.Right, AllowColor: false},
    )
    _ = table.NewRow("hello", "9")
    _ = table.NewRow("world", "42")
    table.WrapAndPad()

    var buf bytes.Buffer
    if err := table.Render(&buf); err != nil {
        t.Fatalf("Render failed: %v", err)
    }

    if got := buf.String(); got == "" {
        t.Errorf("Expected rendered output, got empty string")
    }
}

type failingWriter struct{}

func (f *failingWriter) Write(p []byte) (int, error) {
    return 0, errors.New("fail")
}

func TestTable_RenderWriterFailure(t *testing.T) {
    table := textfmt.NewTable(textfmt.ColumnDefinition{MaxWidth: 5, Alignment: textfmt.Left, AllowColor: false})
    _ = table.NewRow("fail")
    table.WrapAndPad()

    err := table.Render(&failingWriter{})
    if err == nil {
        t.Error("Expected error from failing writer")
    }
}