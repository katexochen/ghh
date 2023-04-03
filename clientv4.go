package main

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GithubV4Client struct {
	client *githubv4.Client
	logger Logger
}

func NewGithubV4Client(ctx context.Context, token string, logger Logger) *GithubV4Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := githubv4.NewClient(tc)
	return &GithubV4Client{
		client: client,
		logger: logger,
	}
}

func (c *GithubV4Client) QueryUser(ctx context.Context, username string) (*User, error) {
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

func (c *GithubV4Client) QueryProject(ctx context.Context, org string, projectNumber int) (*Project, error) {
	var q struct {
		Organization struct {
			ProjectV2 Project `graphql:"projectV2(number: $number)"`
		} `graphql:"organization(login: $org)"`
	}

	variables := map[string]interface{}{
		"number": githubv4.Int(projectNumber),
		"org":    githubv4.String(org),
	}

	if err := c.client.Query(ctx, &q, variables); err != nil {
		return nil, err
	}

	return &q.Organization.ProjectV2, nil
}

func (c *GithubV4Client) AddProjectV2DraftIssue(ctx context.Context, input githubv4.AddProjectV2DraftIssueInput,
) (ProjectItem, error) {
	var m struct {
		AddProjectV2DraftIssue struct {
			ProjectItem ProjectItem
		} `graphql:"addProjectV2DraftIssue(input: $input)"`
	}

	return m.AddProjectV2DraftIssue.ProjectItem, c.client.Mutate(ctx, &m, input, nil)
}

func (c *GithubV4Client) UpdateProjectV2ItemFieldValueInput(ctx context.Context, project *Project,
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
		default:
			return fmt.Errorf("unsupported field type %q", field.Typename)
		}

		c.logger.PrintJSON("update project fields input", input)
		var m struct {
			UpdateProjectV2ItemFieldValue struct {
				ClientMutationId githubv4.String
			} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
		}

		if err := c.client.Mutate(ctx, &m, input, nil); err != nil {
			return err
		}

	}

	return nil
}

type User struct {
	ID githubv4.ID
}

type Project struct {
	ID     githubv4.ID
	Title  githubv4.String
	Fields struct {
		Nodes []ProjectField
	} `graphql:"fields(first: 100)"`
	URL githubv4.URI
}

type ProjectField struct {
	Typename             githubv4.String `graphql:"__typename"`
	ProjectV2FieldCommon `graphql:"... on ProjectV2FieldCommon"`
	Iteration            struct {
		Configuration struct {
			Duration githubv4.Int
			StartDay githubv4.Int
		}
	} `graphql:"... on ProjectV2IterationField"`
	SingleSelect struct {
		Options []FieldOption
	} `graphql:"... on ProjectV2SingleSelectField"`
}

type ProjectV2FieldCommon struct {
	ID       githubv4.ID
	DataType githubv4.ProjectV2FieldType
	Name     githubv4.String
}

type FieldOption struct {
	ID   githubv4.String
	Name githubv4.String
}

type ProjectItem struct {
	ID         githubv4.ID
	DatabaseID githubv4.Int
}
