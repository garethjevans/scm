package cmd

import (
	"github.com/spf13/cobra"
)

// NewPrCmd creates a new token command.
func NewPrCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "pr",
		Short:        "",
		Long:         "",
		Example:      "",
		Aliases:      []string{"p"},
		SilenceUsage: true,
	}

	cmd.AddCommand(NewPrCreateCmd())

	return cmd
}
