package analyzer

import (
	"go/ast"
	"go/token"

	"github.com/gostaticanalysis/comment/passes/commentmap"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "yylinter finds that introduction of non-compliance with specifications"

var Analyzer = &analysis.Analyzer{
	Name: "yylinter",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
		commentmap.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{	// filter needed nodes: visit only them
		(*ast.GenDecl)(nil),
		(*ast.FuncDecl)(nil),
	}

	// 加载 parser 插件
	parsers := &ParserHandler{}
	parsers.LoadParser()

	// 检查文件行数
	parsers.FilesParser.Parse(pass, nil, nil)

	inspect.Preorder(nodeFilter, func(node ast.Node) {
		if genDecl, ok := node.(*ast.GenDecl); ok {
			switch genDecl.Tok {
			case token.IMPORT:
				parsers.ImportParser.Parse(pass, node, genDecl)
				break
			case token.CONST:
				parsers.ConstParser.Parse(pass, node, genDecl)
				break
			}
		}

		if _, ok := node.(*ast.FuncDecl); ok {
			// 处理 func 相关
		}

		return
	})
	return nil, nil
}