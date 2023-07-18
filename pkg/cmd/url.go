package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	Kind string
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
			if isAzureDevops(Kind, Host) {
				fmt.Fprintf(cmd.OutOrStdout(), "%s/%s/_git/%s", strings.TrimSuffix(Host, "/"), Owner, Repo)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "%s/%s/%s", strings.TrimSuffix(Host, "/"), Owner, Repo)
			}
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

func isAzureDevops(kind string, host string) bool {
	if kind == "azure" {
		return true
	}

	if strings.TrimSuffix(host, "/") == "https://dev.azure.com" {
		return true
	}

	return false
}
