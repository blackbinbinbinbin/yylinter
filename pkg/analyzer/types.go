package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

type Parser interface {
	Init()
	Parse(*analysis.Pass, ast.Node, ast.Decl)
}


type ParserHandler struct {
	initParserList []Parser

	// 需要加载的插件
	FilesParser *FilesParser
	ImportParser *ImportParser
	ConstParser *ConstParser
}

func (p *ParserHandler) LoadParser() {
	// 文件相关检查 parser
	p.FilesParser = &FilesParser{}
	p.AddParser(p.FilesParser)
	// 导入 import 相关 parser
	p.ImportParser = &ImportParser{}
	p.AddParser(p.FilesParser)
	// 常量相关 parser
	p.ConstParser = &ConstParser{}
	p.AddParser(p.ConstParser)

	// 初始化 parser
	p.InitParser()
}

// AddParser 添加 parser 接口类型的指针
func (p *ParserHandler) AddParser(parser Parser) {
	p.initParserList = append(p.initParserList, parser)
}

func (p *ParserHandler) InitParser() {
	for _, item := range p.initParserList {
		item.Init()
	}
}