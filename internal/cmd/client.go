package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-github/v63/github"
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

func (c *githubClient) GetUserRepositories(ctx context.Context) ([]*github.Repository, error) {
	opt := &github.RepositoryListByAuthenticatedUserOptions{
		Affiliation: "owner",
		ListOptions: github.ListOptions{PerPage: 1000},
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := c.client.Repositories.ListByAuthenticatedUser(ctx, opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}

func (c *githubClient) GetUserForks(ctx context.Context) ([]*github.Repository, error) {
	userRepos, err := c.GetUserRepositories(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing user repositories: %w", err)
	}

	var userForks []*github.Repository
	for _, repo := range userRepos {
		if repo.Fork == nil {
			continue
		}
		if *repo.Fork {
			userForks = append(userForks, repo)
		}
	}

	return userForks, nil
}

func (c *githubClient) SyncFork(ctx context.Context, repo *github.Repository, branch string) (*github.RepoMergeUpstreamResult, error) {
	if repo.Fork == nil {
		return nil, errors.New("repo is not a fork")
	}

	req := &github.RepoMergeUpstreamRequest{
		Branch: &branch,
	}
	result, resp, err := c.client.Repositories.MergeUpstream(ctx, repo.GetOwner().GetLogin(), repo.GetName(), req)
	if err != nil {
		return nil, fmt.Errorf("merging upstream into fork: %w", err)
	}

	switch resp.StatusCode {
	case 200:
		return result, nil
	case 409:
		return nil, errors.New("The branch could not be synced because of a merge conflict")
	case 422:
		return nil, errors.New("The branch could not be synced for some other reason")
	default:
		return nil, fmt.Errorf("An unknown statuscode was returned: %s", resp.Status)
	}
}

func (c *githubClient) GetBranch(ctx context.Context, repo *github.Repository, branch string) (*github.Branch, error) {
	result, resp, err := c.client.Repositories.GetBranch(ctx, repo.GetOwner().GetLogin(), repo.GetName(), branch, 10)
	if err != nil {
		return nil, fmt.Errorf("getting branch: %w", err)
	}

	switch resp.StatusCode {
	case 200:
		return result, nil
	case 404:
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("An unknown statuscode was returned: %s", resp.Status)
	}
}

// ErrNotFound is returned when a resource is not found.
var ErrNotFound = errors.New("resource not found")
