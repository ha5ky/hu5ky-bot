/**
 * @Author Nil
 * @Description pkg/config/version.go
 * @Date 2023/3/28 01:08
 **/

package config

var (
	PlatformVersion = "v0.1.1" // Version platform version

	Project           = "/api" // default /api
	APIVersionV1      = "/v1"
	APICurrentVersion = APIVersionV1 // version of phdalios openapi
	GitCommit         = "default"    // get commit from git
)
