package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/davidjspooner/cicd-utilities/pkg/archive"
	"github.com/davidjspooner/cicd-utilities/pkg/command"
	"github.com/davidjspooner/cicd-utilities/pkg/git"
)

type GlobalOptions struct {
	command.LogOptions
}

func main() {
	command.RootCommand = command.NewCommand("", "A utility for CI/CD operations",
		func(ctx context.Context, options *GlobalOptions, args []string) error {
			level, err := options.LogOptions.Parse()
			slog.SetLogLoggerLevel(level)
			if err != nil {
				return err
			}
			return nil
		}, &GlobalOptions{LogOptions: command.LogOptions{Level: "info"}},
		command.LogicalGroup)

	versionCommand := command.VersionCommand()
	gitCommands := git.Commands()
	archiveCommands := archive.Commands()
	command.RootCommand.SubCommands().MustAdd(versionCommand, gitCommands, archiveCommands)

	err := command.Run(context.Background(), os.Args[1:])
	if err != nil {
		slog.Error("Failed Execution", "msg", err)
		os.Exit(1)
	}
}
