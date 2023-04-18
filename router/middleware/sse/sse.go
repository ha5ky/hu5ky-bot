/**
 * @Author Nil
 * @Description router/middleware/sse/sse.go
 * @Date 2023/4/16 22:54
 **/

package sse

import "github.com/gin-gonic/gin"

func EventStream() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "text/event-stream")
		ctx.Writer.Header().Set("Cache-Control", "no-cache")
		ctx.Writer.Header().Set("Connection", "keep-alive")
		ctx.Next()
	}
}
