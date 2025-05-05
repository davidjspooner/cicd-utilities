package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/davidjspooner/cicd-utilities/pkg/command"
	"github.com/davidjspooner/cicd-utilities/pkg/git"
)

type GetGitEnvOptions struct {
}

func init() {
	cmd := command.New(
		"get-build-env",
		"Get the environment variables for the current build",
		executeGetGitEnv,
		&GetGitEnvOptions{},
	)
	commands = append(commands, cmd)
}

func executeGetGitEnv(ctx context.Context, cmd command.Object, option *GetGitEnvOptions, args []string) error {
	err := command.CheckUnparsedOptions(args)
	if err != nil {
		return err
	}
	// Get the current branch
	currentBranch, err := git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}
	fmt.Printf("export BUILD_BRANCH=%s\n", currentBranch)
	fmt.Printf("export BUILD_NAME=%s\n", suggestBuildName())
	fmt.Printf("export BUILD_CONTEXT=%s\n", getBuildContext())
	now := time.Now().UTC()
	fmt.Printf("export BUILD_TIME=%s\n", now.Format(time.RFC1123))
	return nil
}

func suggestBuildName() string {
	// Check for uncommitted changes
	out, err := git.Run("status", "--porcelain")
	if err != nil {
		return "UNKNOWN"
	}
	if strings.TrimSpace(out) != "" {
		return "HEAD." + time.Now().Format("060102.1504")
	}

	// Check for a tag version
	tag, err := git.Run("tag", "--contains", "HEAD")
	if err == nil && tag != "" {
		return tag
	}

	// Fallback to short commit hash
	commitHash, err := git.Run("rev-parse", "--short", "HEAD")
	if err != nil {
		return "UNKNOWN"
	}
	return commitHash
}

func getBuildContext() string {
	// Check for github actions
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return os.Getenv("GITHUB_RUN_ID")
	}

	// Fallback to user@hostname
	user := os.Getenv("USER")
	if user == "" {
		user = "unknown"
	}
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return user + "@" + hostname
}
