package cmd

import (
	"fmt"

	"github.com/katexochen/ghh/internal/logger"
	"github.com/spf13/cobra"
)

// NewListBranchesCmd creates a new command for listing all GitHub branches of the authenticated user.
func NewListBranchesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-branches",
		Short: "List all branches where the last commit has been pushed by the authenticated user",
		RunE:  listBranches,
	}
	return cmd
}

func listBranches(cmd *cobra.Command, _ []string) error {
	flags, err := parseListBranchesFlags(cmd)
	if err != nil {
		return fmt.Errorf("get current repository: %w", err)
	}

	var log loggerI
	if flags.verbose {
		log = &logger.VerboseLogger{}
	} else {
		log = &logger.DefaultLogger{}
	}

	owner, repo, err := findOwnerAndRepo()
	if err != nil {
		return fmt.Errorf("get current repository: %w", err)
	}

	token, err := getToken()
	if err != nil {
		return fmt.Errorf("get personal access token: %w", err)
	}

	c := newGithubClient(cmd.Context(), owner, repo, token)

	user, err := c.GetUser(cmd.Context())
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}
	log.PrintJSON("user", user)

	branches, err := c.GetBranches(cmd.Context(), user.GetLogin(), 50)
	if err != nil {
		return fmt.Errorf("get branches: %w", err)
	}
	log.PrintJSON("branches", branches)

	if len(branches) == 0 {
		cmd.Println("No branches found")
		return nil
	}

	cmd.Println("Your branches:")
	for _, b := range branches {
		cmd.Printf("\t%s - %s\n", b.GetName(), b.GetCommit().GetSHA())
	}
	return nil
}

type listBranchesFlags struct {
	verbose bool
}

func parseListBranchesFlags(cmd *cobra.Command) (listBranchesFlags, error) {
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return listBranchesFlags{}, fmt.Errorf("parse verbose flag: %w", err)
	}

	return listBranchesFlags{
		verbose: verbose,
	}, nil
}
