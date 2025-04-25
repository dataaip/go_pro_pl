package greetings

/*
声明一个 greetings 包来收集相关函数
*/
import (
	"errors"
	// TODO 导入 Go 标准库 errors 包，以便您可以 使用其 errors.New 新功能

	"log"
	// TODO 配置log包，使其在日志消息开头打印命令名称（“greetings:”），但不包含时间戳或源文件信息

	"math/rand"
	// TODO 使用 math/rand 包 以生成用于从切片中选择项目的随机数

	"fmt"
)

// v1.0 实现一个 Hello 函数来返回问候语，此函数接受一个 name 参数，其类型为 string 。该函数还返回一个 string 。 在 Go 中，以大写字母开头的函数可以是 被不在同一包中的函数调用。这在 Go 中称为 导出的名称
// v2.0 更改函数，使其返回两个值：a string 和 error 。你的调用者将检查 第二个值来判断是否发生错误，任何 Go 函数都可以 返回多个值
func Hello(name string) (string, error) {
	// 添加一个 if 语句来检查请求是否无效（名称所在的位置为空字符串），如果请求无效则返回错误 errors.New 函数返回一个 您的消息里面 error
	if name == "" {
		return "", errors.New("empty name")
	}
	// 声明一个 message 变量来保存您的问候语，在 Go 语言中，:=运算符是在一行中声明并初始化变量的快捷方式。（Go 使用右侧的值来确定变量的类型）
	// TODO Sprintf：使用 fmt 包的 Sprintf 函数创建问候消息。第一个参数是格式字符串， Sprintf 会将 name 参数的值替换为 %v 格式动词。插入 name 参数的值即可完成 问候语，将格式化的问候语文本返回给呼叫者
	// 在 Hello中，调用 randomFormat 函数以获取你将返回的消息的格式，然后将格式和 name 值一起使用来创建消息
	message := fmt.Sprintf(randomFormat(), name)
	// message := fmt.Sprint(randomFormat())
	// TODO := 等价于
	// var message string
	// message = fmt.Sprintf("Hi, %v. Welcome!", name)
	// 添加 nil （表示没有错误）作为第二个值 成功返回。这样，调用者就可以看到函数 成功了
	return message, nil
}

// v3.0 添加一个randomFormat函数，该函数返回一个随机选择的问候消息格式。请注意，randomFormat以小写字母开头，这使得它只能被其所在包中的代码访问（换句话说，它没有被导出）
func randomFormat() string {
	// TODO 切片：在randomFormat中，声明一个formats切片，其中包含三种消息格式。在声明切片时，像这样在方括号中省略其大小：[]string。这告诉 Go，切片底层数组的大小可以动态更改
	format := []string{
		"Hi, %v. Welcome!",
		"Great to see you, %v!",
		"Hail, %v! Well met!",
	}
	// 使用math/rand包生成一个随机数，以便从切片中选择一项
	return format[rand.Intn(len(format))]
}

// v4.0 添加一个名为 “Hellos” 的函数，该函数的参数是一个名字切片，而不是单个名字。此外，将其返回类型之一从 “字符串” 更改为 “映射”，以便可以返回映射到问候消息的名字
// TODO 复用：让新的Hellos函数调用现有的Hello函数。这样有助于减少重复，同时让两个函数都保持在原位
func Hellos(names []string) (map[string]string, error) {
	// TODO map：创建一个messages映射，将每个接收到的名称（作为键）与生成的消息（作为值）相关联。在 Go 语言中，你可以使用以下语法初始化一个映射：make(map[key-type]value-type)。让Hellos函数将这个映射返回给调用者
	messages := make(map[string]string)
	// TODO range：遍历你的函数接收的名称，检查每个名称是否具有非空值，然后为每个名称关联一条消息。在这个 for 循环中，range 返回两个值：循环中当前项的索引和该项值的副本。你不需要索引，所以使用 Go 的空白标识符（下划线）来忽略它。
	for _, name := range names {
		message, err := Hello(name)
		if err != nil {
			return nil, err
		}
		// TODO 在映射中，将检索到的消息与名称关联起来
		messages[name] = message
	}
	return messages, nil
}

func Greetings() {
	// 配置log包，使其在日志消息开头打印命令名称（“greetings:”），但不包含时间戳或源文件信息
	log.SetPrefix("greetings: ")
	log.SetFlags(0)
	// v1.0 获取 greeting message 返回值并打印它.
	// message := greetings.Hello("Gladys")
	// v2.0 为 Hello 返回值赋值，包括 错误 ，变量，将 Hello 参数从 Gladys 的名字更改为空 字符串，以便您可以尝试错误处理代码，查找非 nil 错误值
	message, err := Hello("Gladys")
	if err != nil {
		// 使用标准库的日志包 输出错误信息。如果得到错误，则使用 日志包的 致命功能 打印错误并停止程序
		log.Fatal(err)
	}
	// 返回值并打印它
	fmt.Println(message)
	// 一组名字
	names := []string{"Gladys", "Samantha", "Darrin"}
	// 请求获取这些名称的问候消息
	messages, err := Hellos(names)
	if err != nil {
		log.Fatal(err)
	}
	// 如果未返回错误，则将返回的消息映射打印到控制台
	fmt.Println(messages)
}
