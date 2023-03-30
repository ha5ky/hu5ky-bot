/**
 * @Author Nil
 * @Description api/completion/register.go
 * @Date 2023/3/30 10:07
 **/

package completion

import "github.com/ha5ky/hu5ky-bot/router/schema"

func init() {
	scheme := schema.NewSchemeBuilder().Register()

	scheme.GET("/completion", Get)
}
