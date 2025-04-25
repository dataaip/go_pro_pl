package code_organization

import (
	"fmt"

	"github.com/brightlau/go_pro_pl/code_organization/morestrings"

	"github.com/google/go-cmp/cmp"
)

/*
TODO 如何编写 Go 代码

Go 语言表达力强、简洁、清晰且高效。其并发机制使得编写充分利用多核和网络化机器性能的程序变得容易，而其新颖的类型系统则支持灵活且模块化的程序构建。Go 可快速编译为机器码，同时兼具垃圾回收的便利性和运行时反射的强大功能。它是一门快速、静态类型的编译型语言，但却给人一种动态类型、解释型语言的感觉

目录
- 介绍
- 代码组织
- 第一个程序
- 从模块导入包
- 从远程模块导入包
- 测试
- 下一步
- 获取帮助
*/

// TODO 介绍：本文档演示了如何在模块中开发一个简单的 Go 包，并介绍了 go 工具 https://golang.google.cn/cmd/go/ 这是获取、构建和安装 Go 模块、包和命令的标准方法

// TODO 代码组织：
// - Go 程序由包组成：一个包是同一目录中一起编译的源文件的集合。在一个源文件中定义的函数、类型、变量和常量，对同一包中的所有其他源文件都是可见的
// - 一个存储库包含一个或多个模块：模块是一组相关的 Go 包，它们一起发布。一个 Go 存储库通常只包含一个模块，位于存储库的根目录。在那里，一个名为 go.mod 的文件声明了模块路径：模块内所有包的导入路径前缀。该模块包含其 go.mod 文件所在目录中的包，以及该目录的子目录中的包，直到下一个包含另一个 go.mod 文件的子目录（如果有）
// - 请注意，在构建代码之前，你无需将其发布到远程存储库。模块可以在本地定义，而无需属于某个存储库。但是，养成一种习惯，即假设自己有一天会发布代码那样来组织你的代码
// - 每个模块的路径不仅用作其包的导入路径前缀，还指示 go 命令 应在何处查找以下载它，例如，为了下载模块 golang.org/x/tools go 命令将查询由 https://golang.org/x/tools
// - 导入路径是用于导入包的字符串：一个包的导入路径是其模块路径与其在模块中的子目录相结合的结果。例如，模块 github.com/google/go-cmp 包含了一个位于 cmp/ 目录下的包。该包的导入路径为 github.com/google/go-cmp/cmp。标准库中的包没有模块路径前缀

// TODO 第一个程序：Code_organization
// 1）要编译并运行一个简单的程序，首先选择一个模块路径（我们将使用 example/user/hello）并创建一个 go.mod 文件来声明它
// $ mkdir hello # Alternatively, clone it if it already exists in version control.
// $ cd hello
// $ go mod init example/user/hello
// go: creating new go.mod: module example/user/hello
// $ cat go.mod
// module example/user/hello
//
// go 1.16
// $
// 2）Go 源文件中的第一条语句必须是 package name 可执行命令必须始终使用 package main
// 3）在该目录中创建一个名为 hello.go 的文件，其中包含以下 Go 代码
// package main
//
// import "fmt"
//
// func main() {
//     fmt.Println("Hello, world.")
// }
// 4）可以使用 go 工具构建和安装该程序
// go install example/user/hello
// 此命令构建 hello 命令，生成一个可执行二进制文件。然后将该二进制文件安装到 $HOME/go/bin/hello（在 Windows 下为 %USERPROFILE%\go\bin\hello.exe）
// 安装目录由 GOPATH 和 GOBIN 环境变量控制。如果设置了 GOBIN，则二进制文件将安装到该目录。如果设置了 GOPATH，则二进制文件将安装到 GOPATH 列表中第一个目录的 bin 子目录中。否则，二进制文件将安装到默认 GOPATH 的 bin 子目录（$HOME/go 或 %USERPROFILE%\go）
//
// 可以使用 go env 命令以便携方式为将来的 go 命令设置环境变量的默认值
// go env -w GOBIN=/somewhere/else/bin
//
// 要取消设置之前由 go env -w 设置的变量，请使用 go env -u
// go env -u GOBIN
//
// 像 go install 这样的命令会在包含当前工作目录的模块上下文中执行。如果工作目录不在 example/user/hello 模块中，go install 可能会失败
//
// 为方便起见 go 命令接受相对于工作目录的路径，如果未指定其他路径，则默认为当前工作目录中的包。因此在我们的工作目录中，以下命令都是等效的
// $ go install example/user/hello
// $ go install .
// $ go install
//
// 5）接下来，运行该程序以确保它正常工作。为了更加方便，我们将 install 目录添加到 PATH 中，以便轻松运行二进制文件
// # Windows users should consult /wiki/SettingGOPATH
// # for setting %PATH%.
// $ export PATH=$PATH:$(dirname $(go list -f '{{.Target}}' .))
// $ hello
// Hello, world.
// $
//
// 6）如果你正在使用版本控制系统，现在是初始化仓库、添加文件并提交首次更改的好时机。再次说明，此步骤是可选的：编写 Go 代码并不需要使用版本控制系统
// $ git init
// Initialized empty Git repository in /home/user/hello/.git/
// $ git add go.mod hello.go
// $ git commit -m "initial commit"
// [master (root-commit) 0b4507d] initial commit
//  1 file changed, 7 insertion(+)
//  create mode 100644 go.mod hello.go
// $
//

// TODO importpath：go 命令通过请求相应的 HTTPS URL 并读取 HTML 响应中嵌入的元数据，来定位包含给定模块路径的存储库（请参阅go help importpath）https://golang.google.cn/cmd/go/#hdr-Remote_import_paths 许多托管服务已经为包含 Go 代码的存储库提供了该元数据，因此，让其他人能够轻松使用你的模块的最简单方法通常是使其模块路径与存储库的 URL 相匹配

// TODO 从本模块 hello 导入包
// 1）让我们编写一个 morestrings 包，并在 hello 程序中使用它。首先，为该包创建一个名为 $HOME/hello/morestrings 的目录，然后在该目录中创建一个名为 reverse.go 的文件，其内容如下
// package morestrings
//
// func ReverseRunes(s string) string {
//     r := []rune(s)
//     for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
//         r[i], r[j] = r[j], r[i]
//     }
//     return string(r)
// }
// 因为 ReverseRunes 函数以大写字母开头，所以是导出的，并且可以在导入 morestrings 包的其他包中使用
//
// 2）测试一下这个包是否用 go build 编译，这不会生成输出文件。相反，它将编译后的包保存在本地构建缓存中
// $ cd $HOME/hello/morestrings
// $ go build
//
// 3）确认 morestrings 包构建成功后，就可以从 hello 程序中使用它。为此，修改原始的 $HOME/hello/hello.go 以使用 morestrings 包
// package main
//
// import (
//     "fmt"
//     "example/user/hello/morestrings"
// )
//
// func main() {
//     fmt.Println(morestrings.ReverseRunes("!oG ,olleH"))
// }
//
// 4）安装 hello 程序
// go install example/user/hello
//
// 5）运行该程序的新版本，应该会看到一条新的反向消息
// $ hello
// Hello, Go!

// TODO 从远程模块导入包
// 1) 导入路径可以描述如何使用版本控制系统（例如 Git 或 Mercurial）获取包的源代码。go 工具利用这一特性自动从远程仓库中获取包。例如，要在程序中使用 github.com/google/go-cmp/cmp
//
// 2) 既然你依赖一个外部模块，就需要下载该模块，并在你的 go.mod 文件中记录其版本。go mod tidy 命令会为导入的包添加缺失的模块要求，并移除不再使用的模块的要求
// $ go mod tidy
// go: finding module for package github.com/google/go-cmp/cmp
// go: found github.com/google/go-cmp/cmp in github.com/google/go-cmp v0.5.4
// $ go install example/user/hello
// $ hello
// Hello, Go!
//   string(
// -     "Hello World",
// +     "Hello Go",
//   )
// $ cat go.mod
// module example/user/hello
//
// go 1.16
//
// require github.com/google/go-cmp v0.5.4
// $
//
// 3) 模块依赖项会自动下载到由 pkg/mod 环境变量指示的目录的 GOPATH 子目录中。给定模块版本的下载内容会在所有其他 require 该版本的模块之间共享，因此 go 命令会将这些文件和目录标记为只读。要删除所有下载的模块，可以将 -modcache 标志传递给 go clean
// $ go clean -modcache
// $

// TODO 测试
// 1）Go 拥有一个轻量级测试框架，它由 go test 命令和 testing 包组成
//
// 2）通过创建一个文件名以 _test.go 结尾的文件来编写测试，该文件包含名为 TestXXX 且签名为 func (t *testing.T) 的函数。测试框架会运行每个这样的函数；如果该函数调用了诸如 t.Error 或 t.Fail 等失败函数，则认为测试失败
//
// 3）通过创建包含以下 Go 代码的文件 $HOME/hello/morestrings/reverse_test.go 向 morestrings 包添加一个测试
// package morestrings
// import "testing"
// func TestReverseRunes(t *testing.T) {
//     cases := []struct {
//         in, want string
//     }{
//         {"Hello, world", "dlrow ,olleH"},
//         {"Hello, 世界", "界世 ,olleH"},
//         {"", ""},
//     }
//     for _, c := range cases {
//         got := ReverseRunes(c.in)
//         if got != c.want {
//             t.Errorf("ReverseRunes(%q) == %q, want %q", c.in, got, c.want)
//         }
//     }
// }
//
// 4）使用 go test 运行测试
// $ cd $HOME/hello/morestrings
// $ go test
// PASS
// ok  	example/user/hello/morestrings 0.165s
// $
// 运行 go help test https://golang.google.cn/cmd/go/#hdr-Test_packages 并查看 testing 包文档 https://golang.google.cn/pkg/testing/ 了解更多详细信息

// TODO 参阅《Effective Go》https://golang.google.cn/doc/effective_go.html 以获取编写清晰、符合惯用法的 Go 代码的建议
//
// TODO 参加《A Tour of Go》https://golang.google.cn/tour/ 以学习该语言的正确用法
//
// TODO 访问文档页面 https://golang.google.cn/doc/#articles 以获取有关 Go 语言及其库和工具的一系列深度文章

func Code_organization() {
	fmt.Println("Hello, world.")
	fmt.Println(morestrings.ReverseRunes("!oG ,olleH"))
	fmt.Println(cmp.Diff("Hello World", "Hello Go"))
}
