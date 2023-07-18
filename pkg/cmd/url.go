package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

// $ scm url --host=https://github.com --owner=garethjevans --repo=my-repo
//https://github.com/garethjevans/my-repo
//
//$ scm url --host=https://dev.azure.com --owner=garethjevans --repo=my-repo
//https://dev.azure.com/garethjevans/_git/my-repo

var (
	Host  string
	Owner string
	Repo  string
	Kind  string
)

// NewUrlCmd creates a new cluster command.
func NewUrlCmd() *cobra.Command {
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
		Args: cobra.NoArgs,
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
