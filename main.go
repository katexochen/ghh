package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	args := os.Args[1:]
	if len(args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}

	switch args[0] {
	case "setauth":
		return SetupAuth()
	case "delete-all-runs":
		return DeleteRuns(ctx)
	default:
		return fmt.Errorf("invalid command: %s", args[0])
	}
}

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
