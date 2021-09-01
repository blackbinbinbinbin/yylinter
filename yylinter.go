package yylinter

import (
	"bufio"
	"fmt"
	"go/ast"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "yylinter finds that introduction of non-compliance with specifications"

const FILE_MAX_LINES = 500

var Analyzer = &analysis.Analyzer{
	Name: "yylinter",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{	// filter needed nodes: visit only them
		(*ast.GenDecl)(nil),
	}

	// 检查文件行数
	var fileNames []string
	for _, f := range pass.Files {
		pos := pass.Fset.PositionFor(f.Pos(), false)
		fileNames = append(fileNames, pos.Filename)
		checkFileMaxLine(pass, pos.Filename, FILE_MAX_LINES)
	}

	inspect.Preorder(nodeFilter, func(node ast.Node) {
		if genDecl, ok := node.(*ast.GenDecl); ok {
			checkImport(pass, node, genDecl)
		}

		return
	})
	return nil, nil
}

func checkImport(pass *analysis.Pass, node ast.Node, genDecl *ast.GenDecl) {
	specslist := genDecl.Specs
	for _, importspec := range specslist {
		ispec := importspec.(*ast.ImportSpec)
		if ispec.Name != nil {
			name := ispec.Name.Name
			if name == "." {
				pass.Reportf(node.Pos(), "【golang-rule-1.3.1】formatting import '%s' should not use '%s' before",
					ispec.Path.Value, ispec.Name.Name)
			}
			importPath := ispec.Path.Value
			if strings.Contains(importPath, "./") {
				pass.Reportf(node.Pos(), "【golang-rule-1.3.2】formatting import path '%s'. It is forbidden to use relative path import (./), all import paths must conform to the go get standard",
					importPath)
			}
		}
	}
}

func checkFileMaxLine(pass *analysis.Pass, fileName string, maxLines int) {
	file, _ := os.Open(fileName)
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}

	if lineCount > maxLines {
		fmt.Printf("【golang-rule-1.5.1】file：%s, max line length=%d，but file line=%d\n", fileName, maxLines, lineCount)
	}
}