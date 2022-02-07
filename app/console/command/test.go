package command

import (
	"github.com/spf13/cobra"
)

// init version子命令
func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "test",                 // 子命令名称
		Short: "test command",         // 子命令简短说明
		Long:  "test command for dev", // 子命令完整说明
		Run: func(cmd *cobra.Command, args []string) {

		},
	})
}
