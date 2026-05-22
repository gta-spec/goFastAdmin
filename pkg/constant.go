package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	supportMinGoVer = uint64(26)
	DS              = string(filepath.Separator)
	AppPath         = "internal/"
	ConfPath        = "pkg/config/"
	InstallPath     = "internal/admin/command/install/"
)

var (
	EnvGinMode = gin.DebugMode
	MainPath   string
	GoVersion  string
	rootPath   string
)

func init() {
	rootPath, _ = os.Getwd()
	info, ok := debug.ReadBuildInfo()
	if ok {
		MainPath = info.Main.Path
		GoVersion = info.GoVersion
		if v, e := getMinVer(GoVersion); e == nil && v < supportMinGoVer {
			log.Fatal(fmt.Sprintf("[ERROR] Now Gin requires Go 1.%d+.", supportMinGoVer))
		}
	}
}

func RootPath() string {
	return rootPath
}

func getMinVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseUint(v[first+1:], 10, 64)
	}
	return strconv.ParseUint(v[first+1:last], 10, 64)
}
