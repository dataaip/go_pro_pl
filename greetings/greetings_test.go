package greetings

// TODO 测试：现在，你的代码已经处于稳定状态（顺便说一句，做得很好），添加一个测试。在开发过程中测试你的代码可以暴露在你进行更改时出现的程序缺陷。在本主题中，你将为Hello函数添加一个测试
// Go 的内置单元测试支持使你在开发过程中更容易进行测试。具体来说，使用命名约定、Go 的 testing 包以及go test命令，你可以快速编写和执行测试
// 在 greetings 目录中，创建一个名为 greetings_test. go 的文件。以_test.go 结尾的文件名告诉 go test 命令 这个文件包含测试功能

import (
	"regexp"
	// TODO 正则包：regexp 是 Go 标准库中的一个包，用于处理正则表达式。它提供了强大的工具来匹配、搜索和替换字符串中的模式。例如：使用正则表达式验证输入格式（如邮箱、电话号码等）、在文本中查找特定模式的子串、替换文本中符合正则表达式的部分
	"testing"
	// TODO testing 包：testing 是 Go 标准库中用于编写单元测试的包。它提供了一组工具来定义测试函数、基准测试函数以及示例测试函数。测试文件通常以 _test.go 结尾，并通过 go test 命令运行
)

// TODO 测试T testing.T： 在测试代码所在的包中实现测试函数，创建两个测试函数来测试 greeting.hello 函数。测试函数名称的形式为 TestName，其中 Name 表示有关特定测试的内容。此外，测试函数还使用指向测试包的 testing.T 类型作为参数。您可以使用此参数的方法进行报告 和测试日志
// 实施两项测试：TestHelloName 调用 Hello 函数，传递一个名称值，函数应该能够使用该值返回有效的响应消息。如果调用返回错误或意外的响应消息（不包括传入的名称），则使用 t 参数的 方法将消息打印到控制台
// TestHelloEmpty 使用空字符串调用 Hello 函数。此测试旨在确认您的错误处理是否有效。如果调用返回一个非空字符串或没有错误，则使用 t 参数的 方法将消息打印到控制台
// TODO go test：在 greetings 目录中的命令行中，运行 执行测试命令 来执行测试，go test 命令执行测试文件（其名称以_test.go 结尾）中的测试函数（其名称开始以 Test 开头）。您可以添加 -v 标志以获得详细输出， 列出了所有的测试和结果，测试应该通过
// $ go test
// PASS
// ok      github.com/greetings   0.364s
// $ go test -v
// === RUN   TestHelloName
// --- PASS: TestHelloName (0.00s)
// === RUN   TestHelloEmpty
// --- PASS: TestHelloEmpty (0.00s)
// PASS
// ok      github.com/greetings   0.372s
// TODO 中断 greeting.Hello 功能以查看失败的测试，TestHelloName 测试函数检查您指定为 Hello 函数参数的名称的返回值。若要查看失败的测试结果，请将 greeting.Hello 函数更改为 它不再包括名字
// $ go test
// --- FAIL: TestHelloName (0.00s)
// greetings_test.go:15: Hello("Gladys") = "Hail, %v! Well met!", <nil>, want match for `\bGladys\b`, nil
// FAIL
// exit status 1
// FAIL    example.com/greetings   0.182s

// TestHelloName 调用 greetings.Hello 并传入一个名字，检查是否有有效的返回值
// TODO 参数 t *testing.T 是测试框架提供的一个对象，用于记录测试失败的信息或标记测试状态
func TestHelloName(t *testing.T) {
	name := "Gladys"
	// TODO 正则表达式对象 want: 使用 regexp.MustCompile 创建了一个正则表达式对象，匹配单词边界（\b）包围的 name 值，\b 是正则表达式中的单词边界锚点，确保匹配的是完整的单词，而不是部分子串
	want := regexp.MustCompile(`\b` + name + `\b`)
	// 调用了 Hello 函数，并传入参数 "Gladys"
	msg, err := Hello("Gladys")
	// TODO 断言测试条件
	// TODO 使用正则表达式：want.MatchString(msg)，使用正则表达式 want 检查返回的消息 msg 是否包含完整的单词 "Gladys"，如果 msg 中不包含 "Gladys"，或者包含的部分不符合单词边界规则，则条件为假
	// 检查 Hello 函数是否返回了错误。如果 err 不为 nil，说明函数执行过程中出现了问题
	// t.Errorf 如果上述条件之一不满足（即测试失败），使用 t.Errorf 记录错误信息
	if !want.MatchString(msg) || err != nil {
		// %q: 打印 msg 的值（带引号的字符串），%v: 打印 err 的值，%#q: 打印正则表达式 want 的值（带转义字符的字符串）
		t.Errorf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}

// TestHelloEmpty 调用 greetings.Hello 函数传入空字符串，检查是否有错误。
func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")
	if msg != "" || err == nil {
		t.Errorf(`Hello("") = %q, %v, want "", error`, msg, err)
	}
}
