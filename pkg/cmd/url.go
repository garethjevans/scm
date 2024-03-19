package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// NewURLCmd creates a new cluster command.
func NewURLCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "url",
		Short:   "Calculates the url for an scm provider",
		Long:    "",
		Example: "scm url --host=https://github.com --owner=garethjevans --repo=my-repo",
		Aliases: []string{"u"},
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(cmd.OutOrStdout(), GetURL(Kind, Host, Owner, Repo))
			return nil
		},
		Args:         cobra.NoArgs,
		SilenceUsage: true,
	}

	cmd.Flags().StringVarP(&Host, "host", "", "", "The host of the scm provider, including scheme")
	cmd.Flags().StringVarP(&Owner, "owner", "o", "", "The owner of the repository")
	cmd.Flags().StringVarP(&Repo, "repo", "r", "", "The name of the repository")
	cmd.Flags().StringVarP(&Kind, "kind", "k", "", "The kind of the scm provider")

	return cmd
}

func GetURL(kind string, host string, owner string, repo string) string {
	if isAzureDevops(kind, host) {
		return fmt.Sprintf("%s/%s/_git/%s", strings.TrimSuffix(host, "/"), owner, repo)
	}

	return fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(host, "/"), owner, repo)
}

func isAzureDevops(kind string, host string) bool {
	if kind == "azure" {
		return true
	}

	if strings.TrimSuffix(host, "/") == "https://dev.azure.com" {
		return true
	}

	return false
}
