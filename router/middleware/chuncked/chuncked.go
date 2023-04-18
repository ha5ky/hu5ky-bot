/**
 * @Author Nil
 * @Description router/middleware/chuncked.go
 * @Date 2023/4/16 17:29
 **/

package chuncked

import (
	"github.com/gin-gonic/gin"
)

func Chunked() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Transfer-Encoding", "chunked")
		ctx.Next()
	}
}
