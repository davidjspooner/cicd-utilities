package git

import (
	"fmt"
)

func GetCurrentBranch() (string, error) {
	branch, err := Run("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %v", err)
	}
	return branch, nil
}
