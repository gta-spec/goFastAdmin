package internal

import "github.com/gin-gonic/gin"

func init() {
	RegisterRoutes = func() {
		if index, ok := Modules["index"]; ok {
			index.Group("/c", func(c *gin.Context) {

			})
		}
	}
}
