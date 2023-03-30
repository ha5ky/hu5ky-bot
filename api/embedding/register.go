/**
 * @Author Nil
 * @Description api/embedding/register.go
 * @Date 2023/3/28 21:12
 **/

package embedding

import "github.com/ha5ky/hu5ky-bot/router/schema"

func init() {
	scheme := schema.NewSchemeBuilder().Register()

	scheme.POST("/embedding", Create)
	scheme.GET("/embeddings", List)
}
