package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"strings"
)

// $ scm url --host=https://github.com --owner=garethjevans --repo=my-repo
//https://github.com/garethjevans/my-repo
//
//$ scm url --host=https://dev.azure.com --owner=garethjevans --repo=my-repo
//https://dev.azure.com/garethjevans/_git/my-repo

var (
	Path string
)

// NewTokenCmd creates a new token command.
func NewTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "token",
		Short:   "Determines the token to use for an scm provider",
		Long:    "",
		Example: "scm token --host=https://github.com --path .git-credentials",
		Aliases: []string{"t"},
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := os.ReadFile(Path)
			if err != nil {
				return err
			}
			token, err := DetermineToken(string(b), Host)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), token)
			return nil
		},
		Args: cobra.NoArgs,
	}

	cmd.Flags().StringVarP(&Host, "host", "", "", "The host of the scm provider, including scheme")
	cmd.Flags().StringVarP(&Path, "path", "p", "", "The path to the git-credentials file")

	return cmd
}

func DetermineToken(credentials string, host string) (string, error) {
	lines := strings.Split(credentials, "\n")
	for _, line := range lines {
		u, err := url.Parse(strings.TrimSpace(line))
		if err != nil {
			return "", err
		}

		if host != "" {
			h, err := url.Parse(strings.TrimSpace(host))
			if err != nil {
				return "", err
			}

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

	return "", nil
}
