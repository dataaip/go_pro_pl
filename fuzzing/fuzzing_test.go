package fuzzing_test

import (
	"testing"
	"unicode/utf8"

	"github.com/brightlau/go_pro_pl/fuzzing"
)

/*
TODO 1、在 fuzzing 目录中创建一个名为 fuzzing_test.go 的文件
TODO 2、添加单元测试 TestReverse
TODO 3、添加模糊测试 FuzzReverse
*/

// TODO 这个简单的测试将断言列出的输入字符串将被正确地反转
func TestReverse(t *testing.T) {
	// 创建一个 测试样例 结构体 in , want
	testcases := []struct{ in, want string }{
		{"Hello, world", "dlrow ,olleH"},
		{" ", " "},
		{"!12345", "54321!"},
	}
	// for 循环调用 fuzzing.Reverse 反转
	for _, tc := range testcases {
		rev := fuzzing.Reverse(tc.in)
		// 判断 反转后的结果 rev 与 测试样例中的 tc.want 是否一致
		if rev != tc.want {
			t.Errorf("Reverse: %q, want %q", rev, tc.want)
		}
	}
}

// TODO 模糊测试（Fuzzing）的实现，使用了 Go 语言的 testing.F 框架。模糊测试是一种自动化测试技术，通过生成大量随机输入数据来发现程序中的潜在错误或漏洞
// TODO 工作流程：
// - 初始化种子测试用例 将 "Hello, world", " ", "!12345" 添加到模糊测试的种子集合中
// - 运行模糊测试 模糊测试工具会基于种子生成大量随机输入字符串，并调用回调函数进行测试，每次调用都会对随机输入字符串执行 Reverse 函数，并验证其行为是否符合预期
// - 发现问题 如果在测试过程中发现任何不符合预期的情况（如 orig != doubleRev 或 UTF-8 编码被破坏），模糊测试工具会记录问题并提供导致问题的输入数据
// TODO 模糊测试的优点在于它能够自动生成大量随机输入数据，覆盖开发者可能未考虑到的边界情况，从而提高代码的健壮性和可靠性

// TODO 参数 f *testing.F 是模糊测试框架提供的上下文对象，用于管理测试用例和模糊测试的行为
func FuzzReverse(f *testing.F) {
	// TODO 测试用例种子：testcases 是一个字符串切片，包含了一些初始测试用例（称为 种子测试用例）。种子测试用例是模糊测试的起点，帮助模糊测试工具生成更多类似的输入数据
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		// TODO f.Add(tc) 将每个种子测试用例添加到模糊测试的输入集合中。这些种子会被模糊测试工具用来生成更多的变异输入（mutations）
		f.Add(tc)
	}
	// TODO f.Fuzz 方法 f.Fuzz 是模糊测试的核心方法，用于定义模糊测试的具体逻辑。它接收一个回调函数作为参数，该回调函数会被多次调用，每次传入不同的随机输入数据
	// TODO 回调函数的签名为 func(t *testing.T, orig string)：t *testing.T 是标准的测试上下文，用于报告测试结果。orig string 是模糊测试生成的随机输入字符串
	f.Fuzz(func(t *testing.T, orig string) {
		// TODO 测试逻辑 两次反转字符串 两次反转应该恢复原始字符串，即 orig == doubleRev
		rev := fuzzing.Reverse(orig)
		doubleRev := fuzzing.Reverse(rev)
		// TODO 功能正确性：检查两次反转是否恢复原字符串，如果 orig 和 doubleRev 不相等，说明 Reverse 函数存在逻辑错误。使用 t.Errorf 报告错误，并打印原始字符串和最终结果
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		// TODO 编码正确性：检查 UTF-8 编码有效性，如果 orig 是有效的 UTF-8 字符串，但 rev 不是，则说明 Reverse 函数破坏了 UTF-8 编码规则。使用 utf8.ValidString 检查字符串是否符合 UTF-8 编码
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
	})
}
