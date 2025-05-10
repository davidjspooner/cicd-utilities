package command

import (
	"testing"
)

func TestGetDefinedArgs(t *testing.T) {
	type TestStruct struct {
		Field1 string  `flag:"--field1|-1,Help for field1"`
		Field2 int     `flag:"--field2|-2,Help for field2"`
		Field3 bool    `flag:"--field3|-3,Help for field3"`
		Field4 float32 `flag:"--field4|-4,Help for field4"`
	}

	defaults := &TestStruct{}
	args, err := getFlagDefinitions(defaults)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(args) != 4 {
		t.Fatalf("expected 4 args, got %d", len(args))
	}
}

func TestGetDefinedArgsNegative(t *testing.T) {
	type InvalidStruct struct {
		Field1 string
	}
	_, err := getFlagDefinitions(&InvalidStruct{})
	if err != nil {
		t.Errorf("expected zero options( which is ok but got), got error: %v", err)
	}
}

func TestMalformedTags(t *testing.T) {
	type Malformed struct {
		Bad1 string `flag:"--onlyflag"`                 // missing help
		Bad2 string `flag:"--name|,Missing short flag"` // bad name
		Bad3 string `flag:"--name|-n"`                  // no help
	}
	_, err := getFlagDefinitions(&Malformed{})
	if err == nil {
		t.Error("expected error for malformed tags, got nil")
	}
}

func TestDuplicateNames(t *testing.T) {
	type Duplicate struct {
		Field1 string `flag:"--flag|-f,Flag one"`
		Field2 string `flag:"--flag|-g,Flag two"`
	}
	_, err := getFlagDefinitions(&Duplicate{})
	if err == nil {
		t.Error("expected error for duplicate flags, got nil")
	}
}

func TestEnvVarTag(t *testing.T) {
	type EnvTest struct {
		FromEnv string `flag:"$MY_ENV_VAR,Read from environment"`
	}

	args, err := getFlagDefinitions(&EnvTest{})
	if err != nil {
		t.Errorf("unexpected error for env var tag: %v", err)
	}

	if len(args) != 1 || args[0].aliases[0] != "$MY_ENV_VAR" {
		t.Errorf("unexpected parsing of env var tag: %+v", args)
	}
}
