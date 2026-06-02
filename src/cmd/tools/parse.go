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

type Ast struct {
	Package     string
	Imports     []*ImportSpec
	File        *ast.File
	StructTypes []*StructType
	FuncDecl    []*FuncDecl
}

type ImportSpec struct {
	Name *ast.Ident
	Path *ast.BasicLit
}

type StructType struct {
	Name      string
	Docs      string
	FuncDecls []*FuncDecl
}

func NewStructType(nodeType *ast.TypeSpec) (*StructType, error) {
	if _, isStruct := nodeType.Type.(*ast.StructType); !isStruct {
		return nil, fmt.Errorf("type %s is not a struct", nodeType.Name.Name)
	}
	return &StructType{
		Name:      nodeType.Name.Name,
		FuncDecls: make([]*FuncDecl, 0),
	}, nil
}

type FuncDecl struct {
	Name     string
	Docs     string
	Receiver string
	FuncDecl *ast.FuncDecl
}

func NewFuncDecl(nodeType *ast.FuncDecl) (*FuncDecl, error) {
	if nodeType.Name.Name == "init" {
		return nil, fmt.Errorf("skip init function")
	}
	method := &FuncDecl{
		Name: nodeType.Name.Name,
		Docs: nodeType.Doc.Text(),
	}
	if nodeType.Recv != nil && len(nodeType.Recv.List) > 0 {
		recv := nodeType.Recv.List[0]
		switch t := recv.Type.(type) {
		case *ast.Ident:
			method.Receiver = t.Name
		case *ast.StarExpr:
			if ident, ok := t.X.(*ast.Ident); ok {
				method.Receiver = "*" + ident.Name
			}
		}
		method.FuncDecl = nodeType
	}
	return method, nil
}

// NewFileAst 解析单个 Go 文件
func NewFileAst(filename string, src io.Reader) (*Ast, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("解析文件 %s 失败: %w", filename, err)
	}
	return &Ast{
		Package: file.Name.Name,
		File:    file,
	}, nil
}

func NewPackageAst(dir string) {

}

func (p *Ast) ParseFile() {
	var methods []*FuncDecl
	structMap := make(map[string]string)

	ast.Inspect(p.File, func(n ast.Node) bool {
		switch nodeType := n.(type) {
		case *ast.ImportSpec:
			p.Imports = append(p.Imports, &ImportSpec{
				Name: nodeType.Name,
				Path: nodeType.Path,
			})
		case *ast.GenDecl:
			if nodeType.Tok == token.TYPE {
				for _, spec := range nodeType.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
							obj, _ := NewStructType(typeSpec)
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
			if obj, err := NewStructType(nodeType); err == nil {
				p.StructTypes = append(p.StructTypes, obj)
			}
		case *ast.FuncDecl:
			if method, err := NewFuncDecl(nodeType); err == nil {
				methods = append(methods, method)
			}
		}
		return true
	})

	for _, s := range p.StructTypes {
		if docs, ok := structMap[s.Name]; ok {
			s.Docs = docs
		}
		for i := 0; i < len(methods); {
			receiverName := methods[i].Receiver
			// 跳过没有接收器的普通函数
			if receiverName == "" {
				i++
				continue
			}
			// 去掉指针前缀
			if receiverName[0] == '*' {
				receiverName = receiverName[1:]
			}
			if receiverName == s.Name {
				s.FuncDecls = append(s.FuncDecls, methods[i])
				methods = append(methods[:i], methods[i+1:]...)
			} else {
				i++
			}
		}
	}

	p.FuncDecl = append(p.FuncDecl, methods...)
}
