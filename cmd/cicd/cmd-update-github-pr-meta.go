package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/davidjspooner/go-cicd/pkg/semantic"
)

func init() {
	registerCommand(
		"update-github-pr-meta",
		"Updates GitHub PR metadata based on commit messages",
		executeUpdateGithubPRMeta,
	)
}

func executeUpdateGithubPRMeta(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: update-github-pr-meta <pr-number> [--dry-run]")
	}

	subCommand := flag.NewFlagSet("update-github-pr-meta", flag.ExitOnError)
	prNumber := subCommand.String("pr", "", "Pull request number")
	dryRun := subCommand.Bool("dry-run", false, "Perform a dry run without updating the PR")
	err := subCommand.Parse(args)
	if err != nil {
		return fmt.Errorf("error parsing arguments: %v", err)
	}
	if *prNumber == "" {
		return fmt.Errorf("pull request number is required")
	}

	token := os.Getenv("GITHUB_TOKEN")
	repo := os.Getenv("GITHUB_REPOSITORY") // e.g., "owner/repo"

	if token == "" || repo == "" {
		return fmt.Errorf("GITHUB_TOKEN and GITHUB_REPOSITORY environment variables are required")
	}

	// Fetch commit messages
	commitMessages := getCommitMessages(*prNumber, token, repo)

	// Determine version bump
	bump, err := semantic.Bumps.GetVersionBump(commitMessages)
	if err != nil {
		return fmt.Errorf("error determining version bump: %v", err)
	}

	// Compose new title and body
	newTitle := fmt.Sprintf("%s: update based on commits", bump)
	newBody := strings.Join(commitMessages, "\n")

	if *dryRun {
		fmt.Println("Dry run enabled.")
		fmt.Printf("Would update PR #%s with title: %s\n", *prNumber, newTitle)
		fmt.Println("New body:")
		fmt.Println(newBody)
		return nil
	}

	// Update PR via GitHub API
	return updatePullRequest(*prNumber, token, repo, newTitle, newBody)
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

func updatePullRequest(prNumber, token, repo, title, body string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%s", repo, prNumber)

	pr := struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}{Title: title, Body: body}
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

	fmt.Printf("Updated PR #%s with new title and description.\n", prNumber)
	return nil
}
