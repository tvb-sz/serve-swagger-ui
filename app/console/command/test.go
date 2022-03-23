package command

import (
	"github.com/spf13/cobra"
)

// init test sub-command
func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "test",
		Short: "test command",
		Long:  "test command for dev",
		Run: func(cmd *cobra.Command, args []string) {

		},
	})
}
