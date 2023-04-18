/**
 * @Author Nil
 * @Description api/chat/register.go
 * @Date 2023/4/10 19:23
 **/

package chat

import (
	"github.com/ha5ky/hu5ky-bot/router/schema"
)

func init() {
	scheme := schema.NewSchemeBuilder().Register()

	scheme.GET("/chat", Get)
}
