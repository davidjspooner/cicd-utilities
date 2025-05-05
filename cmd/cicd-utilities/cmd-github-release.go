package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/davidjspooner/cicd-utilities/pkg/command"
)

type GithubReleaseOptions struct {
	TagName    string `arg:"--tag,Tag name for the release"`
	Name       string `arg:"--name,Name/Title of the release"`
	Body       string `arg:"--body,Description of the release"`
	Draft      bool   `arg:"--draft,Create the release as a draft"`
	Prerelease bool   `arg:"--prerelease,Mark the release as a prerelease"`
}

func init() {
	cmd := command.New(
		"github-release",
		"Create a GitHub release",
		executeGithubRelease,
		&GithubReleaseOptions{},
	)
	commands = append(commands, cmd)
}

func executeGithubRelease(ctx context.Context, cmd command.Object, option *GithubReleaseOptions, args []string) error {
	err := command.CheckUnparsedOptions(args)
	if err != nil {
		return err
	}

	files, err := globFiles(args)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no files found matching the pattern")
	}

	token := os.Getenv("GITHUB_TOKEN")
	repo := os.Getenv("GITHUB_REPOSITORY") // e.g., "owner/repo"

	if token == "" || repo == "" {
		return fmt.Errorf("GITHUB_TOKEN and GITHUB_REPOSITORY environment variables are required")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/releases", repo)

	release := struct {
		TagName    string `json:"tag_name,omitempty"`
		Name       string `json:"name"`
		Body       string `json:"body"`
		Draft      bool   `json:"draft"`
		Prerelease bool   `json:"prerelease"`
	}{
		TagName:    option.TagName,
		Name:       option.Name,
		Body:       option.Body,
		Draft:      true,
		Prerelease: option.Prerelease,
	}

	data, err := json.Marshal(release)
	if err != nil {
		return fmt.Errorf("failed to marshal release data: %v", err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create release: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("failed to create release, status code: %d", resp.StatusCode)
	}

	// Parse the response to get the release ID
	var releaseResponse struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releaseResponse); err != nil {
		return fmt.Errorf("failed to parse release response: %v", err)
	}
	releaseID := releaseResponse.ID

	// Logic to upload the file to the release
	for _, file := range files {
		uploadFileToGubHubRelease(file, releaseID, token, repo)
	}

	fmt.Println("GitHub release created successfully.")
	return nil
}

func uploadFileToGubHubRelease(file string, releaseID int, token string, repo string) {
	url := fmt.Sprintf("https://uploads.github.com/repos/%s/releases/%d/assets?name=%s", repo, releaseID, file)

	fileData, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("failed to read file %s: %v\n", file, err)
		return
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(fileData)))
	if err != nil {
		fmt.Printf("failed to create request for file %s: %v\n", file, err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("failed to upload file %s: %v\n", file, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		fmt.Printf("failed to upload file %s, status code: %d\n", file, resp.StatusCode)
		return
	}

	fmt.Printf("File %s uploaded successfully.\n", file)
}
