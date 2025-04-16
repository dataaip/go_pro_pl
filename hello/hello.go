package hello

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
	// TODO Go 命令在 www.example.com 指定的 hello 目录中查找命令行中指定的 www.example.com 模块 文件，并类似地使用 www.example.com 文件解析 golang.org/x/example/hello/reverse 导入
	// go.work 指令跨多个模块工作。由于这两个模块位于同一个工作区中，因此很容易在一个模块中进行更改并在另一个模块中使用它
	// 现在，为了正确发布这些模块，我们需要发布golang.org/x/example/hello模块，例如发布为v0.1.0版本。这通常是通过在模块的版本控制仓库中给一个提交打标签来完成的。查看模块发布工作流文档以获取更多详细信息。一旦发布完成，我们可以在hello/go.mod中增加对golang.org/x/example
)

func Hello() {
	// 多模块开发
	fmt.Println(reverse.String("Hello"))
	fmt.Println(reverse.String("Hello"), reverse.Int(24601))
}
