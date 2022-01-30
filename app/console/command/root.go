package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tvb-sz/serve-swagger-ui/app/console"
	"github.com/tvb-sz/serve-swagger-ui/client/initializer"
	"github.com/tvb-sz/serve-swagger-ui/conf"
	"os"
)

// RootCmd 基于cobra的命令行根节点定义
var (
	RootCmd = &cobra.Command{
		Use:   "serve-swagger-ui --path [--host --port --app-key= --app-secret=]",
		Short: "serve-swagger-ui service manage",
		Long: `serve-swagger-ui command example
----------------------------------------------------------------------
serve-swagger-ui --path=--Your-swagger-file-PATH-- [--host --port --app-key= --app-secret]
serve-swagger-ui --config=--Your-TOML-config-file-PATH--
----------------------------------------------------------------------`,
		Run: func(cmd *cobra.Command, args []string) {
			console.BootStrap()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// step1、init config
			conf.Init()

			// step2、init global client handle
			initializer.Init()

			return nil
		},
	}
)

func init() {
	// command line args
	RootCmd.PersistentFlags().StringVar(&conf.Cmd.ConfigFile, "config", "", "指定本地配置文件")
	RootCmd.PersistentFlags().StringVar(&conf.Cmd.Host, "host", "", "指定配置文件类型")
	RootCmd.PersistentFlags().IntVar(&conf.Cmd.Port, "port", 0, "指定优先级高于配置文件的日志存储位置：stderr|stdout|目录路径")
	RootCmd.PersistentFlags().StringVar(&conf.Cmd.SwaggerPath, "path", "", "指定优先级高于配置文件的日志存储位置：stderr|stdout|目录路径")
	RootCmd.PersistentFlags().StringVar(&conf.Cmd.GoogleClientID, "client_id", "", "指定优先级高于配置文件的日志级别：debug|info|warn|error|panic|fatal")
	RootCmd.PersistentFlags().StringVar(&conf.Cmd.GoogleClientSecret, "client_secret", "", "指定优先级高于配置文件的日志级别：debug|info|warn|error|panic|fatal")
	RootCmd.PersistentFlags().StringVar(&conf.Cmd.LogLevel, "log_level", "", "指定优先级高于配置文件的日志级别：debug|info|warn|error|panic|fatal")
	RootCmd.PersistentFlags().StringVar(&conf.Cmd.LogPath, "log_path", "", "指定优先级高于配置文件的日志级别：debug|info|warn|error|panic|fatal")
}

// Start 启动应用
func Start() {
	err := RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
