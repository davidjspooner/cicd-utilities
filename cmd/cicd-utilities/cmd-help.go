package main

import (
	"context"
	"fmt"

	"github.com/davidjspooner/cicd-utilities/pkg/command"
)

func helpCommand(ctx context.Context, cmd command.Object, option *HelpOptions, args []string) error {
	err := command.CheckUnparsedOptions(args)
	if err != nil {
		return err
	}
	return fmt.Errorf("help command not implemented")
}

type HelpOptions struct {
}

func init() {
	cmd := command.New(
		"help",
		"Display help information",
		helpCommand,
		&HelpOptions{},
	)
	commands = append(commands, cmd)
}
