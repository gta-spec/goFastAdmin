//go:build debug

package pkg

import "github.com/gin-gonic/gin"

func init() {
	EnvGinMode = gin.DebugMode
}
