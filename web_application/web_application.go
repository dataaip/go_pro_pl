package web_application

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

// Page 结构体表示一个页面
type Page struct {
	Title string
	Body  []byte
}

// save 将页面保存到文件
func (p *Page) save() error {
	filename := p.Title + ".txt"
	log.Printf("save %s", filename)
	return os.WriteFile(filename, p.Body, 0600)
}

// loadPage 从文件中加载页面
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	log.Printf("load %s", filename)
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// viewHandler 处理查看页面的请求
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// editHandler 处理编辑页面的请求
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// saveHandler 处理保存页面的请求
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

// templates 用于存储解析后的模板
var templates = template.Must(template.ParseFiles("/Users/minghui.liu/vscode/go_pro/go_pro_pl/web_application/edit.html", "/Users/minghui.liu/vscode/go_pro/go_pro_pl/web_application/view.html"))

// renderTemplate 用于渲染模板
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// validPath 用于匹配 URL 路径
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// TODO 闭包
// makeHandler 用于创建处理函数
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

// Web_application 启动 web 应用
func Web_application() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
