package git

import "github.com/davidjspooner/cicd-utilities/pkg/command"

func Commands() []command.Command {
	commands := []command.Command{}
	cmd1 := command.NewCommand(
		"git-suggest-build-env",
		"Get the environment variables for the current build",
		executeGetGitEnv,
		&GetGitEnvOptions{},
	)
	commands = append(commands, cmd1) // Initialize all commands
	cmd2 := command.NewCommand(
		"git-update-tag",
		"Automatically increment Git tags based on commit messages (e.g., fix:, feat:, breaking:)",
		executeBumpGitTag,
		&BumpGitTagOptions{
			Remote: "origin",
			Prefix: "v",
		},
	)
	commands = append(commands, cmd2) // Initialize all commands
	return commands
}
