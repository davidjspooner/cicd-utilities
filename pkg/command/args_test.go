package command

import (
	"reflect"
	"testing"
)

func TestGetDefinedArgs(t *testing.T) {
	type TestStruct struct {
		Field1 string  `arg:"--field1|-1,Help for field1"`
		Field2 int     `arg:"--field2|-2,Help for field2"`
		Field3 bool    `arg:"--field3|-3,Help for field3"`
		Field4 float32 `arg:"--field4|-4,Help for field4"`
	}

	defaults := &TestStruct{}
	args, err := getDefinedOptions(defaults)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 4 {
		t.Fatalf("expected 4 args, got %d", len(args))
	}

	if args[0].Names[0] != "--field1" || args[0].Names[1] != "-1" {
		t.Errorf("unexpected names for Field1: %v", args[0].Names)
	}
	if args[0].field.Type.Kind() != reflect.String {
		t.Errorf("unexpected type for Field1: %s", args[0].field.Type.Kind())
	}
	if args[1].Names[0] != "--field2" || args[1].Names[1] != "-2" {
		t.Errorf("unexpected names for Field2: %v", args[1].Names)
	}
	if args[1].field.Type.Kind() != reflect.Int {
		t.Errorf("unexpected type for Field2: %s", args[1].field.Type.Kind())
	}
	if args[2].Names[0] != "--field3" || args[2].Names[1] != "-3" {
		t.Errorf("unexpected names for Field3: %v", args[2].Names)
	}
	if args[2].field.Type.Kind() != reflect.Bool {
		t.Errorf("unexpected type for Field3: %s", args[2].field.Type.Kind())
	}
	if args[3].Names[0] != "--field4" || args[3].Names[1] != "-4" {
		t.Errorf("unexpected names for Field4: %v", args[3].Names)
	}
	if args[3].field.Type.Kind() != reflect.Float32 {
		t.Errorf("unexpected type for Field4: %s", args[3].field.Type.Kind())
	}
}

func TestGetDefinedArgsNegative(t *testing.T) {
	// Test with invalid struct
	type InvalidStruct struct {
		Field1 string
	}
	_, err := getDefinedOptions(&InvalidStruct{})
	if err == nil {
		t.Errorf("expected error for struct without arg tags, got nil")
	}
}

func TestExpandInputArgs(t *testing.T) {
	input := []string{"--key=value", "-abc=harry", "positional"}
	expected := []string{"--key", "value", "-a", "-b", "-c", "harry", "positional"}

	expanded, err := expandInputArgs(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(expanded, expected) {
		t.Errorf("expected %v, got %v", expected, expanded)
	}
}
