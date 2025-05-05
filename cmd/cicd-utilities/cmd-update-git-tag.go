package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/davidjspooner/cicd-utilities/pkg/command"
	"github.com/davidjspooner/cicd-utilities/pkg/git"
	"github.com/davidjspooner/cicd-utilities/pkg/semantic"
)

type BumpGitTagOptions struct {
	Prefix string `arg:"--prefix,Prefix string"`
	Suffix string `arg:"--suffix,Suffix string"`
	DryRun bool   `arg:"--dry-run,Do not push the tag"`
	Remote string `arg:"--remote,Remote to push the tag to"`
}

func init() {
	cmd := command.New(
		"update-git-tag",
		"Automatically increment Git tags based on commit messages (e.g., fix:, feat:, breaking:)",
		executeBumpGitTag,
		&BumpGitTagOptions{
			Remote: "origin",
			Prefix: "v",
		},
	)
	commands = append(commands, cmd)
}

func executeBumpGitTag(ctx context.Context, cmd command.Object, option *BumpGitTagOptions, args []string) error {
	// Get the current branch
	currentBranch, err := git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}

	// Get the latest tag
	latestTag, err := getLatestTag(currentBranch)
	if err != nil {
		return fmt.Errorf("failed to get the latest tag: %v", err)
	}

	_, _, currentVersion, err := semantic.ExtractVersionFromTag(latestTag)
	if err != nil {
		return fmt.Errorf("failed to extract version from tag: %v", err)
	}

	if global.Verbose {
		fmt.Printf("Latest tag: %s\n", latestTag)
	}
	fmt.Printf("Current version: %s\n", currentVersion.String())

	// Get commit messages since the latest tag
	commitMessages, err := git.Run("log", fmt.Sprintf("%s..HEAD", latestTag), "--pretty=format:%s")
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

	fmt.Printf("Increment reason: %s\n", increment)

	if option.DryRun {
		fmt.Println("Dry run enabled.")
		fmt.Printf("Would create new tag: %s\n", newTag)
		return nil
	}

	// Create and push the new tag
	if _, err := git.Run("tag", newTag); err != nil {
		return fmt.Errorf("failed to create tag: %v", err)
	}
	if _, err := git.Run("push", option.Remote, newTag); err != nil {
		return fmt.Errorf("failed to push tag: %v", err)
	}

	fmt.Printf("Successfully created and pushed tag: %s\n", newTag)
	return nil
}

func getLatestTag(branch string) (string, error) {

	commits, err := git.Run("rev-list", "--tags", "--no-walk", "--abbrev=0", "--date-order", branch)
	if err != nil {
		return "", fmt.Errorf("failed to get latest tags: %v", err)
	}
	var bestVersion semantic.Version
	var bestTag string
	for _, commit := range splitLines(commits) {
		if commit == "" {
			continue
		}
		commit = strings.TrimSpace(commit)
		checkCmd, err := git.Run("merge-base", "--is-ancestor", commit, branch)
		if err != nil {
			continue
		}
		if global.Verbose {
			fmt.Printf("Commit %s is ancestor of %s: %s\n", commit, branch, checkCmd)
		}
		tagsForCommit, err := git.Run("tag", "--contains", commit)
		if err != nil {
			continue
		}
		for _, tag := range splitLines(tagsForCommit) {
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
