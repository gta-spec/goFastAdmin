//go:generate go run src/cmd/tools/gen_route.go
//go:generate swag init -o ./src/docs

package main

import (
	"fmt"
	_ "gota/internal/api/controller"
	_ "gota/internal/index/controller"
	"gota/src"
	"gota/src/app"
	_ "gota/src/app/autoload"
	"gota/src/cmd/tools"
	"gota/src/config"
	"gota/src/database"
	"gota/src/database/redis"
	_ "gota/src/i18n"
	"gota/src/logger"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

// @title OpenAPI
// @version 1.0
// @description Gota框架
// @termsOfService https://www.swagger.io/terms/
// @host 127.0.0.1:8080
// @BasePath /
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://github.com/swaggo/swag/blob/master/README_zh-CN.md
// @SecurityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description API 认证方式：在请求头中添加 Authorization 字段，值为 Bearer + 空格 + token
func main() {
	root := src.AppPath
	//mainPath := pkg.MainPath
	//outFilename := filepath.Join("./", root, "route_gen.go")
	//out, err := os.Create(outFilename)
	//if err != nil {
	//	panic(err)
	//}
	//defer out.Close()

	fs := os.DirFS(filepath.Join("./", root))

	filenames, err := doublestar.Glob(fs, "*/controller/**.go")
	//fmt.Println(mainPath)
	if err != nil {
		panic(err)
	}

	for _, filename := range filenames {
		if filename != "api/controller/Demo.go" {
			continue
		}
		file, err := fs.Open(filename)
		if err != nil {
			fmt.Printf("打开文件失败 %s: %v\n", filename, err)
			continue
		}

		ast, err := tools.NewFileAst(filename, file)
		if err != nil {
			fmt.Println(err)
		}
		ast.ParseFile()
	}
	//out.WriteString(`package internal`)
	//if mainPath != "" {
	//	return
	//}

	defer logger.Close()
	defer database.Close()
	defer redis.Close()
	// 加载配置
	_ = config.SetGlobalConfigFile("src/config/config.yaml")
	cfgInst := config.Viper()
	// 加载次要配置
	config.LoadConfigGlob(strings.Join([]string{src.AppPath, "extra", "*.yaml"}, "/"))

	application := app.New(cfgInst)

	application.Run()
}
