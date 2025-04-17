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
- 探索其他资源

TODO Go 模糊测试文档 https://golang.google.cn/security/fuzz/#requirements 中，未来会增加对更多内置类型的支持

TODO 3、添加测试代码：将添加一个函数来反转字符串 Reverse

TODO 4、添加单元测试：在此步骤中，将为 Reverse 函数编写一个基本单元测试

*/

import (
	"fmt"
)

// TODO 这个函数将接受一个字符串 ，一次循环一个字节 ，并在结尾返回反转的字符串，此代码基于 golang.org/x/example 中的 stringutil.Reverse 函数
func Reverse(s string) string {
	b := []byte(s)
	// TODO for ; ; { ... }循环，多个变量 i, j := , 初始化
	for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
		// TODO 反转 第 n 个字节 与 倒数第 n 个字节
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

// TODO 函数初始化一个字符串，反转它，打印输出，然后重复
func ReversePrint() {
	input := "The quick brown fox jumped over the lazy dog"
	// 调用，运行，打印字符串
	rev := Reverse(input)
	doubleRev := Reverse(rev)
	fmt.Printf("original: %q\n", input)
	fmt.Printf("reversed: %q\n", rev)
	fmt.Printf("reversed again: %q\n", doubleRev)
}
