package vuln_tutorial

import (
	"fmt"
	"os"

	"golang.org/x/text/language"
	// TODO golang.org/x/text/language 是 Go 语言官方提供的一个子模块，属于 golang.org/x/text 包的一部分。它主要用于处理与语言相关的操作，例如语言标识、语言匹配、区域设置（locale）等。这个包的核心功能是帮助开发者在国际化（i18n）和本地化（l10n）场景中处理多语言支持
	// 语言标识符（Language Tags）提供了对 BCP 47 标准的支持，用于表示语言标签（Language Tag）。例如：en-US（美国英语）、zh-CN（简体中文，中国）、fr-FR（法国法语）。使用 language.Tag 类型来表示语言标签，并支持解析和标准化
	// 语言匹配（Language Matching）提供了语言匹配算法，可以根据用户偏好和可用的语言资源找到最佳匹配的语言。这对于实现动态语言切换或根据用户设置选择合适的语言非常有用
	// 语言覆盖范围（Coverage）提供了工具来检查某种语言是否被特定的资源或服务支持
	// 区域和脚本信息 支持语言的区域（Region）和脚本（Script）信息。例如，区分 zh-Hans（简体中文）和 zh-Hant（繁体中文）
	// language.Parse(string) (Tag, error)：解析语言标签
	// language.Make(string) Tag：快速创建语言标签
	// language.NewMatcher([]Tag) Matcher：创建语言匹配器
	// Matcher.Match(...Tag) (Tag, int, confidence)：匹配最佳语言
)

/*
TODO govulncheck 是 Go 官方提供的一款工具，用于检测代码中是否存在已知的安全漏洞。它通过分析你的项目依赖（即 go.mod 文件中列出的模块）和已发布的漏洞数据库（Go Vulnerability Database），来识别可能影响你的代码的已知安全问题

TODO Go Vulnerability Database 是一个由 Go 官方维护的开源漏洞数据库，记录了与 Go 模块相关的已知安全漏洞信息。这些漏洞通常包括：
- 第三方库中的安全问题
- 常见的编程错误或不安全的实现
- 其他可能导致程序被攻击的风险

TODO govulncheck 的主要功能是帮助开发者快速发现项目依赖中存在的安全漏洞。具体来说，它可以：
- 扫描依赖：检查你的项目依赖（go.mod 和 go.sum 文件中列出的模块）是否存在已知漏洞
- 定位问题：不仅报告哪些模块有漏洞，还会告诉你这些漏洞是否会影响你的代码路径
- 提供修复建议：在某些情况下，govulncheck 会提示如何升级依赖以修复漏洞

TODO govulncheck 的工作原理
- 解析依赖：读取你的项目的 go.mod 和 go.sum 文件，提取所有依赖模块及其版本
- 查询漏洞数据库：将这些依赖模块与 Go Vulnerability Database 中的漏洞记录进行匹配
- 分析调用链：通过静态分析技术，检查你的代码是否实际调用了存在漏洞的函数或方法
- 生成报告：输出漏洞信息，包括受影响的模块、漏洞描述以及修复建议

TODO 本教程将引导您完成以下步骤：
- 创建具有易受攻击依赖项的示例 Go 模块 Vuln_tutorial，将 golang.org/x/text 版本降级至 0.3.5 版，其中包含已知漏洞 go get golang.org/x/text@v0.3.5
- 安装并运行 govulncheck，使用 go install 命令安装 govulncheck ：go install golang.org/x/vuln/cmd/govulncheck@latest ，运行从要分析的文件夹（在本例中为 vuln_tutorial）中运行 $GOPATH/bin/govulncheck ./vuln_tutorial
```
govulncheck is an experimental tool. Share feedback at https://go.dev/s/govulncheck-feedback.

Using go1.20.3 and govulncheck@v0.0.0 with
vulnerability data from https://vuln.go.dev (last modified 2023-04-18 21:32:26 +0000 UTC).

Scanning your code and 46 packages across 1 dependent module for known vulnerabilities...
Your code is affected by 1 vulnerability from 1 module.

Vulnerability #1: GO-2021-0113
  Due to improper index calculation, an incorrectly formatted
  language tag can cause Parse to panic via an out of bounds read.
  If Parse is used to process untrusted user inputs, this may be
  used as a vector for a denial of service attack.

  More info: https://pkg.go.dev/vuln/GO-2021-0113

  Module: golang.org/x/text
    Found in: golang.org/x/text@v0.3.5
    Fixed in: golang.org/x/text@v0.3.7

    Call stacks in your code:
      main.go:12:29: vuln.tutorial.main calls golang.org/x/text/language.Parse

=== Informational ===

Found 1 vulnerability in packages that you import, but there are no call
stacks leading to the use of this vulnerability. You may not need to
take any action. See https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck
for details.

Vulnerability #1: GO-2022-1059
  An attacker may cause a denial of service by crafting an
  Accept-Language header which ParseAcceptLanguage will take
  significant time to parse.
  More info: https://pkg.go.dev/vuln/GO-2022-1059
  Found in: golang.org/x/text@v0.3.5
  Fixed in: golang.org/x/text@v0.3.8
```
代码受到一个漏洞 GO-2021-0113，因为它在易受攻击的版本 （v0.3.5） 直接调用 golang.org/x/text/language 的 Parse 函数，另一个漏洞 GO-2022-1059 存在于 v0.3.5 的 golang.org/x/text 模块中。但是，它被报告为 “Informational”，因为我们的代码从未（直接或间接）调用其任何易受攻击的函数

- 评估漏洞
TODO 1）评估漏洞：首先，阅读漏洞描述，并确定它是否真的适用于您的代码和使用场景。如果需要更多信息，请访问“更多信息”链接，根据描述，漏洞 GO-2021-0113 在使用 Parse 处理不受信任的用户输入时可能引发崩溃。假设我们希望我们的程序能够承受不受信任的输入，并且我们担心拒绝服务问题，那么该漏洞很可能适用，GO-2022-1059 可能不会影响我们的代码，因为我们的代码不会从该报告中调用任何易受攻击的函数
TODO 2）决定采取行动：为了缓解 GO-2021-0113 漏洞，我们有以下几个选项：
- 选项1：升级到修复版本。如果已经有修复可用，我们可以通过升级到模块的修复版本来移除易受攻击的依赖
- 选项2：停止使用易受攻击的符号。我们可以选择从代码中移除对易受攻击函数的所有调用。我们需要找到替代方案或自行实现
在本例中，修复已经可用，并且 Parse 函数对我们程序至关重要。因此，我们将依赖升级到 已修复 版本，即 v0.3.7
我们决定优先修复信息级别的漏洞 GO-2022-1059，但由于它与 GO-2021-0113 处于同一个模块中，且其修复版本为 v0.3.8，因此我们可以通过升级到 v0.3.8 同时移除这两个漏洞

- 升级易受攻击的依赖项：幸运的是，升级易受攻击的依赖项相当简单
TODO golang.org/x/text 升级到 v0.3.8 运行 go get golang.org/x/text@v0.3.8
TODO 再次运行 govulncheck 检查 $GOPATH/bin/govulncheck ./vuln_tutorial
```
govulncheck is an experimental tool. Share feedback at https://go.dev/s/govulncheck-feedback.

Using go1.20.3 and govulncheck@v0.0.0 with
vulnerability data from https://vuln.go.dev (last modified 2023-04-06 19:19:26 +0000 UTC).

Scanning your code and 46 packages across 1 dependent module for known vulnerabilities...
No vulnerabilities found.
```

*/

// TODO Govulncheck 是一种低噪声工具，可帮助您查找和修复 Go 项目中易受攻击的依赖项。它通过扫描项目的依赖项以查找已知漏洞，然后识别代码中对这些漏洞的任何直接或间接调用来实现此目的
// TODO 在本教程中，将学习如何使用 govulncheck 扫描简单的程序以查找漏洞。您还将学习如何确定漏洞的优先级和评估漏洞，以便您可以首先专注于修复最重要的漏洞
// TODO 要了解有关 govulncheck 的更多信息，请参阅 govulncheck 文档 https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck 以及这篇关于 Go 漏洞管理的博客文章
func Vuln_tutorial() {
	// TODO os.Args 是一个字符串切片，存储了程序运行时传递的命令行参数。os.Args[0] 是程序本身的路径或名称。os.Args[1:] 表示从索引 1 开始的所有参数（即用户输入的实际参数），使用 range 遍历每个参数，并将参数值赋给变量 arg
	for _, arg := range os.Args[1:] {
		// TODO 调用 language.Parse 方法解析字符串 arg，尝试将其转换为一个有效的语言标签（language.Tag 类型）
		// TODO 如果 language.Parse 返回的 err 不为 nil，说明解析失败
		tag, err := language.Parse(arg)
		if err != nil {
			fmt.Printf("%s: error: %v\n", arg, err)
			// TODO 检查是否为未定义标签 如果解析成功但返回的语言标签是 language.Und，说明该标签表示“未定义”（undefined）
		} else if tag == language.Und {
			fmt.Printf("%s: undefined\n", arg)
		} else {
			// TODO 输出有效语言标签，如果解析成功且返回的语言标签不是 language.Und，说明这是一个有效的语言标签
			fmt.Printf("%s: tag %s\n", arg, tag)
		}
	}
}
