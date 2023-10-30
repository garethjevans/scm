package cmd

import (
	"context"
	"fmt"
	"github.com/jenkins-x/go-scm/scm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	PreviousCommit string
	Branch         string
	SubPath        string
)

// NewMonoRepoChangeCmd creates a new token command.
func NewMonoRepoChangeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "monorepo-change",
		Short:   "Determines the sha to clone for the supplied path on the monorepository",
		Long:    "",
		Example: "scm monorepo-change --host=https://github.com --path .git-credentials",
		Aliases: []string{"t"},
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := os.ReadFile(Path)
			if err != nil {
				return err
			}
			token, err := DetermineMonoRepoChange(string(b), Kind, Host, Owner, Repo)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), token)
			return nil
		},
		Args:         cobra.NoArgs,
		SilenceUsage: true,
	}

	// common flags
	cmd.Flags().StringVarP(&Path, "path", "p", "", "The path to the git-credentials file")
	cmd.Flags().StringVarP(&Host, "host", "", "", "The host of the scm provider, including scheme")
	cmd.Flags().StringVarP(&Owner, "owner", "o", "", "The owner of the repository")
	cmd.Flags().StringVarP(&Repo, "repo", "r", "", "The name of the repository")
	cmd.Flags().StringVarP(&Kind, "kind", "k", "", "The kind of the scm provider")

	// local flags
	cmd.Flags().StringVarP(&PreviousCommit, "previous-commit", "", "", "The previous commit to search from (optional)")
	cmd.Flags().StringVarP(&Branch, "branch", "", "main", "The branch to search on (default: main)")
	cmd.Flags().StringVarP(&SubPath, "subpath", "", "", "The subPath to look for changes in (default: \"\")")

	_ = cobra.MarkFlagRequired(cmd.Flags(), "path")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "host")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "owner")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "repo")

	return cmd
}

func DetermineMonoRepoChange(credentials string, kind string, host string, owner string, repo string) (string, error) {
	repositoryURL := GetURL(kind, host, owner, repo)

	scmClient, _, _, err := GetScmClient(repositoryURL)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create scm client")
	}

	return DetermineClonePoint(scmClient, fmt.Sprintf("%s/%s", owner, repo), Branch, PreviousCommit, SubPath)
}

func DetermineClonePoint(client *scm.Client, repository string, branch string, previousCommit string, subPath string) (string, error) {
	ctx := context.Background()

	logrus.Debugf("previousCommit=%s\n", previousCommit)
	logrus.Debugf("repository=%s\n", repository)

	ref, _, err := client.Git.FindBranch(ctx, repository, branch)
	if err != nil {
		return "", err
	}

	logrus.Debugf("latest commit on branch is %s, ref=%+v\n", branch, ref)

	latestCommitOnBranch := ref.Sha

	// this is the first time we are seeing this repository, so we need to clone it all
	if previousCommit == "" {
		return latestCommitOnBranch, nil
	}

	// our understanding of the repository is up to date, so there is no work to be done
	if latestCommitOnBranch == previousCommit {
		return previousCommit, nil
	}

	// if subPath is not set, we want the whole repository
	if subPath == "" {
		return latestCommitOnBranch, nil
	}

	changes, _, err := client.Git.CompareCommits(ctx, repository, previousCommit, latestCommitOnBranch, &scm.ListOptions{})
	if err != nil {
		return "", err
	}

	logrus.Debugf("checking %d changes", len(changes))

	for _, change := range changes {
		logrus.Debugf("%s", change.Path)
		if strings.HasPrefix(change.Path, subPath) {
			return latestCommitOnBranch, nil
		}
	}

	return previousCommit, nil
}
