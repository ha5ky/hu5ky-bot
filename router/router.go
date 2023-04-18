/**
 * @Author Nil
 * @Description router/router.go
 * @Date 2023/3/28 13:38
 **/

package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ha5ky/hu5ky-bot/pkg/config"
	"github.com/ha5ky/hu5ky-bot/pkg/logger"
	"github.com/ha5ky/hu5ky-bot/router/schema"
	"net/http"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "github.com/ha5ky/hu5ky-bot/api/chat"
	_ "github.com/ha5ky/hu5ky-bot/api/collection"
	_ "github.com/ha5ky/hu5ky-bot/api/completion"
	_ "github.com/ha5ky/hu5ky-bot/api/embedding"
	_ "github.com/ha5ky/hu5ky-bot/api/file"
	_ "github.com/ha5ky/hu5ky-bot/api/image"
)

var (
	r   *gin.Engine
	svr *http.Server
)

func Registry() {
	r = schema.Registry()
}

func ListenAndRun() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	svr = &http.Server{
		Addr:           ":" + strconv.Itoa(int(config.SysCache.HttpConfig.Port)),
		Handler:        r,
		ReadTimeout:    time.Duration(config.SysCache.HttpConfig.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.SysCache.HttpConfig.WriteTimeout) * time.Second,
		MaxHeaderBytes: config.SysCache.HttpConfig.MaxHeaderBytes << 20,
	}

	go func() {
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	logger.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svr.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", err.Error())
	}
}
