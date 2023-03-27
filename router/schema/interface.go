/**
 * @Author Nil
 * @Description router/schema/interface.go
 * @Date 2023/3/28 01:01
 **/

package schema

import (
	"github.com/ka5ky/hu5ky-bot/pkg/config"
	"sync"

	"github.com/ka5ky/hu5ky-bot/router/middleware/cors"

	"github.com/gin-gonic/gin"
)

var (
	r        *gin.Engine
	apiGroup *gin.RouterGroup
	once     sync.Once
)

func init() {
	r = gin.Default()
	r.Use(cors.Cors())
	apiGroup = r.Group(config.Project)
}

type SchemeBuilder struct {
	Version    string
	Middleware []gin.HandlerFunc
	Scheme     *gin.RouterGroup
}

func NewSchemeBuilder() *SchemeBuilder {
	return &SchemeBuilder{
		Version: config.APIVersionV1,
		Scheme:  apiGroup,
		// todo could be controlled by plugins (config.yaml)
		Middleware: []gin.HandlerFunc{
			//authorization.Authorization(),
		},
	}
}

func (sb *SchemeBuilder) Register() *gin.RouterGroup {
	return sb.Scheme.Group(sb.Version, sb.Middleware...)
}

func Registry() *gin.Engine {
	gin.SetMode(config.PhdaliosCache.ServerConfig.Mode)
	return r
}
