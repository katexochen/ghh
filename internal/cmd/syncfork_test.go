package cmd

import (
	"strconv"
	"testing"

	"github.com/google/go-github/v63/github"
	"github.com/stretchr/testify/assert"
)

func TestFilterIgnoredRepos(t *testing.T) {
	testCases := []struct {
		repos       []*github.Repository
		ignoreRepos []string
		expected    []*github.Repository
	}{
		{
			repos: []*github.Repository{
				{Name: toPtr("foo")},
				{Name: toPtr("bar")},
				{Name: toPtr("baz")},
			},
			ignoreRepos: []string{"bar"},
			expected: []*github.Repository{
				{Name: toPtr("foo")},
				{Name: toPtr("baz")},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert := assert.New(t)

			filtered := filterIgnoredRepos(tc.repos, tc.ignoreRepos)
			assert.Equal(tc.expected, filtered)
		})
	}
}
