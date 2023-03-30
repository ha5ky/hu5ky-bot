/**
 * @Author Nil
 * @Description cmd/root.go
 * @Date 2023/3/27 20:02
 **/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func Execute() {
	rootCmd := GetRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "botd",
		Short: "http server of hu5ky bot",
		Long:  "http server of hu5ky bot, a ChatGPT based server, we can use ours user case to train the bot",
	}
	rootCmd.AddCommand(versionCmd, serverCmd)
	return rootCmd
}
