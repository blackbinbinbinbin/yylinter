package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

type ImportParser struct {
}

func (i *ImportParser) Init() {
}

func (i *ImportParser) Parse(pass *analysis.Pass, node ast.Node, decl ast.Decl) {
	genDecl := decl.(*ast.GenDecl)
	specslist := genDecl.Specs
	for _, importspec := range specslist {
		ispec := importspec.(*ast.ImportSpec)
		if ispec.Name != nil {
			name := ispec.Name.Name
			if name == "." {
				// 不能用 . 来进行引入包
				pass.Reportf(node.Pos(), "golang-rule error【golang-rule-1.3.1】")
			}
			importPath := ispec.Path.Value
			if strings.Contains(importPath, "./") {
				// 引入不能用相对路径
				pass.Reportf(node.Pos(), "golang-rule error【golang-rule-1.3.2】")
			}
		}
	}
}