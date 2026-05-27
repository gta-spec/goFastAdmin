package controller

import (
	"gota/internal/common/controller"
	"gota/pkg/app/route"

	"github.com/gin-gonic/gin"
)

func init() {
	route.Register(&Index{
		NoNeedLogin: []string{"*"},
		NoNeedRight: []string{"*"},
	})
}

type Index struct {
	controller.Frontend
	NoNeedLogin []string
	NoNeedRight []string
}

func (t *Index) Index(c *gin.Context) {
	//t.View.Fetch(c)
}
