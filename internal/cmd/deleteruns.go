package cmd

import (
	"fmt"

	"github.com/google/go-github/v53/github"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// NewDeleteAllRunsCmd creates a new command for deleting workflow runs.
func NewDeleteAllRunsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-all-runs",
		Short: "Delete all workflow runs",
		RunE:  deleteRuns,
	}
	return cmd
}

func deleteRuns(cmd *cobra.Command, _ []string) error {
	owner, repo, err := findOwnerAndRepo()
	if err != nil {
		return err
	}

	token, err := getToken()
	if err != nil {
		return err
	}

	c := newGithubClient(cmd.Context(), owner, repo, token)

	workflows, err := c.GetWorkflows(cmd.Context())
	if err != nil {
		return err
	}

	workflow, err := selectWorkflow(workflows)
	if err != nil {
		return err
	}

	runs, err := c.GetWorkflowRuns(cmd.Context(), workflow.GetID())
	if err != nil {
		return err
	}

	fmt.Printf("Deleting %d runs...\n", len(runs))
	if err := c.DeleteWorkflowRuns(cmd.Context(), runs); err != nil {
		return err
	}

	fmt.Println("Done.")

	return nil
}

func selectWorkflow(workflows []*github.Workflow) (*github.Workflow, error) {
	names := workflowNames(workflows)
	prompt := promptui.Select{
		Label: "Select workflow",
		Items: names,
	}
	idx, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return workflows[idx], nil
}

func workflowNames(workflows []*github.Workflow) []string {
	var names []string
	for _, w := range workflows {
		names = append(names, w.GetName())
	}
	return names
}
