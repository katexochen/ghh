package cmd

import (
	"context"

	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
)

type githubClient struct {
	client *github.Client
	owner  string
	repo   string
}

func newGithubClient(ctx context.Context, owner, repo, token string) *githubClient {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return &githubClient{
		client: github.NewClient(tc),
		owner:  owner,
		repo:   repo,
	}
}

func (c *githubClient) GetWorkflows(ctx context.Context) ([]*github.Workflow, error) {
	opt := &github.ListOptions{PerPage: 1000}
	var allWorkflows []*github.Workflow
	for {
		workflows, resp, err := c.client.Actions.ListWorkflows(ctx, c.owner, c.repo, opt)
		if err != nil {
			return nil, err
		}
		allWorkflows = append(allWorkflows, workflows.Workflows...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allWorkflows, nil
}

func (c *githubClient) GetWorkflowRuns(ctx context.Context, workflowID int64) ([]*github.WorkflowRun, error) {
	opt := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{PerPage: 1000},
	}
	var allRuns []*github.WorkflowRun
	for {
		runs, resp, err := c.client.Actions.ListWorkflowRunsByID(ctx, c.owner, c.repo, workflowID, opt)
		if err != nil {
			return nil, err
		}
		allRuns = append(allRuns, runs.WorkflowRuns...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRuns, nil
}

func (c *githubClient) DeleteWorkflowRuns(ctx context.Context, runs []*github.WorkflowRun) error {
	for _, run := range runs {
		_, err := c.client.Actions.DeleteWorkflowRun(ctx, c.owner, c.repo, run.GetID())
		if err != nil {
			return err
		}
	}
	return nil
}
