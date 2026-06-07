package router

import "text/template"

// AnnotationHandler 注解处理器接口
// 所有注解的「解析逻辑」和「生成行为」都实现该接口
type AnnotationHandler interface {
	// Name 注解名称，如 Controller、GetMapping、Middleware
	Name() string
	// Scope 作用域：struct(结构体) / method(方法)
	Scope() string
	// Parse 解析注释，提取注解参数，存入元数据
	Parse(comment string, meta *NodeMeta) bool
	// Generate 生成该注解对应的 Go 代码片段（核心：定义注解行为）
	Generate(tpl *template.Template, data map[string]interface{}) (string, error)
}

// NodeMeta 通用节点元数据（结构体/方法共用）
type NodeMeta struct {
	// 结构体维度
	StructName string
	Prefix     string   // @Controller 路由前缀
	Middleware []string // @Middleware 中间件列表

	// 方法维度
	MethodName string
	HttpMethod string // GET/POST
	Path       string // 接口路径
}
