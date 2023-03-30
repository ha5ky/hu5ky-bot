/**
 * @Author Nil
 * @Description cmd/server.go
 * @Date 2023/3/28 01:00
 **/

package cmd

import (
	"github.com/ha5ky/hu5ky-bot/model"
	"github.com/ha5ky/hu5ky-bot/pkg/config"
	"github.com/ha5ky/hu5ky-bot/router"
	"github.com/spf13/cobra"
	"gorm.io/gorm/schema"
)

var (
	tables = []schema.Tabler{
		&model.File{},
	}
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "run http server",
		Run: func(cmd *cobra.Command, args []string) {
			config.Watcher(mode)
			svrCfg := &config.SysCache.ServerConfig
			svrCfg.Mode = mode
			svrCfg.LogLevel = logLevel
			model.Registry()
			router.Registry()
			router.ListenAndRun()
		},
	}
)

func init() {
	serverCmd.Flags().StringVarP(
		&mode,
		"mode",
		"m",
		"debug",
		"set server mode, including debug , release , test",
	)
	if err := serverCmd.MarkFlagRequired("mode"); err != nil {
		panic(err)
	}

	serverCmd.Flags().StringVarP(
		&logLevel,
		"loglevel",
		"l",
		"debug",
		"including:\n debug\n info\n warn\n error\n dpanic\n panic\n fatal",
	)
}
