package github

import "github.com/davidjspooner/cicd-utilities/pkg/command"

func Commands() []command.Command {
	githubCommand := command.NewCommand(
		"github",
		"GitHub commands",
		nil,
		&command.NoopOptions{},
	)
	cmd1 := command.NewCommand(
		"release",
		"Create a GitHub release",
		executeGithubRelease,
		&GithubReleaseOptions{},
	)
	cmd2 := command.NewCommand(
		"update-pull-request",
		"Update a GitHub pull request with the latest changes from the base branch",
		executeUpdateGithubPRMeta,
		&GithubPRUpdateOptions{},
	)
	githubCommand.SubCommands().MustAdd(cmd1, cmd2)
	return []command.Command{githubCommand}
}
