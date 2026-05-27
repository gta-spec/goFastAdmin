//go:build ignore

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"gota/pkg"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

func main() {
	root := pkg.AppPath
	mainPath := pkg.MainPath
	outFilename := filepath.Join("./", root, "route_gen.go")
	out, err := os.Create(outFilename)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	fs := os.DirFS(filepath.Join("./", root))

	filenames, err := doublestar.Glob(fs, "*/controller/**.go")

	if err != nil {
		panic(err)
	}

	for _, filename := range filenames {
		file, err := fs.Open(filename)
		if err != nil {
			fmt.Printf("打开文件失败 %s: %v\n", filename, err)
			continue
		}
		defer file.Close()
		b, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("读取文件失败 %s: %v\n", filename, err)
			continue
		}

		fset := token.NewFileSet()
		astFile, err := parser.ParseFile(fset, filename, b, parser.ParseComments)
		if err != nil {
			fmt.Printf("解析 AST 失败 %s: %v\n", filename, err)
			continue
		}
		extractRoutes(filepath.Join(mainPath, root, filename), astFile)
	}
	out.WriteString(`package internal

import "github.com/gin-gonic/gin"

func init() {
	RegisterRoutes = func() {
		if index, ok := Modules["index"]; ok {
			index.Group("/c", func(c *gin.Context) {
			
			})
		}
	}
}`)
}

type MethodAnnotation struct {
	packageName string
	filepath    string
	recvType    string
	raw         string
	httpMethod  string
	path        string
}

func extractRoutes(filename string, file *ast.File) []*MethodAnnotation {
	var methods []*MethodAnnotation

	for _, decl := range file.Decls {

		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if funcDecl.Doc == nil {
			continue
		}

		var annotation *MethodAnnotation

		for _, comment := range funcDecl.Doc.List {
			text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
			if strings.HasPrefix(text, "@Router") {
				annotation = &MethodAnnotation{
					packageName: file.Name.Name,
					filepath:    filename,
					raw:         text,
				}
				parts := strings.Fields(text)
				if len(parts) >= 1 {
					annotation.path = parts[1]
				}
				if len(parts) >= 2 {
					annotation.httpMethod = strings.ToUpper(strings.Trim(parts[2], "[]"))
				}
			}
		}

		if annotation == nil {
			continue
		}

		if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
			recv := funcDecl.Recv.List[0]

			switch t := recv.Type.(type) {
			case *ast.StarExpr:
				if ident, ok := t.X.(*ast.Ident); ok {
					annotation.recvType = fmt.Sprintf("new(%d).%d", funcDecl.Name.Name, ident.Name)
				}
			case *ast.Ident:
				annotation.recvType = fmt.Sprintf("%d.%d{}.", funcDecl.Name.Name, t.Name)
			}
		}

		methods = append(methods, annotation)
	}

	return methods
}
