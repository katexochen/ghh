package cmd

import "github.com/shurcooL/githubv4"

// User is a GitHub user, see https://docs.github.com/en/graphql/reference/objects#user.
type User struct {
	ID githubv4.ID
}

// Project is a GitHub project, see https://docs.github.com/en/graphql/reference/objects#projectv2.
type Project struct {
	ID     githubv4.ID
	Title  githubv4.String
	Fields struct {
		Nodes []ProjectField
	} `graphql:"fields(first: 100)"`
	URL githubv4.URI
}

// ProjectField is a GitHub project field, see https://docs.github.com/en/graphql/reference/objects#projectv2field.
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
		Options []ProjectSingleSelectFieldOption
	} `graphql:"... on ProjectV2SingleSelectField"`
}

// ProjectV2FieldCommon is a GitHub project field common, see https://docs.github.com/en/graphql/reference/interfaces#projectv2fieldcommon.
type ProjectV2FieldCommon struct {
	ID       githubv4.ID
	DataType githubv4.ProjectV2FieldType
	Name     githubv4.String
}

// ProjectSingleSelectFieldOption is a GitHub field option, see https://docs.github.com/en/graphql/reference/objects#projectv2singleselectfieldoption.
type ProjectSingleSelectFieldOption struct {
	ID   githubv4.String
	Name githubv4.String
}

// ProjectItem is a GitHub project item, see https://docs.github.com/en/graphql/reference/objects#projectv2item.
type ProjectItem struct {
	ID         githubv4.ID
	DatabaseID githubv4.Int
}
