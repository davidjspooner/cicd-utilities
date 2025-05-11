package git

import "github.com/davidjspooner/cicd-utilities/pkg/command"

func Commands() []command.Command {

	gitCommand := command.NewCommand(
		"git",
		"Git commands",
		nil,
		&command.NoopOptions{},
	)

	cmd1 := command.NewCommand(
		"suggest-build-env",
		"Get the environment variables for the current build",
		executeGetGitEnv,
		&GetGitEnvOptions{},
	)
	cmd2 := command.NewCommand(
		"update-tag",
		"Automatically increment Git tags based on commit messages (e.g., fix:, feat:, breaking:)",
		executeBumpGitTag,
		&BumpGitTagOptions{
			Remote: "origin",
			Prefix: "v",
		},
	)

	gitCommand.SubCommands().MustAdd(cmd1, cmd2)
	return []command.Command{gitCommand}
}
