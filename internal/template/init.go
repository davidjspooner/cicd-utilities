package template

import (
	"github.com/davidjspooner/cicd-utilities/pkg/command"
)

func Commands() []command.Command {
	templateCommand := command.NewCommand(
		"template",
		"Template commands",
		nil,
		&command.NoopOptions{},
	)
	cmd1 := command.NewCommand(
		"expand",
		"Expand a template file",
		expandTemplate,
		&expandOptions{},
	)
	cmd2 := command.NewCommand(
		"man",
		"Manual for template file syntax",
		templateManPage,
		&command.NoopOptions{},
	)
	templateCommand.SubCommands().MustAdd(cmd1, cmd2)
	return []command.Command{templateCommand}
}
