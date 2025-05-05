package main

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/davidjspooner/cicd-utilities/pkg/command"
)

type ChecksumOptions struct {
	Algorithm string `arg:"--algorithm,Checksum algorithm (e.g., sha256, md5)"`
	Extension string `arg:"--extension,File extension override for the checksum file"`
}

func init() {
	cmd := command.New(
		"checksum",
		"Generate checksum(s) for file(s) using a specified algorithm",
		executeChecksum,
		&ChecksumOptions{
			Algorithm: "sha256",
		},
	)
	commands = append(commands, cmd)
}

func executeChecksum(ctx context.Context, cmd command.Object, option *ChecksumOptions, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("no files specified")
	}

	files, err := globFiles(args)
	if err != nil {
		return fmt.Errorf("error globbing files: %s", err)
	}
	extension := option.Extension
	if extension == "" {
		extension = option.Algorithm
	}
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}

	for _, file := range files {
		checksum, err := generateChecksum(file, option.Algorithm, extension)
		if err != nil {
			return fmt.Errorf("error generating checksum for %s: %v", file, err)
		}
		fmt.Printf("%s: %s\n", file, checksum)
	}

	return nil
}

func generateChecksum(file string, algorithm, extension string) (string, error) {
	stat, err := os.Stat(file)
	if err != nil {
		return "", fmt.Errorf("failed to stat file: %v", err)
	}
	if stat.IsDir() {
		return "", fmt.Errorf("file is a directory: %s", file)
	}
	if stat.Size() == 0 {
		return "", fmt.Errorf("file is empty: %s", file)
	}
	input, err := os.Open(file)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	var checksum string
	defer input.Close()
	switch algorithm {
	case "sha256":
		//generate md5 checksum
		h := sha256.New()
		if _, err := io.Copy(h, input); err != nil {
			return "", fmt.Errorf("failed to generate checksum: %v", err)
		}
		checksum = fmt.Sprintf("%x", h.Sum(nil))
	case "md5":
		//generate md5 checksum
		h := md5.New()
		if _, err := io.Copy(h, input); err != nil {
			return "", fmt.Errorf("failed to generate checksum: %v", err)
		}
		checksum = fmt.Sprintf("%x", h.Sum(nil))
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
	if checksum == "" {
		return "", fmt.Errorf("failed to generate checksum")
	}
	output, err := os.Create(file + extension)
	if err != nil {
		return "", fmt.Errorf("failed to create checksum file: %v", err)
	}
	defer output.Close()
	_, err = output.WriteString(checksum)
	return checksum, err
}
