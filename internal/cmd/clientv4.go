package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type githubV4Client struct {
	client *githubv4.Client
	logger loggerI
}

func newGithubV4Client(ctx context.Context, token string, logger loggerI) *githubV4Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := githubv4.NewClient(tc)
	return &githubV4Client{
		client: client,
		logger: logger,
	}
}

func (c *githubV4Client) QueryUser(ctx context.Context, username string) (*User, error) {
	var q struct {
		User User `graphql:"user(login: $user)"`
	}

	variables := map[string]interface{}{
		"user": githubv4.String(username),
	}

	if err := c.client.Query(ctx, &q, variables); err != nil {
		return nil, err
	}

	return &q.User, nil
}

func (c *githubV4Client) QueryProject(ctx context.Context, owner string, isOrg bool, projectNumber int) (*Project, error) {
	if isOrg {
		var q struct {
			Organization struct {
				ProjectV2 Project `graphql:"projectV2(number: $number)"`
			} `graphql:"organization(login: $org)"`
		}

		variables := map[string]interface{}{
			"number": githubv4.Int(projectNumber),
			"org":    githubv4.String(owner),
		}

		if err := c.client.Query(ctx, &q, variables); err != nil {
			return nil, err
		}

		return &q.Organization.ProjectV2, nil
	}

	var q struct {
		User struct {
			ProjectV2 Project `graphql:"projectV2(number: $number)"`
		} `graphql:"user(login: $user)"`
	}

	variables := map[string]interface{}{
		"number": githubv4.Int(projectNumber),
		"user":   githubv4.String(owner),
	}

	if err := c.client.Query(ctx, &q, variables); err != nil {
		return nil, err
	}

	return &q.User.ProjectV2, nil
}

func (c *githubV4Client) AddProjectV2DraftIssue(ctx context.Context, input githubv4.AddProjectV2DraftIssueInput,
) (ProjectItem, error) {
	var m struct {
		AddProjectV2DraftIssue struct {
			ProjectItem ProjectItem
		} `graphql:"addProjectV2DraftIssue(input: $input)"`
	}

	return m.AddProjectV2DraftIssue.ProjectItem, c.client.Mutate(ctx, &m, input, nil)
}

func (c *githubV4Client) UpdateProjectV2ItemFieldValueInput(ctx context.Context, project *Project,
	itemID githubv4.ID, fieldValues map[string]string,
) error {
	for fieldName, value := range fieldValues {
		input := githubv4.UpdateProjectV2ItemFieldValueInput{
			ProjectID: project.ID,
			ItemID:    itemID,
		}

		var field ProjectField
		for i, f := range project.Fields.Nodes {
			if f.Name == githubv4.String(fieldName) {
				field = project.Fields.Nodes[i]
				input.FieldID = field.ID
			}
		}
		if field.Typename == "" {
			return fmt.Errorf("field %q not found", fieldName)
		}
		c.logger.PrintJSON("found field", field)

		switch field.DataType {
		case githubv4.ProjectV2FieldTypeSingleSelect:
			for i, option := range field.SingleSelect.Options {
				if option.Name == githubv4.String(value) {
					input.Value = githubv4.ProjectV2FieldValue{
						SingleSelectOptionID: &field.SingleSelect.Options[i].ID,
					}
				}
			}
			if input.Value.SingleSelectOptionID == nil {
				return fmt.Errorf("option %q not found", value)
			}
		case githubv4.ProjectV2FieldTypeText:
			input.Value = githubv4.ProjectV2FieldValue{
				Text: toPtr(githubv4.String(value)),
			}
		case githubv4.ProjectV2FieldTypeDate:
			input.Value = githubv4.ProjectV2FieldValue{
				Date: toPtr(githubv4.Date{Time: parseDate(value)}),
			}
		default:
			return fmt.Errorf("unsupported field type %q", field.Typename)
		}

		c.logger.PrintJSON("update project fields input", input)
		var m struct {
			UpdateProjectV2ItemFieldValue struct {
				ClientMutationID githubv4.String
			} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
		}

		if err := c.client.Mutate(ctx, &m, input, nil); err != nil {
			return err
		}

	}

	return nil
}

func parseDate(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
