package main

import (
	"app/cmd/migrate"
	"app/cmd/server"
	"app/internal/web"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
)

func main() {
	log.SetFlags(log.Default().Flags() | log.Llongfile)

	var rootCmd = &cobra.Command{
		Use: path.Base(os.Args[0]),
		Run: func(cmd *cobra.Command, args []string) {
			// 默认执行 server 命令
			serverCmd, _, _ := cmd.Find([]string{"server"})
			if serverCmd != nil {
				serverCmd.Run(cmd, args)
			}
		},
	}
	rootCmd.PersistentFlags().StringP("config", "c", "./config.yaml", "apiserver config file path.")

	server.Register(rootCmd, web.BuildFS, web.IndexPage)
	migrate.Register(rootCmd)
	rootCmd.Execute()
}
