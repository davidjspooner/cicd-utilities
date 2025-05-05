package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/davidjspooner/cicd-utilities/pkg/command"
	"github.com/davidjspooner/cicd-utilities/pkg/semantic"
)

type GithubPRUpdateOptions struct {
	PRNumber string `arg:"<pr-number>,Pull request number"`
	DryRun   bool   `arg:"--dry-run,Do not update the PR title"`
}

func init() {
	cmd := command.New(
		"github-pr-update",
		"Update GitHub PR metadata (title) based on commit messages (e.g., fix:, feat:, breaking:)",
		executeUpdateGithubPRMeta,
		&GithubPRUpdateOptions{},
	)
	commands = append(commands, cmd)
}

func executeUpdateGithubPRMeta(ctx context.Context, cmd command.Object, option *GithubPRUpdateOptions, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: update-github-pr-meta <pr-number> [--dry-run]")
	}

	if option.PRNumber == "" {
		return fmt.Errorf("pull request number is required")
	}

	token := os.Getenv("GITHUB_TOKEN")
	repo := os.Getenv("GITHUB_REPOSITORY") // e.g., "owner/repo"

	if token == "" || repo == "" {
		return fmt.Errorf("GITHUB_TOKEN and GITHUB_REPOSITORY environment variables are required")
	}

	// Fetch commit messages
	commitMessages := getCommitMessages(option.PRNumber, token, repo)

	// Determine version bump
	bump, err := semantic.Bumps.GetVersionBump(commitMessages)
	if err != nil {
		return fmt.Errorf("error determining bump : %v", err)
	}

	//get the current PR title
	prTitle, err := getPullRequestTitle(option.PRNumber, token, repo)
	if err != nil {
		return fmt.Errorf("error fetching PR title: %v", err)
	}
	if strings.Contains(prTitle, bump) {
		fmt.Printf("PR #%s already has the bump %s in the title.\n", option.PRNumber, bump)
		return nil
	}

	// Compose new title
	newTitle := fmt.Sprintf("%s: update based on commits", bump)

	if option.DryRun {
		fmt.Println("Dry run enabled.")
		fmt.Printf("Would update PR #%s with title: %s\n", option.PRNumber, newTitle)
		return nil
	}

	// Update PR via GitHub API
	return updatePullRequest(option.PRNumber, token, repo, newTitle)
}

func getCommitMessages(prNumber, token, repo string) []string {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%s/commits", repo, prNumber)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Fatalf("Failed to fetch commits: %v", err)
	}
	defer resp.Body.Close()

	var result []struct {
		Message string `json:"commit.message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Fatalf("Decode error: %v", err)
	}

	var messages []string
	for _, c := range result {
		messages = append(messages, c.Message)
	}
	return messages
}

func updatePullRequest(prNumber, token, repo, title string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%s", repo, prNumber)

	pr := struct {
		Title string `json:"title"`
	}{Title: title}
	data, _ := json.Marshal(pr)

	req, _ := http.NewRequest("PATCH", url, strings.NewReader(string(data)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to update PR: %v", err)
	}

	fmt.Printf("Updated PR #%s with new title.\n", prNumber)
	return nil
}

func getPullRequestTitle(prNumber, token, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%s", repo, prNumber)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch PR title: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode error: %v", err)
	}

	return result.Title, nil
}
