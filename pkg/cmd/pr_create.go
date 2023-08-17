package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/garethjevans/scm/pkg/client"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/pkg/errors"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/spf13/cobra"
)

var (
	CommitBranch string
	BaseBranch   string
	PrTitle      string
	CommitTitle  string
	GitUser      string
	GitEmail     string
)

// NewPrCreateCmd creates a pr_create command.
func NewPrCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create",
		Short:        "Commits and pushes any local changes, if the base & commit branches are different, then a PR/MR will be created.",
		Long:         "",
		Example:      "scm pr create --commit-branch update .git-credentials",
		Aliases:      []string{"c"},
		RunE:         CreatePullRequest,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(&Path, "path", "p", "", "The path to the git-credentials file")
	cmd.Flags().StringVarP(&CommitBranch, "commit-branch", "", "", "The branch to push the commits to")
	cmd.Flags().StringVarP(&BaseBranch, "base-branch", "", "main", "The branch to target the PR to")
	cmd.Flags().StringVarP(&CommitTitle, "commit-title", "", "", "The title of any commits to push")
	cmd.Flags().StringVarP(&PrTitle, "pr-title", "", "", "The title of the PR/MR to create")
	cmd.Flags().StringVarP(&GitUser, "git-user", "", "", "The author of any git commits")
	cmd.Flags().StringVarP(&GitEmail, "git-email", "", "", "The author of any git commits")
	cmd.Flags().StringVarP(&Kind, "kind", "", "", "In case we are unable to determine the type of scm server, this can provide hints")

	_ = cmd.MarkFlagRequired("path")
	_ = cmd.MarkFlagRequired("commit-branch")
	_ = cmd.MarkFlagRequired("git-user")
	_ = cmd.MarkFlagRequired("git-email")

	return cmd
}

func CreatePullRequest(cmd *cobra.Command, args []string) error {
	repository, err := git.PlainOpen(".")
	if err != nil {
		return errors.Wrapf(err, "unable to open repository")
	}

	origin, err := repository.Remote("origin")
	if err != nil {
		return errors.Wrapf(err, "unable to get remote 'origin'")
	}

	repositoryURL := origin.Config().URLs[0]

	fmt.Printf("[DEBUG] Determined repository remote URL as: %s\n", repositoryURL)

	scmClient, username, token, err := GetScmClient(repositoryURL)
	if err != nil {
		return errors.Wrapf(err, "failed to create scm client")
	}

	fmt.Printf("[DEBUG] Got an scm client talking to: %s\n", scmClient.BaseURL)

	workTree, err := repository.Worktree()
	if err != nil {
		return errors.Wrapf(err, "unable to access worktree")
	}

	status, err := workTree.Status()
	if err != nil {
		return errors.Wrapf(err, "unable to access status")
	}

	fmt.Printf("[DEBUG] Ensure we are working on branch %s\n", CommitBranch)

	err = workTree.Checkout(&git.CheckoutOptions{
		Create: true,
		Branch: plumbing.NewBranchReferenceName(CommitBranch),
		Keep:   true,
	})
	if err != nil {
		err = workTree.Checkout(&git.CheckoutOptions{
			Create: false,
			Branch: plumbing.NewBranchReferenceName(CommitBranch),
			Keep:   true,
		})
		if err != nil {
			return errors.Wrapf(err, "unable to checkout branch %s", CommitBranch)
		}
	}

	if !status.IsClean() {
		fmt.Printf("[DEBUG] There are local changes that we need to commit\n")
		err = workTree.AddGlob("**")
		if err != nil {
			return errors.Wrapf(err, "unable to add files")
		}

		hash, err := workTree.Commit(CommitTitle, &git.CommitOptions{
			AllowEmptyCommits: false,
			Author: &object.Signature{
				Name:  GitUser,
				Email: GitEmail,
				When:  time.Now(),
			},
		})
		if err != nil {
			return errors.Wrapf(err, "unable to create commit")
		}

		obj, err := repository.CommitObject(hash)
		if err != nil {
			return errors.Wrapf(err, "unable to create repository commit")
		}
		fmt.Printf("[DEBUG] obj %+v\n", obj)

		// push using default options
		err = repository.Push(&git.PushOptions{
			RemoteName: "origin",
			Progress:   os.Stdout,
			RefSpecs: []config.RefSpec{
				config.RefSpec(fmt.Sprintf("+refs/heads/%s:refs/heads/%s", CommitBranch, CommitBranch)),
			},
			Auth: &http.BasicAuth{
				Username: username,
				Password: token,
			},
		})
		if err != nil {
			return errors.Wrapf(err, "unable to push to remote repository")
		}
	}

	// if a PR already exists for this branch we should skip as this is being updated
	r, err := url.Parse(repositoryURL)
	if err != nil {
		return err
	}

	fullName := strings.TrimPrefix(strings.TrimSuffix(r.Path, ".git"), "/")

	if CommitBranch != BaseBranch {
		ctx := context.Background()

		exists, prNumber, prURL := existingPr(ctx, CommitBranch, BaseBranch, scmClient, fullName)

		if exists {
			fmt.Printf("[DEBUG] nothing to do, PR-%d already exists at url %s\n", prNumber, prURL)
		} else {
			pullRequestInput := &scm.PullRequestInput{
				Title: PrTitle,
				Body:  "",
				Head:  CommitBranch,
				Base:  BaseBranch,
			}

			res, _, err := scmClient.PullRequests.Create(ctx, fullName, pullRequestInput)
			if err != nil {
				return errors.Wrapf(err, "failed to create a pull request in the repository '%s' with the title '%s'", r.Path, PrTitle)
			}

			fmt.Printf("[DEBUG] PR-%d created at url %s\n", res.Number, res.Link)
		}
	}
	return nil
}

func GetScmClient(repositoryURL string) (*scm.Client, string, string, error) {
	b, err := os.ReadFile(Path)
	if err != nil {
		return nil, "", "", err
	}

	scmClient, username, token, err := client.FromRepoURL(repositoryURL, string(b), Kind)
	return scmClient, username, token, err
}

func existingPr(ctx context.Context, head string, base string, scmClient *scm.Client, fullName string) (bool, int, string) {
	return FindOpenPullRequestByBranches(ctx, head, base, scmClient, fullName)
}

func FindOpenPullRequestByBranches(ctx context.Context, head string, base string, scmClient *scm.Client, fullName string) (bool, int, string) {
	var openPullRequests []*scm.PullRequest
	page := 1

	for {
		pullRequestListOptions := scm.PullRequestListOptions{Page: page, Size: 10, Open: true, Closed: false}

		foundOpenPullRequests, _, err := scmClient.PullRequests.List(ctx, fullName, &pullRequestListOptions)
		if err != nil {
			fmt.Printf("[WARN] listing pull requests in repo '%s' failed: %s", fullName, err)
			return false, 0, ""
		}

		if len(foundOpenPullRequests) == 0 {
			break
		}

		openPullRequests = append(openPullRequests, foundOpenPullRequests...)

		page++
	}

	for _, openPullRequest := range openPullRequests {
		if openPullRequest.Head.Ref == head && openPullRequest.Base.Ref == base {
			return true, openPullRequest.Number, openPullRequest.Link
		}
	}

	return false, 0, ""
}
