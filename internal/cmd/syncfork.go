package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-github/v57/github"
	"github.com/katexochen/ghh/internal/logger"
	"github.com/spf13/cobra"
)

// NewSyncForksCmd creates a new command for syncing forks.
func NewSyncForksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync-forks",
		Short: "Sync all forks of a user",
		Long: `
This command will sync all forks of a user with their upstream repository. It will
fast-forward the fork if possible, otherwise it will merge the upstream branch.

Per default, the target of the merge is the default branch of the fork.
		`,
		RunE: syncForks,
	}

	cmd.Flags().StringSliceP(
		"target-branches",
		"t",
		[]string{},
		"Target branches to sync. If empty, the default branch of the fork will be used.",
	)
	cmd.Flags().StringSliceP(
		"ignore-repos",
		"i",
		[]string{},
		"Repositories to ignore.",
	)
	cmd.Flags().Bool(
		"dont-target-default",
		false,
		"Don't target the default branch of the fork. 'target-branches' must be set. "+
			"If non of the target branches exist, the repo will not be synced.",
	)

	return cmd
}

func syncForks(cmd *cobra.Command, _ []string) error {
	flags, err := parseSyncForksFlags(cmd)
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

	c := newGithubClient(cmd.Context(), "", "", token) // TODO: refactor client to not have owner and repo

	log.Debugf("listing forks")
	forks, err := c.GetUserForks(cmd.Context())
	if err != nil {
		return fmt.Errorf("listing forks: %w", err)
	}

	for _, fork := range forks {
		log.Debugf("discovered fork %s", fork.GetFullName())
	}
	log.Debugf("%d forks found", len(forks))

	forks = filterIgnoredRepos(forks, flags.ignoreRepos)
	log.Debugf("%d remaining after filtering", len(forks))

	var retErr error
	var fastforwarded, merged, notbehind, skipped, failed int
	for _, fork := range forks {
		var branch string
		if !flags.dontTargetDefault {
			branch = fork.GetDefaultBranch()
			log.Debugf("%s: default branch is %s", fork.GetFullName(), branch)
		}

		for _, targetBranch := range flags.targetBranches {
			log.Debugf("%s: checking if branch %q exists", fork.GetFullName(), targetBranch)
			if _, err := c.GetBranch(cmd.Context(), fork, targetBranch); err == nil {
				branch = targetBranch
				log.Debugf("%s: using target branch %q", fork.GetFullName(), branch)
				break
			}
		}

		if branch == "" {
			log.Warnf("%s: no target branch found, skipping", fork.GetFullName())
			skipped++
			continue
		}

		log.Infof("%s: syncing fork branch %q with upstream", fork.GetFullName(), branch)
		result, err := c.SyncFork(cmd.Context(), fork, branch)
		if errors.Is(err, context.Canceled) {
			return err
		}
		if err != nil {
			log.Errorf("%s: syncing fork: %s", fork.GetFullName(), err)
			retErr = errors.Join(retErr, fmt.Errorf("syncing fork %s: %w", fork.GetFullName(), err))
			failed++
			continue
		}

		log.Infof("synced fork %s: %s", fork.GetFullName(), result.GetMessage())

		switch result.GetMergeType() {
		case "fast-forward":
			fastforwarded++
		case "merge":
			merged++
		case "none":
			notbehind++
		}
	}

	log.Infof(
		"synced %d forks: %d up-to-date, %d fast-forwarded, %d merged, %d skipped, %d failed",
		len(forks), notbehind, fastforwarded, merged, skipped, failed,
	)
	return retErr
}

func filterIgnoredRepos(repos []*github.Repository, ignoreRepos []string) []*github.Repository {
	var filtered []*github.Repository
outer:
	for _, repo := range repos {
		for _, ignore := range ignoreRepos {
			if strings.EqualFold(repo.GetName(), ignore) {
				continue outer
			}
		}
		filtered = append(filtered, repo)
	}
	return filtered
}

type syncForksFlags struct {
	verbose           bool
	ignoreRepos       []string
	targetBranches    []string
	dontTargetDefault bool
}

func parseSyncForksFlags(cmd *cobra.Command) (*syncForksFlags, error) {
	flags := &syncForksFlags{}

	var err error
	flags.verbose, err = cmd.Flags().GetBool("verbose")
	if err != nil {
		return nil, err
	}
	flags.ignoreRepos, err = cmd.Flags().GetStringSlice("ignore-repos")
	if err != nil {
		return nil, err
	}
	flags.targetBranches, err = cmd.Flags().GetStringSlice("target-branches")
	if err != nil {
		return nil, err
	}
	flags.dontTargetDefault, err = cmd.Flags().GetBool("dont-target-default")
	if err != nil {
		return nil, err
	}

	if flags.dontTargetDefault && len(flags.targetBranches) == 0 {
		return nil, errors.New("'--target-branches' must be set when using '--dont-target-default'")
	}

	return flags, nil
}
