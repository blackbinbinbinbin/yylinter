package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type ConstParser struct {
	constLinePos int
}

func (c *ConstParser) Init() {
	c.constLinePos = 0
}

func (c *ConstParser) Parse(pass *analysis.Pass, node ast.Node, decl ast.Decl) {
	genDecl := decl.(*ast.GenDecl)
	// 根据 token.Pos 来获取行号
	v := genDecl.Specs[0].(*ast.ValueSpec)
	p := pass.Fset.Position(v.Pos())
	if c.constLinePos == 0 {
		c.constLinePos = p.Line
	} else {
		// 比较当前 line 和上一次 const 所在行号
		if p.Line - c.constLinePos > 1 {
			pass.Reportf(v.Pos(), "golang-rule suggest【golang-rule-2.5.2】")
		}
		c.constLinePos = p.Line
	}

	// 检查 const 常量命名
	specslist := genDecl.Specs
	for _, valuespec := range specslist {
		if v, ok := valuespec.(*ast.ValueSpec) ; ok {
			constName := v.Names[0].String()
			if strings.Contains(constName, "_") {
				pass.Reportf(node.Pos(), "golang-rule suggest【golang-rule-2.5.1】")
			}
		}
	}
}