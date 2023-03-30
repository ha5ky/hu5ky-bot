/**
 * @Author Nil
 * @Description cmd/version.go
 * @Date 2023/3/28 16:50
 **/

package cmd

import (
	"fmt"
	"github.com/ha5ky/hu5ky-bot/pkg/config"
	"github.com/spf13/cobra"
	"runtime"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "get version of phdaliosd and openapi",
		Long:  "check the information such as: phdaliosd, openapi, os, go, gitCommit.",
		Run: func(cmd *cobra.Command, args []string) {
			phdaliosdVersion := config.PlatformVersion
			goVersion := runtime.Version()
			gitCommit := config.GitCommit
			apiVersion := config.APICurrentVersion
			osAndArch := fmt.Sprintf("os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
			fmt.Printf(
				"phdaliosd version: %s\ngo version: %s\ngit commit: %s\napi version: %s\n%s\n",
				phdaliosdVersion,
				goVersion,
				gitCommit,
				apiVersion,
				osAndArch,
			)
		},
	}
)
