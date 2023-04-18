/**
 * @Author Nil
 * @Description api/file/register.go
 * @Date 2023/4/2 15:47
 **/

package file

import "github.com/ha5ky/hu5ky-bot/router/schema"

func init() {
	scheme := schema.NewSchemeBuilder().Register()

	scheme.POST("/file", Upload)
}
