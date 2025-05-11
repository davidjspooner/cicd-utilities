package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/davidjspooner/cicd-utilities/pkg/archive"
	"github.com/davidjspooner/cicd-utilities/pkg/command"
	"github.com/davidjspooner/cicd-utilities/pkg/git"
	"github.com/davidjspooner/cicd-utilities/pkg/github"
)

type GlobalOptions struct {
	command.LogOptions
}

func main() {
	root := command.NewCommand("", "A utility for CI/CD operations",
		func(ctx context.Context, options *GlobalOptions, args []string) error {
			level, err := options.LogOptions.Parse()
			slog.SetLogLoggerLevel(level)
			if err != nil {
				return err
			}
			return nil
		}, &GlobalOptions{LogOptions: command.LogOptions{Level: "info"}},
		command.LogicalGroup)

	command.RootCommand = root
	versionCommand := command.VersionCommand()
	gitCommands := git.Commands()
	archiveCommands := archive.Commands()
	subcommands := command.RootCommand.SubCommands()
	githubCommands := github.Commands()
	subcommands.MustAdd(versionCommand, gitCommands, archiveCommands, githubCommands)

	err := command.Run(context.Background(), os.Args[1:])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
