# Golangci-lint 自定义插件
GolangCI-Lint是一个lint聚合器，它的速度很快，平均速度是gometalinter的5倍。它易于集成和使用，具有良好的输出并且具有最小数量的误报。而且它还支持go modules。最重要的是免费开源

在 GolangCI-Lint 中还可以通过自定义的 linter ，来进行所需要的自定义规范检查，如果在官方的 linter 没有找到您所需要的规范检查，可以在本项目基础上开发符合自己项目风格的规范检查

## 程序目录
```
.
├── README.md
├── cmd
│   └── yylinter
│       └── main.go
├── go.mod
├── go.sum
├── pkg
│   └── analyzer
│       ├── constparser.go      // const 静态常量的相关检测
│       ├── filesparser.go      // file 文件相关的检测
│       ├── importparser.go     // import 相关规范检测
│       ├── types.go            // 加载需要检测的 parser
│       ├── yylinter.go         // 处理被检测文件的主流程逻辑
│       └── yylinter_test.go    // 测试文件
├── plugin
└── testdata
    └── src
        └── a
            └── a.go

```
- cmd 可执行文件入口
- pkg 进行规范插件的主要逻辑部分
- testdata 测试相关文件

## 安装
```
go get -u github.com/blackbinbinbinbin/yylinter
```
如果需要对源码进行调试修改，或者新增 parser，可以直接将修改的项目源码直接复制到本地的路径：
```
$GOPATH/pkg/mod/github.com/blackbinbinbinbin/yylinter@v0.0.0-XXXX
```
在此路径下放入您修改的源码即可

## 使用
```
go run ./cmd/yylinter/main.go -- ./testdata/src/a/a.go
```
执行 main 文件，后面跟着指定需要检测的源文件

## 新增检查 parser
在文件 `pkg/analyzer/types.go` 中添加新的 parser
```
func (p *ParserHandler) LoadParser() {
    ...
    // 新增 parser
    p.NewParser = &NewParser{}
    p.AddParser(p.NewParser)

    ...
}
```
添加 parser 插件：
```
type ParserHandler struct {
    ...
    // 需要加载的插件
    NewParser *NewParser
    ...
}
```

新增的 `NewParser` 需要实现接口：
```
type Parser interface {
    Init()
    Parse(*analysis.Pass, ast.Node, ast.Decl)
}
```
具体的可以参考其他的 parser

## 集成 golangci-lint 
在此之前，建议您先去看看 golangci-lint 官网，关于 new linter 的介绍文档，里面有比较详细的描述：
[golangci-lint new linter](https://golangci-lint.run/contributing/new-linters/)

而如何编写一个符合规范的 linter，这里有一篇非常好的入门教程：[linter教程](https://disaev.me/p/writing-useful-go-analysis-linter/)

除此之外，因为 linster 是根据 `ast` 和 `analysis` 来进行规范检查开发的，所以如果不熟悉 ast 接口的，可以根据这个网站工具来进行开发：[开发工具](http://goast.yuroyoro.net/)

- 声明 linter

代码仓库地址：https://github.com/golangci/golangci-lint

在 golangci-linter 项目中新增：`pkg/golinters/yylinter.go`
```
package golinters

import (
	"github.com/blackbinbinbinbin/yylinter/pkg/analyzer"
	"github.com/golangci/golangci-lint/pkg/config"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"golang.org/x/tools/go/analysis"
)

func NewYylinter(settings *config.YylinterSettings) *goanalysis.Linter {
	return goanalysis.NewLinter(
		"yylinter",
		"Checks that yy golang code format",
		[]*analysis.Analyzer{analyzer.Analyzer},
		nil,
	).WithLoadMode(goanalysis.LoadModeSyntax)
}
```
这里调用 `goanalysis.NewLinter` 方法，声明 linter 名称，描述，和传入的配置参数


- 新增 linter 配置

在 `pkg/lint/lintersdb/manager.go` 中添加 linter 相关的配置
```
func (m Manager) GetAllSupportedLinterConfigs() []*linter.Config {
    ...
    // 声明配置
    var yylinterCfg *config.YylinterSettings
    ...
    
    if m.cfg != nil {
        ...
        yylinterCfg = &m.cfg.LintersSettings.Yylinter
    }
    
    // 初始化配置
    lcs := []*linter.Config{
        ...
        // yylinter
        linter.NewConfig(golinters.NewYylinter(yylinterCfg)).
            WithSince("v1.26.0").
            WithPresets(linter.PresetStyle).
            WithURL("https://github.com/blackbinbinbinbin/yylinter"),
    }
}
```

- 添加测试

golangci-linter 要求新加入的 linter 必须测试能够通过，所以在新增的 linter 中，需要测试文件全部通过。比如当前文件 `pkg/analyzer/yylinter_test.go`

添加测试用例进入测试目录：`test/testdata/yylinter.go`
```
package a	// want "golang-rule error【golang-rule-1.5.1】"

import (	//want "golang-rule error【golang-rule-1.3.1】"
	. "fmt"
)

const ABC = 1
const EFGF = 2
var d int
const EFG = 3	//want "golang-rule suggest【golang-rule-2.5.2】"
const ABC_EFG = 4	//want "golang-rule suggest【golang-rule-2.5.1】"

func main() {
	aaaa := 0
	Println(aaaa)
	Println("yylinter testdata")
}
```
文件内容就是此目录下的 `pkg/analyzer/yylinter_test.go`

保证测试用例通过，然后运行
```
go run ./cmd/golangci-lint/ run --no-config --disable-all --enable=yourlintername ./test/testdata/yourlintername.go
```

- 打包
```
go build -o golangci-lint ./cmd/golangci-lint 
```
`golangci-lint` 生成的可执行文件，可以直接使用