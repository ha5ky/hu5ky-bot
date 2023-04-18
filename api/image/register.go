/**
 * @Author Nil
 * @Description api/image/register.go
 * @Date 2023/4/10 20:43
 **/

package image

import "github.com/ha5ky/hu5ky-bot/router/schema"

func init() {
	scheme := schema.NewSchemeBuilder().Register()

	scheme.GET("/image", Generate)
}
