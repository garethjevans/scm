package cmd

import (
	"fmt"
	"github.com/garethjevans/scm/pkg/client"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/pkg/errors"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/spf13/cobra"
)

// $ scm pr --host=https://github.com --owner=garethjevans --repo=my-repo
//https://github.com/garethjevans/my-repo
//
//$ scm url --host=https://dev.azure.com --owner=garethjevans --repo=my-repo
//https://dev.azure.com/garethjevans/_git/my-repo

var (
	CommitBranch string
	BaseBranch   string
	PrTitle      string
	GitUser      string
	GitEmail     string
)

// NewPrCreateCmd creates a pr_create command
func NewPrCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "create",
		Short:        "Creates a PR/MR from the local changes",
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
	cmd.Flags().StringVarP(&PrTitle, "pr-title", "", "", "The title of the PR/MR to create")
	cmd.Flags().StringVarP(&GitUser, "git-user", "", "", "The author of any git commits")
	cmd.Flags().StringVarP(&GitEmail, "git-email", "", "", "The author of any git commits")

	_ = cmd.MarkFlagRequired("path")
	_ = cmd.MarkFlagRequired("commit-branch")
	_ = cmd.MarkFlagRequired("pr-title")
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

	fmt.Printf("[DEBUG] %s\n", repositoryURL)

	scmClient, err := GetScmClient(repositoryURL)
	if err != nil {
		return errors.Wrapf(err, "failed to create scm client")
	}

	fmt.Printf("[DEBUG] Got scm client %+v\n", scmClient)

	//ctx := context.Background()

	workTree, err := repository.Worktree()
	if err != nil {
		return errors.Wrapf(err, "unable to access worktree")
	}

	status, err := workTree.Status()
	if err != nil {
		return errors.Wrapf(err, "unable to access status")
	}

	fmt.Printf("[DEBUG] isClean %t\n%+v", status.IsClean(), status)

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
		err = workTree.AddGlob("**")
		if err != nil {
			return errors.Wrapf(err, "unable to add files")
		}

		hash, err := workTree.Commit(PrTitle, &git.CommitOptions{
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
	}

	//pullRequestInput := &scm.PullRequestInput{
	//	Title: o.Title,
	//	Body:  o.Body,
	//	Head:  o.Head,
	//	Base:  o.Base,
	//}

	// shouldUpdate, existingPullRequestNumber := updateNecessary(ctx, o.Head, o.Base, o.AllowUpdate, scmClient, fullName)

	//if shouldUpdate {
	//	res, _, err := scmClient.PullRequests.Update(ctx, fullName, existingPullRequestNumber, pullRequestInput)
	//	if err != nil {
	//		return errors.Wrapf(err, "failed to update existing pull request #%d in repo '%s'", existingPullRequestNumber, fullName)
	//	}
	//
	//	log.Logger().Infof("updated pull request #%d in repo '%s'. url: %s", res.Number, res.Base.Repo.FullName, res.Link)
	//
	//	return nil
	//}

	//res, _, err := scmClient.PullRequests.Create(ctx, fullName, pullRequestInput)
	//if err != nil {
	//	return errors.Wrapf(err, "failed to create a pull request in the repository '%s' with the title '%s'", fullName, o.Title)
	//}
	//
	//log.Logger().Infof("created pull request #%d in repo '%s'. url: %s", res.Number, res.Base.Repo.FullName, res.Link)

	return nil
}

func GetScmClient(repoUrl string) (*scm.Client, error) {
	b, err := os.ReadFile(Path)
	if err != nil {
		return nil, err
	}

	scmClient, err := client.FromRepoURL(repoUrl, string(b))
	return scmClient, nil
}

func DeterminePr(credentials string, kind string, host string, owner string, repo string) (string, error) {
	lines := strings.Split(credentials, "\n")
	for _, line := range lines {
		u, err := url.Parse(strings.TrimSpace(line))
		if err != nil {
			return "", err
		}

		if host != "" {
			h, err := url.Parse(GetURL(kind, host, owner, repo))
			if err != nil {
				return "", err
			}

			//if u.Path != "" {
			//	fmt.Println("[DEBUG] we have a path: " + u.Path)
			//}

			if h.Host == u.Host {
				// we have found a host that matches
				password, ok := u.User.Password()

				if ok {
					return password, nil
				}
			}
		} else {
			// get the first in the list if no host is specified
			password, ok := u.User.Password()

			if ok {
				return password, nil
			}
		}
	}

	return "", fmt.Errorf("unable to locate a token for %s", host)
}

//func updateNecessary(ctx context.Context, head string, base string, updateAllowed bool, scmClient *scm.Client, fullName string) (bool, int) {
//	if !updateAllowed {
//		return false, 0
//	}
//
//	return FindOpenPullRequestByBranches(ctx, head, base, scmClient, fullName)
//}

//func FindOpenPullRequestByBranches(ctx context.Context, head string, base string, scmClient *scm.Client, fullName string) (bool, int) {
//	var openPullRequests []*scm.PullRequest
//	page := 1
//
//	for {
//		pullRequestListOptions := scm.PullRequestListOptions{Page: page, Size: 10, Open: true, Closed: false}
//
//		foundOpenPullRequests, _, err := scmClient.PullRequests.List(ctx, fullName, &pullRequestListOptions)
//		if err != nil {
//			log.Logger().Errorf("listing pull requests in repo '%s' failed: %s", fullName, err)
//			return false, 0
//		}
//
//		if len(foundOpenPullRequests) == 0 {
//			break
//		}
//
//		openPullRequests = append(openPullRequests, foundOpenPullRequests...)
//
//		page++
//	}
//
//	for _, openPullRequest := range openPullRequests {
//		if openPullRequest.Head.Ref == head && openPullRequest.Base.Ref == base {
//			return true, openPullRequest.Number
//		}
//	}
//
//	return false, 0
//}
