package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func findOwnerAndRepo() (string, string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	remoteURL, err := cmd.Output()
	if err != nil {
		return "", "", err
	}
	ownerRepo := strings.TrimSpace(string(remoteURL))
	ownerRepo = strings.TrimPrefix(ownerRepo, "https://github.com/")
	ownerRepo = strings.TrimSuffix(ownerRepo, ".git")
	if strings.Count(ownerRepo, "/") != 1 {
		return "", "", fmt.Errorf("invalid remote URL: %s", ownerRepo)
	}
	owner, repo, _ := strings.Cut(ownerRepo, "/")
	return owner, repo, nil
}
