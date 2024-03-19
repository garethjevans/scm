package cmd

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

var (
	Sha        string
	Ref        string
	Repository string
)

// NewTagCmd creates a pr_create command.
func NewTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "tag",
		Short:        "Creates a tag using the SCM api, will determine the repository url from the current directory, or can be overridden with --repository.",
		Long:         "",
		Example:      "scm tag --tag 0.0.1 --sha abcdefgabcdefgabcdefgabcdefg --path .git-credentials",
		Aliases:      []string{"t"},
		RunE:         Tag,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(&Path, "path", "p", "", "The path to the git-credentials file")
	cmd.Flags().StringVarP(&Sha, "sha", "s", "", "The sha to attach the tag to")
	cmd.Flags().StringVarP(&Ref, "tag", "t", "", "The tag to apply")
	cmd.Flags().StringVarP(&Repository, "repository", "r", "", "The full https url of the repository")

	_ = cmd.MarkFlagRequired("path")
	_ = cmd.MarkFlagRequired("sha")
	_ = cmd.MarkFlagRequired("tag")

	return cmd
}

func Tag(cmd *cobra.Command, args []string) error {
	var repositoryURL string
	if Repository == "" {
		repository, err := git.PlainOpen(".")
		if err != nil {
			return errors.Wrapf(err, "unable to open repository")
		}

		origin, err := repository.Remote("origin")
		if err != nil {
			return errors.Wrapf(err, "unable to get remote 'origin'")
		}

		repositoryURL = origin.Config().URLs[0]
	} else {
		repositoryURL = Repository
	}

	fmt.Printf("[DEBUG] Determined repository remote URL as: %s\n", repositoryURL)

	scmClient, _, _, err := GetScmClient(repositoryURL)
	if err != nil {
		return errors.Wrapf(err, "failed to create scm client")
	}

	fmt.Printf("[DEBUG] Got an scm client talking to: %s\n", scmClient.BaseURL)

	// if a PR already exists for this branch we should skip as this is being updated
	r, err := url.Parse(repositoryURL)
	if err != nil {
		return err
	}

	fullName := strings.TrimPrefix(strings.TrimSuffix(r.Path, ".git"), "/")

	ctx := context.Background()
	_, _, err = scmClient.Git.CreateRef(ctx, fullName, fmt.Sprintf("refs/tags/%s", Ref), Sha)
	if err != nil {
		return err
	}

	return nil
}
