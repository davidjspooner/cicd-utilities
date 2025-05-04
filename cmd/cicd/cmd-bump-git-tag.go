package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/davidjspooner/go-cicd/pkg/git"
	"github.com/davidjspooner/go-cicd/pkg/semantic"
)

func init() {
	registerCommand(
		"bump-git-tag",
		"Bumps the Git tag based on commit messages",
		executeBumpGitTag,
	)
}

func executeBumpGitTag(args []string) error {
	bumpCommand := flag.NewFlagSet("bump-git-tag", flag.ExitOnError)
	prefix := bumpCommand.String("prefix", "", "Prefix for the new tag")
	suffix := bumpCommand.String("suffix", "", "Suffix for the new tag")
	dryRun := bumpCommand.Bool("dry-run", false, "Perform a dry run without updating the tag")
	verbose := bumpCommand.Bool("verbose", false, "Enable verbose output")
	remote := bumpCommand.String("remote", "origin", "Remote to push the tag to")
	bumpCommand.Parse(args)

	// Get the current branch
	currentBranch, err := git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("failed to get current branch: %v", err)
	}

	// Get the latest tag
	latestTag, err := getLatestTag(currentBranch, *verbose)
	if err != nil {
		return fmt.Errorf("failed to get the latest tag: %v", err)
	}

	_, _, currentVersion, err := semantic.ExtractVersionFromTag(latestTag)
	if err != nil {
		return fmt.Errorf("failed to extract version from tag: %v", err)
	}

	if *verbose {
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
		fmt.Println("No changes deteced, no version bump needed.")
		return nil
	}

	// Determine the version bump
	bump, err := semantic.Bumps.GetVersionBump(commits)
	if err != nil {
		return fmt.Errorf("failed to determine version bump: %v", err)
	}

	// Increment the version
	newVersion, err := currentVersion.Increment(bump)
	if err != nil {
		return fmt.Errorf("failed to increment version: %v", err)
	}

	// Construct the new tag
	newTag := fmt.Sprintf("%s%s%s", *prefix, newVersion.String(), *suffix)

	fmt.Printf("Bumping reason: %s\n", bump)

	if *dryRun {
		fmt.Println("Dry run enabled.")
		fmt.Printf("Would create new tag: %s\n", newTag)
		return nil
	}

	// Create and push the new tag
	if _, err := git.Run("tag", newTag); err != nil {
		return fmt.Errorf("failed to create tag: %v", err)
	}
	if _, err := git.Run("push", *remote, newTag); err != nil {
		return fmt.Errorf("failed to push tag: %v", err)
	}

	fmt.Printf("Successfully created and pushed tag: %s\n", newTag)
	return nil
}

func getLatestTag(branch string, verbose bool) (string, error) {

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
		if verbose {
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
