package main

import (
	"github.com/blackbinbinbinbin/yylinter"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(yyimport.Analyzer)
}