package web_application

/*
TODO 这段代码实现了一个简单的 Web 应用，具有以下功能：

1）页面管理：
- 查看页面（/view/）
- 编辑页面（/edit/）
- 保存页面（/save/）

2）路径验证：
- 使用正则表达式确保请求路径符合预期格式

3）模板渲染：
- 使用 HTML 模板动态生成页面内容

4）文件操作：
- 读取和保存页面内容文件

TODO 通过闭包（makeHandler）和模板引擎的结合，代码实现了高度模块化的设计，便于扩展和维护
TODO 设计 web 应用程序：go 开发后端 + vue.js 开发前端 前后端分离 + 多模块微服务架构

TODO 编写 Web 应用程序
- 创建数据结构
- 使用 load 和 save 方法操作数据结构
- 使用 net/http 包构建 Web 应用程序
- 使用 html/template 包处理 HTML 模板
- 使用 regexp 包验证用户输入
- 使用闭包
*/

import (
	"html/template"
	// 用于安全地渲染 HTML 模板

	"log"
	// 用于记录程序运行时的日志信息

	"net/http"
	// 提供 HTTP 协议支持，用于构建 Web 应用

	"os"
	// 用于文件系统操作（如读取和保存文件）

	"regexp"
	// 用于正则表达式匹配，验证 URL 路径
)

// 资源路径 定义了一个全局变量 resource_path，表示资源文件的存储路径。所有页面内容文件（.txt 文件）和 HTML 模板文件都存储在此路径下
var resource_path = "/Users/minghui.liu/vscode/go_pro/go_pro_pl/web_application/resource/"

// Page 结构体表示一个页面 定义了一个结构体 Page，用于表示一个页面的数据
// Title：页面的标题，通常是文件名的一部分
// Body：页面的内容，以字节数组的形式存储
// TODO 数据结构：从定义数据结构开始，Wiki 由一系列相互连接的页面组成，每个页面都有一个标题和一个正文（页面内容）。在这里，将 Page 定义为一个结构体，其中两个字段分别表示 title 和 body
// TODO 类型 []byte 表示 字节切片：（有关 slices 的更多信息，请参见 Slices： usage and internals https://golang.google.cn/doc/articles/slices_usage_and_internals.html Body 元素是一个 []字节，而不是 string 的 String，因为这是将使用的 IO 库所需的类型
// TODO Page 结构体描述了页面数据将如何存储在内存中
type Page struct {
	Title string
	Body  []byte
}

// save 将页面保存到文件 定义了一个方法 save，用于将页面内容保存到文件中
// 功能：
// - 构造文件名：filename 是基于 resource_path 和 Title 拼接而成的完整路径
// - 记录日志：使用 log.Printf 记录保存操作
// - 写入文件：调用 os.WriteFile 将 Body 写入文件，权限设置为 0600（仅当前用户可读写）
// 返回值：
// - 如果文件写入成功，返回 nil；否则返回错误
// TODO 但是持久存储呢？可以通过创建一个 save 方法，此方法的签名为：这是一个名为 save 的方法，采用指向 Page 的指针 p 作为其接收器。不需要任何参数，并返回 error 类型的值
// TODO 此方法会将 Page 的 Body 保存到文本文件中。为简单起见，将使用 Title 作为文件名
// TODO save 方法返回一个错误值，因为这是 WriteFile（将字节切片写入文件的标准库函数）的返回类型。save 方法返回 error 值，以便在写入文件时出现任何问题时让应用程序处理它。如果一切顺利，Page.save() 将返回 nil（指针、接口和一些其他类型的零值）
// TODO 八进制整数文本 0600，作为第三个参数传递给 WriteFile 表示应使用 仅当前用户的读写权限。（请参阅 Unix 手册页 open（2） 了解详情
func (p *Page) save() error {
	filename := resource_path + p.Title + ".txt"
	log.Printf("save %s", filename)
	return os.WriteFile(filename, p.Body, 0600)
}

// loadPage 从文件中加载页面 定义了一个函数 loadPage，用于从文件中加载页面内容
// 功能：
// - 构造文件名：基于 resource_path 和 Title 拼接出文件路径
// - 记录日志：使用 log.Printf 记录加载操作
// - 读取文件：调用 os.ReadFile 读取文件内容
// - 错误处理：如果文件不存在或读取失败，返回错误
// - 返回结果：如果成功，返回一个包含页面标题和内容的 Page 对象
// 返回值：
// - 成功时返回 (*Page, nil)
// - 失败时返回 (nil, error)
// TODO 除了保存页面之外，还需要加载页面，函数 loadPage 从 title 参数构造文件名，将文件的内容读取到新的变量 body 中，并返回指向使用正确 title 和 body 值构造的 Page 文本的指针
// TODO 函数可以返回多个值。标准库函数 操作系统。ReadFile 返回 []byte 和 error 在 loadPage 中，尚未处理错误;由下划线 （_） 符号表示的“空白标识符”用于丢弃错误返回值（实质上，将值分配给 nothing），但是，如果 ReadFile 遇到错误，会发生什么情况？例如，该文件可能不存在。不应该忽视这样的错误。修改函数以返回 *Page 和 error
// TODO 此函数的调用者现在可以检查第二个参数; 如果是 nil 则它已成功加载一个 Page 否则，它将是一个 错误 （请参阅 language specification 了解详情）https://golang.google.cn/ref/spec#Errors
func loadPage(title string) (*Page, error) {
	filename := resource_path + title + ".txt"
	log.Printf("load %s", filename)
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// TODO 此时，我们已有一个简单的数据结构，并且能够保存到文件并从文件加载。让我们编写一个 main 函数来测试我们编写的内容
// func main() {
//     p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
//     p1.save()
//     p2, _ := loadPage("TestPage")
//     fmt.Println(string(p2.Body))
// }
// go build wiki.go
// $ ./wiki
// This is a sample Page
// TODO 编译并执行此代码后，将创建一个名为 TestPage.txt 的文件，其中包含 p1 的内容。然后，该文件将被读入结构 p2 中，并将其 Body 元素打印到屏幕上

// TODO web_http_demo.go 一个简单 Web 服务器的完整工作示例

// TODO 使用 net/http 提供 wiki 页面
// TODO 1）创建一个处理程序 viewHandler，允许用户查看 wiki 页面。将处理前缀为 /view/ 的 URL，请注意使用 _ 来忽略 loadPage 的错误返回值 。为了简单起见，这里这样做通常被认为是不好的做法
// - 首先，此函数从 r.URL.Path（即请求 URL 的路径部分）中提取页面标题。Path 使用 [len("/view/"):] 重新切片，以去掉请求路径开头的 "/view/"，然后，该函数加载页面数据，使用一段简单的 HTML 字符串对页面进行格式化，并将其写入w，即http.ResponseWriter，要使用此处理程序，我们重写 main 函数，以使用 viewHandler 初始化 http，从而处理路径 /view/ 下的任何请求
// func viewHandler(w http.ResponseWriter, r *http.Request) {
//     title := r.URL.Path[len("/view/"):]
//     p, _ := loadPage(title)
//     fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
// }
// - 创建一些页面数据（如 test.txt），编译我们的代码，并尝试提供 wiki 页面
// - 在编辑器中打开 test.txt 文件，并在其中保存字符串 “Hello world” （不带引号），运行此 Web 服务器后，访问 http://localhost:8080/view/test 应显示标题为“test”的页面，其中包含单词“Hello world”
// func main() {
//     http.HandleFunc("/view/", viewHandler)
//     log.Fatal(http.ListenAndServe(":8080", nil))
// }

// TODO 2）编辑页面，Wiki 不是没有编辑页面能力的 Wiki。让我们创建两个新的处理程序：一个名为 editHandler 来显示 'edit page' 表单，另一个名为 saveHandler 来保存通过表单输入的数据
// 首先，我们将它们添加到 main（） 中：
// func main() {
//     http.HandleFunc("/view/", viewHandler)
//     http.HandleFunc("/edit/", editHandler)
//     http.HandleFunc("/save/", saveHandler)
//     log.Fatal(http.ListenAndServe(":8080", nil))
// }
// 函数 editHandler 加载页面（或者，如果它不存在，则创建一个空的 Page 结构），并显示一个 HTML 表单，这个函数可以正常工作，但所有硬编码的 HTML 都很丑陋
// func editHandler(w http.ResponseWriter, r *http.Request) {
//     title := r.URL.Path[len("/edit/"):]
//     p, err := loadPage(title)
//     if err != nil {
//         p = &Page{Title: title}
//     }
//     fmt.Fprintf(w, "<h1>Editing %s</h1>"+
//         "<form action=\"/save/%s\" method=\"POST\">"+
//         "<textarea name=\"body\">%s</textarea><br>"+
//         "<input type=\"submit\" value=\"Save\">"+
//         "</form>",
//         p.Title, p.Title, p.Body)
// }

// TODO html/template 包：html/template 包是 Go 标准库的一部分。可以使用 html/template 将 HTML 保存在一个单独的文件中，这样就可以在不修改底层 Go 代码的情况下更改编辑页面的布局
// - 首先，我们必须将 html/template 添加到导入列表中 import "html/template"
// - 创建一个包含 HTML 表单的模板文件，打开名为 edit.html 的新文件，并添加以下行
// <h1>Editing {{.Title}}</h1>
// <form action="/save/{{.Title}}" method="POST">
// 	<div><textarea name="body" rows="20" cols="80">{{printf "%s" .Body}}</textarea></div>
// 	<div><input type="submit" value="Save"></div>
// </form>
// - 修改 editHandler 以使用模板，而不是硬编码的 HTML
// func editHandler(w http.ResponseWriter, r *http.Request) {
//     title := r.URL.Path[len("/edit/"):]
//     p, err := loadPage(title)
//     if err != nil {
//         p = &Page{Title: title}
//     }
//     t, _ := template.ParseFiles("edit.html")
//     t.Execute(w, p)
// }
// - 函数 template.ParseFiles 将读取 edit.html 的内容并返回一个 *template.Template

// TODO

// viewHandler 处理查看页面的请求 定义了一个处理函数 viewHandler，用于处理查看页面的请求
// 功能：
// - 加载页面：调用 loadPage 加载指定标题的页面内容
// 错误处理：
// - 如果页面不存在（loadPage 返回错误），重定向到 /edit/ 页面进行编辑
// - 使用 http.Redirect 实现重定向，状态码为 302 Found
// - 渲染模板：如果页面存在，调用 renderTemplate 渲染 view 模板
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// editHandler 处理编辑页面的请求 定义了一个处理函数 editHandler，用于处理编辑页面的请求
// 功能：
// - 加载页面：调用 loadPage 加载指定标题的页面内容
// 错误处理：
// - 如果页面不存在（loadPage 返回错误），创建一个新的 Page 对象，仅包含标题
// - 渲染模板：调用 renderTemplate 渲染 edit 模板
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// saveHandler 处理保存页面的请求 定义了一个处理函数 saveHandler，用于处理保存页面的请求
// 功能：
// - 获取表单数据：从请求中提取表单字段 body 的值
// - 创建页面对象：构造一个 Page 对象，包含标题和内容
// - 保存页面：调用 save 方法将页面内容保存到文件中
// 错误处理：
// - 如果保存失败，返回 HTTP 500 错误
// - 重定向：保存成功后，重定向到 /view/ 页面查看内容
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// templates 用于存储解析后的模板 定义了一个全局变量 templates，用于存储解析后的 HTML 模板
// 使用 template.Must 确保模板解析成功，否则会触发 panic
// 模板文件包括：
// - edit.html：用于编辑页面
// - view.html：用于查看页面
var templates = template.Must(template.ParseFiles(resource_path+"edit.html", resource_path+"view.html"))

// renderTemplate 用于渲染模板 定义了一个函数 renderTemplate，用于渲染指定的 HTML 模板
// 功能：
// - 执行模板：调用 ExecuteTemplate 渲染模板，并将页面数据传递给模板
// - 错误处理：如果渲染失败，返回 HTTP 500 错误
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// validPath 用于匹配 URL 路径 定义了一个全局变量 validPath，用于验证 URL 路径是否符合预期格式
// 正则表达式含义：
// - 匹配以 /edit/、/save/ 或 /view/ 开头的路径，且后续部分必须由字母和数字组成
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// TODO 闭包
// makeHandler 用于创建处理函数 定义了一个高阶函数 makeHandler，用于创建统一的请求处理函数
// 功能：
// - 验证路径：使用 validPath 验证请求路径是否合法
// - 错误处理：如果路径不合法，返回 HTTP 404 错误
// - 调用目标函数：如果路径合法，提取动态参数（如页面标题），并调用目标处理函数 fn
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// Web_application 定义了一个函数 Web_application，用于启动 Web 应用
// 功能：
//   - 注册路由：
//     /view/：绑定到 makeHandler(viewHandler)
//     /edit/：绑定到 makeHandler(editHandler)
//     /save/：绑定到 makeHandler(saveHandler)
//   - 启动服务器：调用 http.ListenAndServe 启动 HTTP 服务器，监听端口 8080
//   - 错误处理：如果服务器启动失败，记录错误日志并终止程序
func Web_application() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
