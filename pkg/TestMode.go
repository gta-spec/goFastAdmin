//go:build test

package pkg

import "github.com/gin-gonic/gin"

func init() {
	EnvGinMode = gin.TestMode
}
