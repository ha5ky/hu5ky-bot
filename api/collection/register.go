/**
 * @Author Nil
 * @Description api/collection/register.go
 * @Date 2023/4/2 18:06
 **/

package collection

import "github.com/ha5ky/hu5ky-bot/router/schema"

func init() {
	scheme := schema.NewSchemeBuilder().Register()

	scheme.POST("/collection", Create)
	scheme.GET("/collections", List)
}
