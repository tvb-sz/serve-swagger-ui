package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tvb-sz/serve-swagger-ui/conf"
)

// init version子命令
func init() {
	RootCmd.AddCommand(&cobra.Command{
		Use:   "version",           // 子命令名称
		Short: "show version info", // 子命令简短说明
		Long:  "show version info", // 子命令完整说明
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(conf.Config.Server.Version)
		},
	})
}
