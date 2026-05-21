package pkg

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var (
	EnvGinMode = gin.DebugMode
	rootPath   string
)

const (
	DS           = string(filepath.Separator)
	APP_PATH     = "internal/"
	CONF_PATH    = "pkg/config/"
	INSTALL_PATH = "internal/admin/command/install/"
)

func init() {
	rootPath, _ = os.Getwd()
}

func RootPath() string {
	return rootPath
}
