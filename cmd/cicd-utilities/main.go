package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/davidjspooner/cicd-utilities/internal/archive"
	"github.com/davidjspooner/cicd-utilities/internal/git"
	"github.com/davidjspooner/cicd-utilities/internal/github"
	"github.com/davidjspooner/cicd-utilities/internal/man"
	"github.com/davidjspooner/cicd-utilities/internal/template"
	"github.com/davidjspooner/cicd-utilities/pkg/command"
)

type GlobalOptions struct {
	command.LogOptions
}

func main() {
	root := command.NewCommand("", "A utility for CI/CD operations",
		func(ctx context.Context, options *GlobalOptions, args []string) error {
			level, err := options.LogOptions.Parse()
			opts := slog.HandlerOptions{
				Level: level,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.TimeKey {
						return slog.Attr{} // remove timestamp
					}
					return a
				},
			}

			// Create a TextHandler with those options
			handler := slog.NewTextHandler(os.Stdout, &opts)
			logger := slog.New(handler)

			// Set this logger as the default
			slog.SetDefault(logger)
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
	githubCommands := github.Commands()
	templateCommands := template.Commands()
	manCommands := man.Commands()

	subcommands := command.RootCommand.SubCommands()
	subcommands.MustAdd(
		versionCommand,
		gitCommands,
		archiveCommands,
		githubCommands,
		manCommands,
		templateCommands,
	)

	err := command.Run(context.Background(), os.Args[1:])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
