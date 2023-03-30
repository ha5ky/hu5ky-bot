/**
 * @Author Nil
 * @Description router/middleware/cors/cors.go
 * @Date 2023/3/28 01:02
 **/

package cors

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Headers", "*, Origin, X-Requested-With, Content-ModeType, Accept")
		ctx.Header("content-type", "text/plain")
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
		if ctx.Request.Method == http.MethodOptions {
			ctx.JSON(http.StatusOK, "Options Request!")
		} else {
			ctx.Next()
		}
	}
}
