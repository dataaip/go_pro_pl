package say_hello

import (
	"fmt"
	// TODO 导入热门 fmt 包 ， 其中包含格式化文本的功能，包括打印到 控制台。此包是 您获得的标准库包 当你安装 Go 时

	"golang.org/x/example/hello/reverse"
	// TODO Go 命令在 www.example.com 指定的 hello 目录中查找命令行中指定的 www.example.com 模块 文件，并类似地使用 www.example.com 文件解析 golang.org/x/example/hello/reverse 导入
	// go.work 指令跨多个模块工作。由于这两个模块位于同一个工作区中，因此很容易在一个模块中进行更改并在另一个模块中使用它
	// 现在，为了正确发布这些模块，我们需要发布golang.org/x/example/hello模块，例如发布为v0.1.0版本。这通常是通过在模块的版本控制仓库中给一个提交打标签来完成的。查看模块发布工作流文档以获取更多详细信息。一旦发布完成，我们可以在hello/go.mod中增加对golang.org/x/example

	"rsc.io/quote"
	// TODO 这个包收集了精辟的谚语
	// 导入已发布的模块
	// go mod tidy 或 go get rsc.io/quote@latest 引入
)

var x int = 1

const Pi float32 = 3.14

func Say_hello() {
	// Go 语言中使用 fmt.Sprintf 或 fmt.Printf 格式化字符串并赋值给新串：
	// 	- Sprintf 根据格式化参数生成格式化的字符串并返回该字符串
	// 	- Printf 根据格式化参数生成格式化的字符串并写入标准输出
	fmt.Println("Hello, " + "Go!")
	fmt.Println("Hello Go!")
	var str = fmt.Sprintf("Hello,%d,%f", x, Pi)
	fmt.Println(str)

	// %d 表示整型数字，%s 表示字符串，%f 表示浮点，%v 表示 切片
	var stock_code int = 123
	var end_date string = "2020-12-31"
	var url string = "Code=%d&endDate=%s"
	var target_url string = fmt.Sprintf(url, stock_code, end_date)
	fmt.Println(target_url)

	// 多模块开发 go work
	fmt.Println(reverse.String("Hello"))
	fmt.Println(reverse.String("Hello"), reverse.Int(24601))

	// Go 返回 Go 谚语
	fmt.Println(quote.Go())
}
