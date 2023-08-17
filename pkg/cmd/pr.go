package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// NewPrCmd creates a new token command.
func NewPrCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pr",
		Short:   "",
		Long:    "",
		Example: "",
		Aliases: []string{"p"},
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := os.ReadFile(Path)
			if err != nil {
				return err
			}
			token, err := DeterminePr(string(b), Kind, Host, Owner, Repo)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), token)
			return nil
		},
		Args:         cobra.NoArgs,
		SilenceUsage: true,
	}

	cmd.AddCommand(NewPrCreateCmd())

	return cmd
}
