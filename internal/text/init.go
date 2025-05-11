package text

import (
	"github.com/davidjspooner/cicd-utilities/pkg/command"
)

func Commands() []command.Command {
	textCommand := command.NewCommand(
		"text",
		"text commands ( including templates )",
		nil,
		&command.NoopOptions{},
	)
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
		"support",
		"List supported template types",
		listSupportedTemplateTypes,
		&command.NoopOptions{},
	)
	templateCommand.SubCommands().MustAdd(cmd1)
	textCommand.SubCommands().MustAdd(templateCommand)
	return []command.Command{textCommand, cmd2}
}
