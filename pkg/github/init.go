package github

import "github.com/davidjspooner/cicd-utilities/pkg/command"

func InitCommands() ([]command.Command, error) {
	commands := []command.Command{}
	cmd1 := command.NewCommand(
		"github-release",
		"Create a GitHub release",
		executeGithubRelease,
		&GithubReleaseOptions{},
	)
	commands = append(commands, cmd1) // Initialize all commands
	cmd2 := command.NewCommand(
		"github-pr-update",
		"Update a GitHub pull request with the latest changes from the base branch",
		executeUpdateGithubPRMeta,
		&GithubPRUpdateOptions{},
	)
	commands = append(commands, cmd2)
	return commands, nil
}
