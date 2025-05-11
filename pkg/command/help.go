package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/davidjspooner/cicd-utilities/pkg/textfmt"
)

// renderHelpText displays help information for all commands in the execution plan.
func renderHelpText(plan *plan) error {

	columnSpecs := []*textfmt.WrapSpec{
		{
			ExactWidth: 2,
			Align:      textfmt.Left,
			PadChar:    ' ',
		},
		{
			MaxWidth: 40,
			MinWidth: 16,
			Align:    textfmt.Left,
			PadChar:  ' ',
		},
		{
			ExactWidth: 3,
			Align:      textfmt.Center,
			PadChar:    ' ',
		},
		{
			MaxWidth: 60,
			Align:    textfmt.Left,
			PadChar:  ' ',
		},
	}

	lastCommand := plan.steps[len(plan.steps)-1]
	lastSubCommands := lastCommand.command().SubCommands()
	lastSubCommands.SortAlphabetically()

	os.Stdout.WriteString("Usage:\n")
	commandStr := strings.Builder{}
	commandStr.WriteString("  ")
	for _, step := range plan.steps {
		commandStr.WriteString(fmt.Sprintf("%s ", step.command().Aliases()[0]))
	}
	if lastSubCommands.Count() > 0 {
		commandStr.WriteString("[subcommand] ")
	}

	commandStr.WriteString("<flags...> <args...>")

	os.Stdout.WriteString(commandStr.String())
	os.Stdout.WriteString("\n\n")

	table := textfmt.NewTable(columnSpecs...)
	table.AddBanner("Command Hierarchy:")
	for _, step := range plan.steps {
		table.AddRow("", step.command().Aliases()[0], "-", step.command().Help())
	}

	if lastSubCommands.Count() > 0 {
		table.AddBanner("")
		table.AddBanner(fmt.Sprintf("Available Subcommands for %s:", lastCommand.command().Aliases()[0]))
		for _, subCommand := range lastSubCommands.commands {
			table.AddRow("", subCommand.Aliases()[0], "-", subCommand.Help())
		}
	}
	for i, step := range plan.steps {
		cmd := step.command()
		flags, err := cmd.Flags()
		if flags == nil {
			continue
		}
		table.AddBanner("")
		if i == 0 {
			table.AddBanner("Global Flags:")
		} else {
			table.AddBanner(fmt.Sprintf("Flags for %s:", cmd.Aliases()[0]))
		}
		if err != nil {
			return fmt.Errorf("error getting flags for command %s: %v", cmd.Aliases()[0], err)
		}
		for _, flag := range flags {
			name := flag.Aliases()[0]
			metaVar := flag.MetaVar()
			if metaVar != "" {
				name = fmt.Sprintf("%s %s", name, metaVar)
			}
			defaultValue := flag.DefaultValue()
			help := flag.Help()
			if defaultValue != "" {
				help = fmt.Sprintf("%s (default: %s)", help, defaultValue)
			}

			table.AddRow("", name, "-", help)
		}
	}
	err := table.RenderTo(os.Stdout)
	os.Stdout.WriteString("\n\n")
	return err
}
