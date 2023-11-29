package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	Commit string
	Name   string
	Status string
)

// NewCheckRunCmd creates a new token command.
func NewCheckRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "check-run",
		Short:   "Creates/Updates a checkrun for a repository and sha",
		Long:    "",
		Example: "scm check-run --host=https://github.com --path .git-credentials",
		Aliases: []string{"c"},
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := os.ReadFile(Path)
			if err != nil {
				return err
			}
			url, err := UpdateCheckRun(string(b), Kind, Host, Owner, Repo)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), url)
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
	cmd.Flags().StringVarP(&Commit, "commit", "", "", "The commit to attach the check run to")
	cmd.Flags().StringVarP(&Name, "name", "", "", "The name of the check run")

	cmd.Flags().StringVarP(&Status, "status", "", "in_progress", "The status to set on the check run")

	_ = cobra.MarkFlagRequired(cmd.Flags(), "path")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "host")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "owner")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "repo")

	_ = cobra.MarkFlagRequired(cmd.Flags(), "commit")
	_ = cobra.MarkFlagRequired(cmd.Flags(), "name")

	return cmd
}

func UpdateCheckRun(credentials string, kind string, host string, owner string, repo string) (string, error) {
	repositoryURL := GetURL(kind, host, owner, repo)

	scmClient, _, _, err := GetScmClient(repositoryURL)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create scm client")
	}

	return CreateOrUpdateCheckRun(scmClient, fmt.Sprintf("%s/%s", owner, repo), Commit)
}

func CreateOrUpdateCheckRun(client *scm.Client, repository string, commit string) (string, error) {
	ctx := context.Background()

	status, _, err := client.Commits.UpdateCommitStatus(ctx, repository, commit, &scm.CommitStatusUpdateOptions{
		Name:  Name,
		Sha:   Commit,
		State: Status,
		//ID:          "",
		//Ref:         "",
		//TargetURL:   "",
		//Description: "",
		//Coverage:    0,
		//PipelineID:  nil,
	})
	if err != nil {
		return "", err
	}

	return status.Ref, nil
}
