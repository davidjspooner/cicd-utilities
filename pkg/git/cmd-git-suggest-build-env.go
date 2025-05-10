package git

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

type GetGitEnvOptions struct {
}

func executeGetGitEnv(ctx context.Context, options *GetGitEnvOptions, args []string) error {
	// Get the current branch
	currentBranch, err := GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}
	fmt.Printf("BUILD_BRANCH=%s\n", currentBranch)
	fmt.Printf("BUILD_VERSION=%s\n", suggestBuildName())
	fmt.Printf("BUILD_CONTEXT=%s\n", getBuildContext())
	now := time.Now().UTC()
	fmt.Printf("BUILD_TIME=%s\n", now.Format(time.RFC1123))
	return nil
}

func suggestBuildName() string {
	// Check for uncommitted changes
	out, err := Run("status", "--porcelain")
	if err != nil {
		return "UNKNOWN"
	}
	if strings.TrimSpace(out) != "" {
		return "HEAD." + time.Now().Format("060102.1504")
	}

	// Check for a tag version
	tag, err := Run("tag", "--contains", "HEAD")
	if err == nil && tag != "" {
		return tag
	}

	// Fallback to short commit hash
	commitHash, err := Run("rev-parse", "--short", "HEAD")
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
