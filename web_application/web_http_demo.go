package web_application

import (
	"fmt"
	"log"
	"net/http"
)

// TODO 函数 handler 的类型为 http.HandlerFunc 将 http.ResponseWriter 和 http.Request 作为参数
// TODO 一个 http.ResponseWriter 值用于组合 HTTP 服务器的响应；通过向其写入数据，我们将数据发送到 HTTP 客户端
// TODO http.Request是一种数据结构，代表客户端的 HTTP 请求。r.URL.Path是请求 URL 的路径部分。末尾的[1:]表示 从第一个字符到末尾创建Path的子切片。这会从路径名称中去掉开头的 /
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

// TODO Web_http_demo 函数以调用 http.HandleFunc 开始，该调用告诉 http 包使用 handler 处理对网站根目录（"/"）的所有请求
func Web_http_demo() {
	http.HandleFunc("/", handler)
	// TODO 然后它调用 http.ListenAndServe，指定它应在任何接口上监听 8080 端口（":8080"）。（目前不用担心它的第二个参数 nil）这个函数将阻塞，直到程序终止
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 如果您运行此程序并访问 URL：http://localhost:8080/monkeys
// 该程序将显示一个页面，其中包含：Hi there, I love monkeys!
