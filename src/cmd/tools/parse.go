package tools

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"regexp"
)

var (
	// @Router /path [method]
	routerPattern = regexp.MustCompile(`@Router\s+(\S+)\s*\[(\w+)]`)
)

type File struct {
	packageName string
	file        *ast.File
	Objects     []*StructType
	Function    []*FuncDecl
}

type StructType struct {
	Name    string
	Docs    string
	Methods []*FuncDecl
}

func NewStruct(nodeType *ast.TypeSpec) (*StructType, error) {
	if _, isStruct := nodeType.Type.(*ast.StructType); !isStruct {
		return nil, fmt.Errorf("type %s is not a struct", nodeType.Name.Name)
	}
	return &StructType{
		Name:    nodeType.Name.Name,
		Methods: make([]*FuncDecl, 0),
	}, nil
}

type FuncDecl struct {
	Name      string
	Docs      string
	Receiver  string
	IsPointer bool
	FuncDecl  *ast.FuncDecl
}

func NewFunc(nodeType *ast.FuncDecl) (*FuncDecl, error) {
	if nodeType.Name.Name == "init" {
		return nil, fmt.Errorf("skip init function")
	}
	method := &FuncDecl{
		Name: nodeType.Name.Name,
		Docs: nodeType.Doc.Text(),
	}
	if nodeType.Recv != nil && len(nodeType.Recv.List) > 0 {
		recv := nodeType.Recv.List[0]
		var receiverType string
		isPointer := false
		switch t := recv.Type.(type) {
		case *ast.Ident:
			receiverType = t.Name
			isPointer = false
		case *ast.StarExpr:
			if ident, ok := t.X.(*ast.Ident); ok {
				receiverType = ident.Name
				isPointer = true
			}
		default:
		}
		method.Receiver = receiverType
		method.IsPointer = isPointer
		method.FuncDecl = nodeType
	}
	return method, nil
}

// NewFileAst 解析单个 Go 文件
func NewFileAst(filename string, src io.Reader) (*File, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("解析文件 %s 失败: %w", filename, err)
	}
	return &File{
		packageName: file.Name.Name,
		file:        file,
	}, nil
}

func NewpPckageAst(dir string) {

}

func (p *File) ParseFile() {
	var methods []*FuncDecl
	structMap := make(map[string]string)
	// 第一遍：收集所有结构体定义
	ast.Inspect(p.file, func(n ast.Node) bool {
		switch nodeType := n.(type) {
		case *ast.GenDecl:
			if nodeType.Tok == token.TYPE {
				for _, spec := range nodeType.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
							obj, _ := NewStruct(typeSpec)
							comment := typeSpec.Doc.Text()
							if comment == "" {
								comment = nodeType.Doc.Text()
							}
							structMap[obj.Name] = comment
						}
					}
				}
			}
		case *ast.TypeSpec:
			if obj, err := NewStruct(nodeType); err == nil {
				p.Objects = append(p.Objects, obj)
			}
		case *ast.FuncDecl:
			if method, err := NewFunc(nodeType); err == nil {
				methods = append(methods, method)
			}
		}
		return true
	})

	for _, s := range p.Objects {
		if docs, ok := structMap[s.Name]; ok {
			s.Docs = docs
		}
		for i := 0; i < len(methods); {
			if methods[i].Receiver == s.Name {
				s.Methods = append(s.Methods, methods[i])
				methods = append(methods[:i], methods[i+1:]...)
			} else {
				i++
			}
		}
	}

	p.Function = append(p.Function, methods...)
}
