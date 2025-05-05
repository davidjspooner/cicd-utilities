package command

import (
	"context"
	"errors"
	"testing"
)

type MockOptions struct {
	Bool   bool    `arg:"--bool|-b,Test bool"`
	Value  string  `arg:"--value|-v,Test value"`
	Float  float64 `arg:"--float|-F,Test float"`
	Number int     `arg:"--number|-n,Test number"`
}

func TestCommand_Execute(t *testing.T) {
	executeFunc := func(ctx context.Context, cmd Object, option *MockOptions, args []string) error {
		if option.Bool && option.Value == "test" && option.Number == 42 {
			return nil
		}
		return errors.New("execution failed")
	}

	cmd := New("test", "Test command", executeFunc, &MockOptions{})

	args := []string{"--bool", "--value", "test", "--float", "42.2", "--number", "42"}
	ctx := context.Background()
	if err := cmd.Execute(ctx, args); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCommand_ExecuteSubCommand(t *testing.T) {
	executeFunc := func(ctx context.Context, cmd Object, option *MockOptions, args []string) error {
		if option.Bool && option.Value == "test" && option.Number == 42 {
			return nil
		}
		return errors.New("execution failed")
	}

	subCmd := New("child", "Child command", executeFunc, &MockOptions{
		Value:  "test",
		Number: 42,
	})

	group := Group{
		subCmd,
	}

	args := []string{"--bool"}
	ctx := context.Background()

	if err := group.Execute(ctx, "child", args); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := group.Execute(ctx, "nonexistent", args); err == nil {
		t.Errorf("expected error for nonexistent subcommand, got nil")
	}
}
