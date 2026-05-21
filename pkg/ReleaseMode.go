//go:build release

package pkg

import "github.com/gin-gonic/gin"

func init() {
	EnvGinMode = gin.ReleaseMode
}
