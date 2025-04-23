package fuzzing

/*
TODO 本教程介绍了 Go 语言中模糊测试的基础知识。通过模糊测试，随机数据会在你的测试中运行，试图找到漏洞或导致崩溃的输入。可以通过模糊测试找到的一些漏洞示例包括 SQL 注入、缓冲区溢出、拒绝服务和跨站脚本攻击
TODO 1、本教程，将为一个简单函数编写模糊测试，运行 go 命令，并调试和修复代码中的问题，有关本教程中术语的帮助，请参阅 Go Fuzzing 词汇表 https://golang.google.cn/security/fuzz/#glossary

TODO 2、将逐步完成以下部分：
- 为您的代码创建一个文件夹
- 添加要测试的代码 Reverse
- 添加单元测试 TestReverse
- 添加模糊测试 FuzzReverse 单元测试有局限性，即每个输入都必须由开发人员添加到测试中。模糊测试的一个好处是它为你的代码生成输入，并且可能识别出你想出的测试用例没有覆盖到的边界情况，在本节中，将把单元测试转换为模糊测试，这样就可以用更少的工作量生成更多的输入，请注意，可以将单元测试、基准测试和模糊测试保存在同一个 “*_test.go” 文件中，但在本示例中，将把单元测试转换为模糊测试
- 修复两个错误
	1）单字节可以按字节反转/多字节不能按字节反转，如果对 Go 语言如何处理字符串感到好奇，请阅读博客文章Go 语言中的字符串、字节、符文和字符 https://golang.google.cn/blog/strings 以获得更深入的理解，更好地理解 bug 后，纠正 Reverse 中的错误 功能
	2）修复双重 reverse 错误，仔细查看反转后的字符串以发现错误。在 Go 语言中，字符串是字节的只读切片，并且可以包含无效的 UTF-8 字节。原始字符串是一个具有一个字节 '\x91' 的字节切片。当输入字符串被设置为[]rune时，Go 将字节切片编码为 UTF-8，并将该字节替换为 UTF-8 字符 � 当我们将替换后的 UTF-8 字符与输入字节切片进行比较时，它们显然不相等
- 探索其他资源

TODO Go 模糊测试文档 https://golang.google.cn/security/fuzz/#requirements 中，未来会增加对更多内置类型的支持

TODO 3、添加测试代码：将添加一个函数来反转字符串 Reverse

TODO 4、添加单元测试：在此步骤中，将为 Reverse 函数编写一个基本单元测试


TODO 5、进行测试：go test 是一个强大的命令行工具，用于运行测试和模糊测试（fuzzing）
go test -run=TestReverse ：-run 参数指定了要运行的测试函数名称（支持正则表达式匹配）。TestReverse 是一个普通的单元测试函数
go test -fuzz=Fuzz -run=FuzzReverse ：-fuzz 参数启用了模糊测试（fuzzing），并指定了要运行的模糊测试函数名称（支持正则表达式匹配）。Fuzz 是模糊测试函数的前缀
go test -run=FuzzReverse ：-run 参数指定了要运行的测试函数名称（支持正则表达式匹配）。FuzzReverse 是一个模糊测试函数，但由于这里没有使用 -fuzz 参数，Go 只会将其视为普通测试函数运行
go test -fuzz=Fuzz -fuzztime 30s -run=FuzzReverse ：-fuzz=Fuzz 启用了模糊测试，并指定了模糊测试函数的前缀为 Fuzz。-fuzztime 30s 设置了模糊测试的最大运行时间为 30 秒。-run=FuzzReverse 表示只运行与正则表达式 FuzzReverse 匹配的测试函数

TODO 6、除了 -fuzz 标志外，几个新的标志已被添加到 go test 中，并且可以在文档中查看 https://golang.google.cn/security/fuzz/#custom-settings
请参阅 Go 模糊测试以了解有关模糊测试输出中使用的术语的更多信息。例如，新的有趣输入 是指那些扩展现有模糊测试语料库代码覆盖范围的输入。随着模糊测试的开始，新的有趣输入 的数量预计会急剧增加，随着新代码路径的发现会多次飙升，然后随着时间的推移逐渐减少
*/

import (
	"errors"
	"fmt"
	"log"
	"unicode/utf8"
)

// TODO 这个函数将接受一个字符串 ，一次循环一个字节 ，并在结尾返回反转的字符串，此代码基于 golang.org/x/example 中的 stringutil.Reverse 函数
func Reverse(s string) (string, error) {
	// TODO 字节切片
	// b := []byte(s)
	// TODO 1）关键区别在于，Reverse 现在正在遍历字符串中的每个 rune 而不是每个byte 请注意，这只是一个示例，并且不能正确处理组合字符 https://en.wikipedia.org/wiki/Combining_character
	// TODO 字符串的本质：在Go语言中，字符串是以 UTF-8 编码存储的字节序列。这意味着字符串中的每个字符可能占用1到4个字节（取决于字符的Unicode编码）
	// TODO 在Go语言中，rune 是一个内置类型，实际上是 int32 的别名。它用来表示一个Unicode码点（即一个字符的数值编码）。通过使用 rune，我们可以方便地处理字符串中的单个字符，而无需担心字符的字节长度，当我们将字符串 s 转换为 []rune 类型时，Go会将字符串中的每个字符解码为对应的Unicode码点，并将这些码点存储在一个 rune 切片中。这样做的好处是可以方便地按字符操作字符串，而不是按字节操作
	// fmt.Printf("input: %q\n", s)
	// TODO 2）为了解决这个问题，如果 Reverse 的输入不是有效的 UTF-8，返回一个错误，如果输入字符串包含无效的 UTF-8 字符，则此更改将返回错误
	if !utf8.ValidString(s) {
		return s, errors.New("input is not valid UTF-8")
	}
	b := []rune(s)
	// fmt.Printf("runes: %q\n", b)
	// TODO for ; ; { ... }循环，多个变量 i, j := , 初始化
	for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
		// TODO 反转 第 n 个字节 与 倒数第 n 个字节
		b[i], b[j] = b[j], b[i]
	}
	return string(b), nil
}

// TODO 函数初始化一个字符串，反转它，打印输出，然后重复
func ReversePrint() {
	input := "The quick brown fox jumped over the lazy dog"
	// 调用，运行，打印字符串
	rev, revErr := Reverse(input)
	if revErr != nil {
		log.Fatal(revErr)
	}
	doubleRev, doubleRevErr := Reverse(rev)
	if doubleRevErr != nil {
		log.Fatal(doubleRevErr)
	}
	fmt.Printf("original: %q\n", input)
	fmt.Printf("reversed: %q\n", rev)
	fmt.Printf("reversed again: %q\n", doubleRev)
}
