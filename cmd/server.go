/**
 * @Author Nil
 * @Description cmd/server.go
 * @Date 2023/3/28 01:00
 **/

package cmd

import "github.com/spf13/cobra"

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "run http server",
		Run: func(cmd *cobra.Command, args []string) {
			configLogs := logger.NewProductionLog()
			configLogs.Level = logLevel
			logger.InitLog(configLogs)

			config.Watcher(mode)
			svrCfg := &config.PhdaliosCache.ServerConfig
			svrCfg.Mode = mode
			svrCfg.LogLevel = logLevel
			svrCfg.Pattern = pattern
			if svrCfg.Pattern == config.InCluster {
				kubeconfig = ""
			}
			svrCfg.Kubeconfig = kubeconfig

			httpCfg := &config.PhdaliosCache.HttpConfig
			httpCfg.Port = port
			httpCfg.ReadTimeout = readTimeout
			httpCfg.WriteTimeout = writeTimeout
			httpCfg.MaxHeaderBytes = maxHeaderBytes

			//logger.Info(util.ConvStr("phdalios cache: ", *config.PhdaliosCache))
			//if err := config.Cache2yaml(config.PhdaliosCache); err != nil {
			//	logger.Error(util.ConvStr("writing yaml fail.", err))
			//}

			// init kube-client
			kube.InitClient(svrCfg.Mode)

			//init mongo
			db := model.GetDB()
			defer func() {
				if err := model.Disconnect(db); err != nil {
					panic(err)
				}
			}()
			queue.InitQueue()
			queue.InitDeploymentQueue()

			ns := namespace.NewNamespace(kube.Kc)
			ns.Namespace.ObjectMeta.Name = config.PhdaliosABTestingNamespace
			_ = ns.Create(&metav1.CreateOptions{})

			ns = namespace.NewNamespace(kube.Kc)
			ns.Namespace.ObjectMeta.Name = config.PhdaliosCanaryNamespace
			_ = ns.Create(&metav1.CreateOptions{})

			router.Registry()
			router.ListenAndRun()
		},
	}
)
