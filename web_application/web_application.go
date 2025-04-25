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

TODO 通过闭包（makeHandler）和 模板引擎 的结合，代码实现了高度模块化的设计，便于扩展和维护
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
//
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
// 首先，我们必须将 html/template 添加到导入列表中 import "html/template"
// 1）创建一个包含 HTML 表单的模板文件，打开名为 edit.html 的新文件，并添加以下行
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
// - 方法 t.Execute 执行模板，将生成的 HTML 写入到 http.ResponseWriter 中。.Title 和 .Body 这种带点的标识符分别指向 p.Title 和 p.Body
// - 模板指令包含在双花括号中 printf "%s" .Body 指令是一个函数调用，将 .Body 作为字符串输出，而不是字节流，这与调用 fmt.Printf 相同。html/template 包有助于确保模板操作仅生成安全且外观正确的 HTML。例如，它会自动转义任何大于号（>），将其替换为 &gt;，以确保用户数据不会破坏表单 HTML
// 2）由于现在使用的是模板，那为 viewHandler 也创建一个名为 view.html 的模板
// <h1>{{.Title}}</h1>
// <p>[<a href="/edit/{{.Title}}">edit</a>]</p>
// <div>{{printf "%s" .Body}}</div>
// - 相应地修改 viewHandler
// func viewHandler(w http.ResponseWriter, r *http.Request) {
//     title := r.URL.Path[len("/view/"):]
//     p, _ := loadPage(title)
//     t, _ := template.ParseFiles("view.html")
//     t.Execute(w, p)
// }
// 请注意，我们在两个处理程序中使用了几乎完全相同的模板代码。可以通过将模板代码独立出来一个模版函数中来消除这种重复
// func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
//     t, _ := template.ParseFiles(tmpl + ".html")
//     t.Execute(w, p)
// }
// - 修改处理程序以使用该函数
// func viewHandler(w http.ResponseWriter, r *http.Request) {
//     title := r.URL.Path[len("/view/"):]
//     p, _ := loadPage(title)
//     renderTemplate(w, "view", p)
// }
// func editHandler(w http.ResponseWriter, r *http.Request) {
//     title := r.URL.Path[len("/edit/"):]
//     p, err := loadPage(title)
//     if err != nil {
//         p = &Page{Title: title}
//     }
//     renderTemplate(w, "edit", p)
// }

// TODO 处理不存在的页面：如果你访问/view/APageThatDoesntExist会怎样？你会看到一个包含 HTML 的页面。这是因为它忽略了loadPage的错误返回值，并继续尝试在没有数据的情况下填充模板。相反，如果请求的页面不存在，它应该将客户端重定向到编辑页面，以便可以创建内容
// func viewHandler(w http.ResponseWriter, r *http.Request) {
//     title := r.URL.Path[len("/view/"):]
//     p, err := loadPage(title)
//     if err != nil {
//         http.Redirect(w, r, "/edit/"+title, http.StatusFound)
//         return
//     }
//     renderTemplate(w, "view", p)
// }
// - http.Redirect 函数会向 HTTP 响应添加一个 HTTP 状态码http.StatusFound（302）以及一个 Location 标头
//
// TODO 保存页面：函数 saveHandler 将处理位于编辑页面上的表单提交。它将提取表单数据并将其保存到文件中。saveHandler 还会将用户重定向到查看页面
// func saveHandler(w http.ResponseWriter, r *http.Request) {
//     title := r.URL.Path[len("/save/"):]
//     body := r.FormValue("body")
//     p := &Page{Title: title, Body: []byte(body)}
//     p.save()
//     http.Redirect(w, r, "/view/"+title, http.StatusFound)
// }
// - 页面标题（在 URL 中提供）和表单的唯一字段 Body 会存储在一个新的 Page 中。然后调用 save() 方法将数据写入文件，并且客户端会被重定向到 /view/ 页面
// - FormValue 返回的值类型为 string 在将其放入 Page 结构体之前，必须将该值转换为[]byte 使用 []byte(body) 进行转换
//
// TODO 错误处理 程序中有几个地方忽略了错误。这是一种不好的做法，尤其是因为当错误真的发生时，程序会出现意外行为。更好的解决方案是处理这些错误，并向用户返回一条错误消息。这样一来，如果真的出现问题，服务器将按我们期望的方式运行，并且可以通知用户
// - 处理 renderTemplate 中的错误
// func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
//     t, err := template.ParseFiles(tmpl + ".html")
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     err = t.Execute(w, p)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//     }
// }
// http.Error 函数发送指定的 HTTP 响应代码（在这种情况下为 “内部服务器错误”）和错误消息。将此操作放在一个单独的函数 renderTemplate 中的决定已经开始奏效
// - 修复 saveHandler
// func saveHandler(w http.ResponseWriter, r *http.Request) {
//     title := r.URL.Path[len("/save/"):]
//     body := r.FormValue("body")
//     p := &Page{Title: title, Body: []byte(body)}
//     err := p.save()
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     http.Redirect(w, r, "/view/"+title, http.StatusFound)
// }
// p.save() 过程中发生的任何错误都将报告给用户

// TODO 模板缓存 这段代码存在一个效率问题：renderTemplate 在每次渲染页面时都会调用 ParseFiles 更好的方法是在程序初始化时调用一次 ParseFiles 将所有模板解析为单个 *Template 然后，可以使用 ExecuteTemplate 方法来渲染特定的模板
// - 创建一个名为templates的全局变量，并用ParseFiles对其进行初始化
// var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
// 函数 template.Must 是一个便捷包装器，当传入非空的 error 值时会引发 panic 否则将原封不动地返回 *Template 在这里引发 panic 是合适的；如果无法加载模板，唯一明智的做法就是退出程序
// ParseFiles 函数接受任意数量的字符串参数，这些参数用于指定模板文件，并将这些文件解析为以基础文件名命名的模板。如果要向程序中添加更多模板，就将其名称添加到 ParseFiles 调用的参数中
// - 修改 renderTemplate 函数，以便使用合适的模板名称调用 templates.ExecuteTemplate 方法
// func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
//     err := templates.ExecuteTemplate(w, tmpl+".html", p)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//     }
// }
// 模板名称是模板文件名，因此必须在 tmpl 参数后附加 .html

// TODO 验证 URL 路径 这个程序存在一个严重的安全漏洞：用户可以提供一个任意路径，在服务器上进行读取 / 写入操作。为了缓解这个问题，可以编写一个函数，使用正则表达式来验证
// - 导入 regexp ，然后可以创建一个全局变量来存储我们的验证表达式
// var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
// 函数 regexp.MustCompile 将解析并编译正则表达式，并返回一个 regexp.Regexp，MustCompile 与 Compile 的不同之处在于，如果表达式编译失败，它会引发 panic 而 Compile 会将 error 作为第二个参数返回
// - 现在，编写一个函数，该函数使用 validPath 表达式来验证路径并提取页面标题，如果标题有效，则返回该标题以及 nil 错误值。如果标题无效，该函数会将“404 Not Found”错误写入 HTTP 连接，并将错误返回给处理程序。要创建新错误，必须导入 errors 包
// func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
//     m := validPath.FindStringSubmatch(r.URL.Path)
//     if m == nil {
//         http.NotFound(w, r)
//         return "", errors.New("invalid Page Title")
//     }
//     return m[2], nil // The title is the second subexpression.
// }
// - 在每个处理程序中调用 getTitle
// func viewHandler(w http.ResponseWriter, r *http.Request) {
//     title, err := getTitle(w, r)
//     if err != nil {
//         return
//     }
//     p, err := loadPage(title)
//     if err != nil {
//         http.Redirect(w, r, "/edit/"+title, http.StatusFound)
//         return
//     }
//     renderTemplate(w, "view", p)
// }
// func editHandler(w http.ResponseWriter, r *http.Request) {
//     title, err := getTitle(w, r)
//     if err != nil {
//         return
//     }
//     p, err := loadPage(title)
//     if err != nil {
//         p = &Page{Title: title}
//     }
//     renderTemplate(w, "edit", p)
// }
// func saveHandler(w http.ResponseWriter, r *http.Request) {
//     title, err := getTitle(w, r)
//     if err != nil {
//         return
//     }
//     body := r.FormValue("body")
//     p := &Page{Title: title, Body: []byte(body)}
//     err = p.save()
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     http.Redirect(w, r, "/view/"+title, http.StatusFound)
// }

// TODO 介绍函数字面量和闭包
// 在每个处理程序中捕 获错误条件 会引入大量重复代码。要是能将每个处理程序包装在一个 执行此验证和错误检查 的函数中，情况会怎样呢？Go 语言的函数字面量提供了一种强大的抽象功能的方式，这在这种情况下能帮到我们
// - 函数字面量 在 Go 语言中，函数字面量（Function Literal） 是一种匿名函数的定义方式。它允许我们在代码中直接定义一个没有名字的函数，并且可以立即调用或将其赋值给变量、作为参数传递等，函数字面量是 Go 中闭包和高阶函数的基础，也是实现灵活编程的重要工具
// - 闭包：闭包（Closure） 是一种非常强大的特性。它允许函数捕获并访问其外部作用域中的变量，即使这些变量的生命周期已经超出了它们的原始作用域范围。闭包的本质是一个匿名函数（或普通函数）与其周围环境的结合
// 闭包可以分为两个部分：
// 函数本身：一个函数（可以是匿名函数或命名函数）
// 外部环境：函数运行时能够访问到的外部变量
// 当一个函数返回时，如果它内部定义的匿名函数仍然引用了外部变量，那么这些变量的值会被保留在内存中，供匿名函数使用。这就是闭包的核心机制
// - Go 语言中的闭包是一种强大且灵活的工具，适用于许多场景，如状态封装、回调函数、工厂函数和中间件等。理解闭包的工作原理及其注意事项，可以帮助你编写更高效、更优雅的代码。同时，合理使用闭包可以避免全局变量的滥用，提高代码的模块化和可维护性
//
// TODO 利用闭包处理 验证URL和错误
// 1）重写每个处理程序的函数定义以接受 title 字符串
// func viewHandler(w http.ResponseWriter, r *http.Request, title string)
// func editHandler(w http.ResponseWriter, r *http.Request, title string)
// func saveHandler(w http.ResponseWriter, r *http.Request, title string)
// 2）定义一个包装函数，接受上述类型的函数，并返回一个类型为 http.HandlerFunc 的函数（适合传递给函数http.HandleFunc）
// func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//         // 从 getTitle 获取代码并在此处使用它（有一些小的修改）
//         m := validPath.FindStringSubmatch(r.URL.Path)
//         if m == nil {
//             http.NotFound(w, r)
//             return
//         }
//         fn(w, r, m[2])
//     }
// }
// makeHandler 接收一个函数 返回一个函数，返回的函数使用了外部变量 fn
// 返回的函数被称为闭包，因为它封装了在其外部定义的值。在这种情况下，变量fn（makeHandler的单个参数）被闭包封装。变量fn将是我们的保存、编辑或查看处理程序之一
// makeHandler 返回的闭包是一个接受 http.ResponseWriter 和 http.Request（换句话说，一个 http.HandlerFunc）的函数。该闭包从请求路径中提取标题，并使用 validPath 正则表达式对其进行验证。如果标题无效，将使用 http.NotFound 函数向 ResponseWriter 写入错误。如果标题有效，将调用内部的处理器函数 fn，并将 ResponseWriter、Request 和标题作为参数传递给它
// 3）在函数中，将处理程序函数用 makeHandler 包装起来，然后再向 http 包注册
// func main() {
//     http.HandleFunc("/view/", makeHandler(viewHandler))
//     http.HandleFunc("/edit/", makeHandler(editHandler))
//     http.HandleFunc("/save/", makeHandler(saveHandler))
//     log.Fatal(http.ListenAndServe(":8080", nil))
// }
// 4）从处理程序函数中删除对 getTitle 的调用，使它们变得更加简单
// func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
//     p, err := loadPage(title)
//     if err != nil {
//         http.Redirect(w, r, "/edit/"+title, http.StatusFound)
//         return
//     }
//     renderTemplate(w, "view", p)
// }
// func editHandler(w http.ResponseWriter, r *http.Request, title string) {
//     p, err := loadPage(title)
//     if err != nil {
//         p = &Page{Title: title}
//     }
//     renderTemplate(w, "edit", p)
// }
// func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
//     body := r.FormValue("body")
//     p := &Page{Title: title, Body: []byte(body)}
//     err := p.save()
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }
//     http.Redirect(w, r, "/view/"+title, http.StatusFound)
// }
// 5）重新编译代码，然后运行应用程序：
// $ go build wiki.go
// $ ./wiki
// 访问 http://localhost:8080/view/ANewPage 时，应该会看到页面编辑表单。然后，应该能够输入一些文本，单击 保存 然后重定向到新创建的页面
// TODO 后续其他任务
// 将模板存储在 tmpl/ 中，将页面数据存储在 data/ 中
// 添加处理程序以使 Web 根重定向到 /view/FrontPage
// 通过使页面模板有效 HTML 并添加一些 CSS 规则来美化页面模板
// 通过转换 [PageName] 设置为 <a href="/view/PageName">PageName</a> 来实现页面间链接（提示：您可以使用 regexp，ReplaceAllFunc 来执行此作）

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
