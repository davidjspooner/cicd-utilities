package command

import (
	"context"
	"fmt"
)

type shellOptions struct {
	Shell string `flag:"--shell,The shell to generate completion for. Supported shells: bash"`
}

type initOptions struct {
	Full bool `flag:"--full,Generate full completion script"`
}

func Completion() Command {
	branch := NewCommand("completion", "Support command line completion for shells", func(ctx context.Context, options *shellOptions, args []string) error {
		switch options.Shell {
		case "bash":
		case "":
			options.Shell = "bash"
		default:
			return fmt.Errorf("unsupported shell %s", options.Shell)
		}
		return nil
	}, &shellOptions{
		Shell: "bash",
	})
	init := NewCommand("init", "Initialize the completion for the shell", func(ctx context.Context, options *initOptions, args []string) error {
		if options.Full {
			print("#TODO FULL")
		} else {
			print("#TODO BASIC")
		}
		return nil
	}, &initOptions{})

	suggestions := NewCommand("suggest", "Get suggestions for the command", func(ctx context.Context, options *NoopOptions, args []string) error {
		print("TODO")
		print("THIS")
		print("GENERATION")
		return nil
	}, &NoopOptions{},
		LogicalGroup,
	)

	branch.SubCommands().MustAdd(init, suggestions)
	return branch
}
