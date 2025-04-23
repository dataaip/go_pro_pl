package fuzzing_test

import (
	"testing"

	"unicode/utf8"
	// TODO unicode/utf8 是一个标准库包，专门用于处理UTF-8编码的文本。它提供了一系列工具函数和常量，帮助开发者正确地操作和解析UTF-8编码的数据。下面我们详细解析这个包的功能和用途
	// UTF-8（8-bit Unicode Transformation Format）是一种可变长度的字符编码方式，广泛用于表示 Unicode 字符。它的特点包括：ASCII字符（0-127）用1个字节表示，其他字符（如中文、日文、韩文等）根据Unicode码点的范围，可能占用2到4个字节
	// 由于 UTF-8 编码的可变长度特性，直接操作字节序列可能会导致错误。因此，unicode/utf8 包提供了许多工具来安全地处理这些数据：字符 'A' 的UTF-8编码是 [65]（1个字节），字符 '你' 的UTF-8编码是 [228, 189, 160]（3个字节）
	// 判断是否为有效的UTF-8编码 utf8.Valid(p []byte) bool
	// 计算字符串或字节切片中的字符数量 utf8.RuneCountInString(s string) int 、utf8.RuneCount(p []byte) int
	// 解码单个UTF-8字符 utf8.DecodeRune(p []byte) (r rune, size int)
	// 编码单个UTF-8字符 utf8.EncodeRune(p []byte, r rune) int

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
		rev, revErr := fuzzing.Reverse(tc.in)
		if revErr != nil {
			// log.Fatal(revErr)
		}
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

// TODO 模糊测试时发生故障，导致问题的输入被写入种子语料库文件，该文件将在下次调用 go test 时运行，即使没有 -fuzz 标志也是如此。要查看导致失败的输入，请在文本编辑器中打开写入 testdata/fuzz/FuzzReverse 目录的语料库文件。您的种子语料库文件可能包含不同的字符串，但格式相同
// 语料库文件的第一行表示编码版本。接下来的每一行都表示构成语料库条目的每种类型的值。由于模糊目标只接受 1 个输入，因此版本后只有 1 个值

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
		rev, revErr := fuzzing.Reverse(orig)
		if revErr != nil {
			// log.Fatal(revErr)
		}
		doubleRev, doubleRevErr := fuzzing.Reverse(rev)
		if doubleRevErr != nil {
			// log.Fatal(doubleRevErr)
		}
		// TODO 整个种子语料库使用的字符串中每个字符都是单个字节。然而，像 "ՙ" 这样的字符可能需要多个字节。因此，逐字节反转字符串将使多字节字符无效
		t.Logf("Number of runes: orig=%d, rev=%d, doubleRev=%d", utf8.RuneCountInString(orig), utf8.RuneCountInString(rev), utf8.RuneCountInString(doubleRev))
		// TODO 功能正确性：检查两次反转是否恢复原字符串，如果 orig 和 doubleRev 不相等，说明 Reverse 函数存在逻辑错误。使用 t.Errorf 报告错误，并打印原始字符串和最终结果
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}
		// TODO 编码正确性：检查 UTF-8 编码有效性，如果 orig 是有效的 UTF-8 字符串，但 rev 不是，则说明 Reverse 函数破坏了 UTF-8 编码规则。使用 utf8.ValidString 检查字符串是否符合 UTF-8 编码
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q, string %q", orig, rev)
		}
	})
}
