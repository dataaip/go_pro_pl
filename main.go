package main

/*
声明一个 main 包（包是一种分组的方式 函数，由同一目录中的所有文件组成）
Go 代码被分组到包中，包又被分组到模块中。模块指定了运行代码所需的依赖项，包括 Go 版本以及它所需的其他模块集

go mod init <module-name> 初始化一个新的 Go 模块，<module-name> 是模块的路径，通常是代码托管平台上的仓库路径（如 github.com/username/repo），如果省略 <module-name>，Go 会根据当前目录名自动生成模块路径，创建一个 go.mod 文件，用于定义模块的依赖关系，go.mod 文件是 Go 模块的核心文件，记录了模块名称和依赖版本

go get <package-path>@<version> 添加或更新模块依赖项，<package-path> 是要安装的包路径，<version> 是可选的版本号（如 v1.2.3 或 latest）。如果未指定版本，默认使用最新版本，下载并添加指定的依赖到 go.mod 和 go.sum 文件中，更新现有依赖项的版本
go mod tidy 清理和同步模块依赖，移除未使用的依赖项，添加缺失的依赖项，确保 go.mod 和 go.sum 文件与代码中的实际依赖一致，当代码中删除了某些依赖后，运行此命令可以清理多余的依赖，当代码中新增了依赖但未显式添加时，运行此命令可以自动补充

go mod edit [flags] 手动编辑 go.mod 文件，-require=<module>@<version>: 添加或修改依赖项，-exclude=<module>@<version>: 排除特定版本的依赖，-droprequire=<module>: 删除指定的依赖项，-replace=<old-module>=<new-module>: 替换依赖项，go mod edit -require=github.com/gin-gonic/gin@v1.8.1

go build [flags] [packages] 编译 Go 包及其依赖项，编译指定的包及其依赖项，不会安装编译结果（即不会将生成的二进制文件移动到 $GOPATH/bin 或其他安装路径），默认情况下，生成的二进制文件会放在当前目录下，go build main.go 编译 main.go 文件，生成可执行文件（默认文件名为当前目录名）
go install [flags] [packages] 编译并安装 Go 包及其依赖项，编译指定的包及其依赖项，将生成的二进制文件安装到 $GOPATH/bin 或 $GOBIN 目录中，通常用于安装可执行工具或 CLI 应用程序，go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2 安装 golangci-lint 工具，并将其二进制文件放到 $GOPATH/bin 中

go env $GOPATH 和 $GOBIN 的设置会影响 go install 的行为。可以通过 go env 查看当前配置，go env -w GOBIN=/path/to/your/bin 更改路径
*/

import (
	"fmt"
	// TODO 导入热门 fmt 包 ， 其中包含格式化文本的功能，包括打印到 控制台。此包是 您获得的标准库包 当你安装 Go 时
	"log"
	// TODO 配置log包，使其在日志消息开头打印命令名称（“greetings:”），但不包含时间戳或源文件信息

	"github.com/brightlau/go_pro_pl/datatype"
	// 导入其他 datatype package 包

	"github.com/greetings"
	// TODO 导入其他 go mod 模块
	// 编辑 github.com/go_test/go_lab_web 模块以使用本地 github.com/greetings 模块，对于生产环境，您需要从其代码库中发布 github.com/greetings 模块（模块路径应反映其发布位置），以便 Go 工具能够找到并下载该模块。目前，由于您尚未发布该模块，因此需要调整 github.com/go_test/go_lab_web 模块，使其能够找到 本地文件系统上的 github.com/greetings 代码
	// 使用 go mod edit 命令来编辑 github.com/go_test/go_lab_web 模块将 Go 工具从其模块路径（模块不在其中）重定向 到本地目录（它所在的位置）, go mod edit -replace github.com/greetings=./greetings
	"rsc.io/quote"
	// 导入已发布的模块
	// go mod tidy 或 go get rsc.io/quote@latest 引入
)

/*
Go 语言中使用 fmt.Sprintf 或 fmt.Printf 格式化字符串并赋值给新串：
	- Sprintf 根据格式化参数生成格式化的字符串并返回该字符串。
	- Printf 根据格式化参数生成格式化的字符串并写入标准输出。
*/

var x int = 1

const Pi float32 = 3.14

func main() { // 实现一个 main 函数，用于将消息打印到控制台。运行 main 包时，默认会执行 main 函数
	fmt.Println("Hello, " + "Go!")
	fmt.Println("Hello Go!")
	var str = fmt.Sprintf("Hello,%d,%f", x, Pi)
	fmt.Println(str)
	// %d 表示整型数字，%s 表示字符串，%f 表示浮点
	var stock_code int = 123
	var end_date string = "2020-12-31"
	var url string = "Code=%d&endDate=%s"
	var target_url string = fmt.Sprintf(url, stock_code, end_date)
	fmt.Println(target_url)

	fmt.Println(quote.Go())

	datatype.Data_type()

	// 配置log包，使其在日志消息开头打印命令名称（“greetings:”），但不包含时间戳或源文件信息
	log.SetPrefix("greetings: ")
	log.SetFlags(0)
	// v1.0 获取 greeting message 返回值并打印它.
	// message := greetings.Hello("Gladys")
	// v2.0 为 Hello 返回值赋值，包括 错误 ，变量，将 Hello 参数从 Gladys 的名字更改为空 字符串，以便您可以尝试错误处理代码，查找非 nil 错误值
	message, err := greetings.Hello("Gladys")
	if err != nil {
		// 使用标准库的日志包 输出错误信息。如果得到错误，则使用 日志包的 致命功能 打印错误并停止程序
		log.Fatal(err)
	}
	// 返回值并打印它
	fmt.Println(message)

	// 一组名字
	names := []string{"Gladys", "Samantha", "Darrin"}
	// 请求获取这些名称的问候消息
	messages, err := greetings.Hellos(names)
	if err != nil {
		log.Fatal(err)
	}
	// 如果未返回错误，则将返回的消息映射打印到控制台
	fmt.Println(messages)
}

// TODO go build：从 go_lab_web 目录中的命令行运行 go build 命令将代码编译为可执行文件，从 go_lab_web 目录中的命令行，运行新的 go_lab_web 可执行以确认代码是否工作
// 你已经将应用程序编译为可执行文件，因此可以运行它。但是目前要运行它，你的命令提示符要么在可执行文件的目录中，要么需要指定可执行文件的路径。接下来，你将安装可执行文件，这样就可以在不指定其路径的情况下运行它
// TODO 查找 Go 安装路径：go命令会在此路径安装当前包。你可以通过运行 go list 命令来查找安装路径，如下例所示：$ go list -f '{{.Target}}' 例如，命令的输出可能是 /Users/minghui.liu/go/bin/go_lab_web，这意味着二进制文件被安装到 /Users/minghui.liu/go/bin
// TODO 将 Go 安装目录添加到系统的 shell 路径中：您就可以运行程序的可执行文件，而无需指定可执行文件的位置，在 Linux 或 Mac 上，运行以下命令 $ export PATH=$PATH:/path/to/your/install/directory
// TODO 如果您已经有一个目录， $HOME/bin 在您的 shell 路径中，您希望安装您的 Go 程序，您可以通过设置 GOBIN 变量使用 go env -w GOBIN=/path/to/your/bin
// TODO go install：更新 shell 路径后，运行 go install 命令 编译并安装软件包，只需键入应用程序的名称即可运行应用程序
