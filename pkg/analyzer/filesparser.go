package analyzer

import (
	"bufio"
	"go/ast"
	"os"

	"golang.org/x/tools/go/analysis"
)

// FileMaxLines 最大文件行数
const FileMaxLines = 5

type FilesParser struct {
	MaxLines int
}

func (s *FilesParser) Init() {
	s.MaxLines = FileMaxLines
}

func (s *FilesParser) Parse(pass *analysis.Pass, node ast.Node, decl ast.Decl) {
	var files []*os.File
	files = make([]*os.File, len(pass.Files))
	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()
	for _, f := range pass.Files {
		pos := pass.Fset.PositionFor(f.Pos(), false)
		fileName := pos.Filename
		file, _ := os.Open(fileName)
		fileScanner := bufio.NewScanner(file)
		lineCount := 0
		for fileScanner.Scan() {
			lineCount++
		}

		if lineCount > FileMaxLines {
			pass.Reportf(f.Pos(), "golang-rule error【golang-rule-1.5.1】")
		}
	}
}