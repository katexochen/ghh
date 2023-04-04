package cmd

import (
	"context"
	"fmt"

	"github.com/google/go-github/v50/github"
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

// getBranches returns all branches with the last commit being authored by the specified author.
// Limit is the per-page limit for the GitHub API. (capped at 100)
func (c *githubClient) GetBranches(ctx context.Context, author string, limit int) ([]*github.Branch, error) {
	if limit > 100 {
		// max limit as per https://docs.github.com/en/rest/branches/branches?apiVersion=2022-11-28
		limit = 100
	}

	opt := &github.BranchListOptions{
		Protected: github.Bool(false),
		ListOptions: github.ListOptions{
			PerPage: limit,
		},
	}

	var allBranches []*github.Branch
	for {
		branches, resp, err := c.client.Repositories.ListBranches(ctx, c.owner, c.repo, opt)
		if err != nil {
			return nil, fmt.Errorf("list branches: %w", err)
		}
		allBranches = append(allBranches, branches...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var filteredBranches []*github.Branch
	for _, branch := range allBranches {
		commit, _, err := c.client.Repositories.GetCommit(ctx, c.owner, c.repo, branch.GetCommit().GetSHA(), nil)
		if err != nil {
			return nil, fmt.Errorf("get commit: %w", err)
		}
		if commit.GetAuthor().GetLogin() == author {
			filteredBranches = append(filteredBranches, branch)
		}
	}

	return filteredBranches, nil
}

func (c *githubClient) GetUser(ctx context.Context) (*github.User, error) {
	user, _, err := c.client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}
