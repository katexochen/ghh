package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v50/github"
	"github.com/manifoldco/promptui"
)

func DeleteRuns(ctx context.Context) error {
	owner, repo, err := findOwnerAndRepo()
	if err != nil {
		return err
	}

	token, err := GetToken()
	if err != nil {
		return err
	}

	c := newGithubClient(ctx, owner, repo, token)

	workflows, err := c.GetWorkflows(ctx)
	if err != nil {
		return err
	}

	workflow, err := selectWorkflow(workflows)
	if err != nil {
		return err
	}

	runs, err := c.GetWorkflowRuns(ctx, workflow.GetID())
	if err != nil {
		return err
	}

	fmt.Printf("Deleting %d runs...\n", len(runs))
	if err := c.DeleteWorkflowRuns(ctx, runs); err != nil {
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
