package git

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/davidjspooner/cicd-utilities/pkg/semantic"
)

type BumpGitTagOptions struct {
	Prefix string `flag:"--prefix,Prefix string"`
	Suffix string `flag:"--suffix,Suffix string"`
	DryRun bool   `flag:"--dry-run,Do not push the tag"`
	Remote string `flag:"--remote,Remote to push the tag to"`
}

func executeBumpGitTag(ctx context.Context, option *BumpGitTagOptions, args []string) error {

	// Get the current branch
	currentBranch, err := GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}

	// Get the latest tag
	latestTag, err := getLatestTag(ctx, currentBranch)
	if err != nil {
		return fmt.Errorf("failed to get the latest tag: %v", err)
	}

	_, _, currentVersion, err := semantic.ExtractVersionFromTag(latestTag)
	if err != nil {
		return fmt.Errorf("failed to extract version from tag: %v", err)
	}

	slog.Debug("Latest tag found", "tag", latestTag)
	slog.Info("Current version", "version", currentVersion.String())

	// Get commit messages since the latest tag
	commitMessages, err := Run("log", fmt.Sprintf("%s..HEAD", latestTag), "--pretty=format:%s")
	if err != nil {
		return fmt.Errorf("failed to get commit messages: %v", err)
	}
	commits := splitLines(commitMessages)
	nCommits := 0
	for i := range commits {
		commits[i] = strings.TrimSpace(commits[i])
		if commits[i] == "" {
			continue
		}
		nCommits++
	}
	if nCommits == 0 {
		fmt.Println("No changes deteced, no version increment needed.")
		return nil
	}

	// Determine the version increment
	increment, err := semantic.Bumps.GetVersionBump(commits)
	if err != nil {
		return fmt.Errorf("failed to determine version increment: %v", err)
	}

	// Increment the version
	newVersion, err := currentVersion.Increment(increment)
	if err != nil {
		return fmt.Errorf("failed to increment version: %v", err)
	}

	// Construct the new tag
	newTag := fmt.Sprintf("%s%s%s", option.Prefix, newVersion.String(), option.Suffix)

	fmt.Printf("Increment: %s\n", increment)

	if option.DryRun {
		fmt.Println("Dry run enabled.")
		fmt.Printf("Would create new tag: %s\n", newTag)
		return nil
	}

	// Create and push the new tag
	if _, err := Run("tag", newTag); err != nil {
		return fmt.Errorf("failed to create tag: %v", err)
	}
	if _, err := Run("push", option.Remote, newTag); err != nil {
		return fmt.Errorf("failed to push tag: %v", err)
	}

	fmt.Printf("Successfully created and pushed tag: %s\n", newTag)
	return nil
}

func getLatestTag(ctx context.Context, branch string) (string, error) {

	commits, err := Run("rev-list", "--tags", "--no-walk", "--abbrev=0", "--date-order", branch)
	if err != nil {
		return "", fmt.Errorf("failed to get latest tags: %v", err)
	}
	var bestVersion semantic.Version
	var bestTag string

	slog.Debug("Latest commits found", "commits", commits)
	commitList := splitLines(commits)
	for _, commit := range commitList {
		if commit == "" {
			continue
		}
		commit = strings.TrimSpace(commit)
		checkCmd, err := Run("merge-base", "--is-ancestor", commit, branch)
		if err != nil {
			continue
		}
		slog.Debug("Checked commit is ancestor", "commit", commit, "branch", branch, "checkCmd", checkCmd)
		tagsForCommit, err := Run("tag", "--contains", commit)
		if err != nil {
			continue
		}
		slog.Debug("Tags for commit found", "commit", commit, "tags", tagsForCommit)

		for _, tag := range splitLines(tagsForCommit) {
			slog.Debug("Tag found", "tag", tag)
			_, _, version, err := semantic.ExtractVersionFromTag(tag)
			if err != nil {
				continue
			}
			if version.IsGreaterThan(bestVersion) {
				bestVersion = version
				bestTag = tag
			}
		}
		if bestVersion.IsNotEmpty() {
			return bestTag, nil
		}
	}
	return "", fmt.Errorf("no valid tags found for branch %s", branch)
}
