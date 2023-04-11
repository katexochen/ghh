package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/katexochen/ghh/internal/logger"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
)

// NewCreateProjectIssueCmd creates a new command for creating a project issue.
func NewCreateProjectIssueCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-project-issue",
		Short: "Create a project issue",
		RunE:  createProjectIssue,
	}
	cmd.Flags().String("metadata", "", "Path to metadata file")
	cmd.Flags().String("body", "", "Path to body file")
	return cmd
}

type metadata struct {
	Organization  string // required
	ProjectNumber int    // required

	IssueTitle string
	Assignees  []string
	Fields     map[string]string
}

func createProjectIssue(cmd *cobra.Command, _ []string) error {
	flags, err := parseCreateProjectIssueFlags(cmd)
	if err != nil {
		return err
	}

	var log loggerI
	if flags.verbose {
		log = &logger.VerboseLogger{}
	} else {
		log = &logger.DefaultLogger{}
	}

	token, err := getToken()
	if err != nil {
		return fmt.Errorf("getting token: %w", err)
	}

	c := newGithubV4Client(cmd.Context(), token, log)

	c.logger.Debugf("searching project %s/%d", flags.Metadata.Organization, flags.Metadata.ProjectNumber)
	project, err := c.QueryProject(cmd.Context(), flags.Metadata.Organization, flags.Metadata.ProjectNumber)
	if err != nil {
		return fmt.Errorf("querying project: %w", err)
	}
	c.logger.PrintJSON("found project", project)

	var assigneeIDs []githubv4.ID
	for _, assignee := range flags.Metadata.Assignees {
		c.logger.Debugf("searching user %s", assignee)
		user, err := c.QueryUser(cmd.Context(), assignee)
		if err != nil {
			return fmt.Errorf("querying user: %w", err)
		}
		c.logger.PrintJSON("found user", user)
		assigneeIDs = append(assigneeIDs, user.ID)
	}

	addDraftIssueInput := githubv4.AddProjectV2DraftIssueInput{
		ProjectID: project.ID,
		Title:     githubv4.String(flags.Metadata.IssueTitle),
	}
	if flags.Body != "" {
		addDraftIssueInput.Body = toPtr(githubv4.String(flags.Body))
	}
	if len(flags.Metadata.Assignees) > 0 {
		addDraftIssueInput.AssigneeIDs = toPtr(assigneeIDs)
	}

	item, err := c.AddProjectV2DraftIssue(cmd.Context(), addDraftIssueInput)
	if err != nil {
		return fmt.Errorf("adding project issue: %w", err)
	}

	if err := c.UpdateProjectV2ItemFieldValueInput(cmd.Context(), project, item.ID, flags.Metadata.Fields); err != nil {
		return fmt.Errorf("updating project issue fields: %w", err)
	}

	itemURL := fmt.Sprintf("%s?pane=issue&itemId=%d", project.URL, item.DatabaseID)
	c.logger.Infof("created project issue:")
	fmt.Println(itemURL)

	return nil
}

type createProjectIssueFlags struct {
	Metadata metadata
	Body     string
	verbose  bool
}

func parseCreateProjectIssueFlags(cmd *cobra.Command) (createProjectIssueFlags, error) {
	metadataPath, err := cmd.Flags().GetString("metadata")
	if err != nil {
		return createProjectIssueFlags{}, err
	}
	metadataBytes, err := os.ReadFile(metadataPath)
	if err != nil {
		return createProjectIssueFlags{}, err
	}
	var metadata metadata
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		return createProjectIssueFlags{}, err
	}
	if metadata.Organization == "" {
		return createProjectIssueFlags{}, errors.New("validating metadata fields: organization is required")
	}
	if metadata.ProjectNumber == 0 {
		return createProjectIssueFlags{}, errors.New("validating metadata fields: project number is required")
	}

	bodyPath, err := cmd.Flags().GetString("body")
	if err != nil {
		return createProjectIssueFlags{}, err
	}
	var bodyBytes []byte
	if bodyPath != "" {
		bodyBytes, err = os.ReadFile(bodyPath)
		if err != nil {
			return createProjectIssueFlags{}, err
		}
	}

	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return createProjectIssueFlags{}, err
	}

	return createProjectIssueFlags{
		Metadata: metadata,
		Body:     string(bodyBytes),
		verbose:  verbose,
	}, nil
}

func toPtr[T any](v T) *T {
	return &v
}
