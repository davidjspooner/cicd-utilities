package command

import "context"

// NOTE this interface cannot be implemented outside this package
type Command interface {
	Name() string
	Aliases() []string
	Help() string
	SubCommands() *CommandGroup
	With(modifiers ...Middeware) Command
	Flags() ([]Flag, error)
}

type ExecuteFunc[T any] func(ctx context.Context, options *T, args []string) error

func (fn ExecuteFunc[T]) IsDefined() bool {
	return fn != nil
}
