package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
)

func newCreateProjectIssueCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-project-issue",
		Short: "Create a project issue",
		RunE:  CreateProjectIssue,
	}
	cmd.Flags().String("metadata", "", "Path to metadata file")
	cmd.Flags().String("body", "", "Path to body file")
	return cmd
}

type Metadata struct {
	Organization  string
	ProjectNumber int

	IssueTitle string
	Assignees  []string
	Fields     map[string]string
}

func CreateProjectIssue(cmd *cobra.Command, _ []string) error {
	flags, err := parseCreateProjectIssueFlags(cmd)
	if err != nil {
		return err
	}

	token, err := GetToken()
	if err != nil {
		return err
	}

	c := NewGithubV4Client(cmd.Context(), token)

	project, err := c.QueryProject(cmd.Context(), flags.Metadata.Organization, flags.Metadata.ProjectNumber)
	if err != nil {
		return err
	}
	c.logger.PrintJSON("found project", project)

	var assigneeIDs []githubv4.ID
	for _, assignee := range flags.Metadata.Assignees {
		user, err := c.QueryUser(cmd.Context(), assignee)
		if err != nil {
			return err
		}
		c.logger.PrintJSON("found user", user)
		assigneeIDs = append(assigneeIDs, user.ID)
	}

	addDraftIssueInput := githubv4.AddProjectV2DraftIssueInput{
		ProjectID:   project.ID,
		Title:       githubv4.String(flags.Metadata.IssueTitle),
		Body:        toPtr(githubv4.String(flags.Body)),
		AssigneeIDs: toPtr(assigneeIDs),
	}
	item, err := c.AddProjectV2DraftIssue(cmd.Context(), addDraftIssueInput)
	if err != nil {
		return err
	}

	if err := c.UpdateProjectV2ItemFieldValueInput(cmd.Context(), project, item.ID, flags.Metadata.Fields); err != nil {
		return err
	}

	itemURL := fmt.Sprintf("%s?pane=issue&itemId=%d", project.URL, item.DatabaseID)
	c.logger.Infof("created project issue %s", itemURL)

	return nil
}

type CreateProjectIssueFlags struct {
	Metadata Metadata
	Body     string
}

func parseCreateProjectIssueFlags(cmd *cobra.Command) (CreateProjectIssueFlags, error) {
	metadataPath, err := cmd.Flags().GetString("metadata")
	if err != nil {
		return CreateProjectIssueFlags{}, err
	}
	metadataBytes, err := os.ReadFile(metadataPath)
	if err != nil {
		return CreateProjectIssueFlags{}, err
	}
	var metadata Metadata
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		return CreateProjectIssueFlags{}, err
	}

	bodyPath, err := cmd.Flags().GetString("body")
	if err != nil {
		return CreateProjectIssueFlags{}, err
	}
	bodyBytes, err := os.ReadFile(bodyPath)
	if err != nil {
		return CreateProjectIssueFlags{}, err
	}

	return CreateProjectIssueFlags{
		Metadata: metadata,
		Body:     string(bodyBytes),
	}, nil
}

func toPtr[T any](v T) *T {
	return &v
}
