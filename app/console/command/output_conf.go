package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tvb-sz/serve-swagger-ui/stubs"
)

// init version子命令
func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "output_conf",              // 子命令名称
		Short: "output all config values", // 子命令简短说明
		Long:  "output all config values", // 子命令完整说明
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("# Copy the following output to create a new Configuration file for .toml suffix")
			fmt.Println(stubs.ConfExample)
		},
	})
}
