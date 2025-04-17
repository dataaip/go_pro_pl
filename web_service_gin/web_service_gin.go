package web_service_gin

/*
TODO REST API开发：本教程介绍了使用 Go 和 Gin Web Framework（Gin）编写 REST 风格 Web 服务 API 的基础知识
TODO Gin框架：Gin 简化了许多与构建 Web 应用程序（包括 Web 服务）相关的编码任务。在本教程中，您将使用 Gin 来路由请求、检索请求详细信息并为响应编组 JSON

TODO 在本教程中，你将构建一个具有两个端点的 RESTful API 服务器。你的示例项目将是一个关于复古爵士唱片数据的存储库
本教程包括以下部分：
- 设计 API 端点
- 为你的代码创建一个文件夹
- 创建数据
- 编写一个处理程序以返回所有项目
- 编写一个处理程序以添加新项目
- 编写一个处理程序以返回特定项目

TODO 1、设计 API 端点
你将构建一个提供对销售黑胶复古唱片商店访问权限的 API。因此，你需要提供端点，客户端可以通过这些端点为用户获取和添加专辑
在开发 API 时，通常首先设计端点。如果端点易于理解，你的 API 的用户将更容易取得成功。以下是你将在本教程中创建的端点
/albums
- GET – 获取所有专辑的列表，以 JSON 格式返回
- POST – 根据作为 JSON 发送的请求数据添加新专辑
/albums/:id
- GET – 通过 ID 获取一张专辑，并以 JSON 格式返回专辑数据

TODO 8、建议的下一个主题：
如果您是 Go 的新手，您可以在 有效的 Go https://golang.google.cn/doc/effective_go 和 如何写 Go 代码 https://golang.google.cn/doc/code
Go Tour https://golang.google.cn/tour/ 是对 Go 基础的一步一步的介绍。
有关 Gin 的更多信息，请参阅 Gin Web Framework 包文档 https://pkg.go.dev/github.com/gin-gonic/gin 和 Gin Web Framework 文档 https://gin-gonic.com/zh-cn/
*/

import (
	// TODO 6、导入需要支持刚刚编写的代码的包，并引入依赖 go get 或 go mod tidy
	"net/http"
	// TODO net/http 包： 是 Go 标准库中的一个核心包，用于处理 HTTP 协议相关的操作。它提供了构建 Web 服务器和客户端的基础功能
	// 创建 HTTP 服务器：提供了一个简单的接口来监听 HTTP 请求并返回响应。支持路由、中间件（通过自定义实现）、静态文件服务等
	// 发送 HTTP 请求：提供了发送 HTTP 请求的功能，例如 http.Get、http.Post 等
	// 处理请求和响应：定义了 http.Request 和 http.ResponseWriter 类型，用于处理客户端请求和生成响应
	// 支持 HTTPS：可以轻松配置 TLS/SSL，用于支持 HTTPS
	// 内置支持：无需额外安装第三方库
	// 灵活性：虽然功能基础，但可以通过自定义代码实现复杂的路由和中间件逻辑
	// 性能一般：对于小型项目足够使用，但在高并发场景下可能需要更高效的框架

	"github.com/gin-gonic/gin"
	// TODO gin 包： 是一个高性能的 HTTP Web 框架，基于 Go 的 net/http 包开发。它是目前 Go 社区中最流行的 Web 框架之一，广泛用于构建 RESTful API 和 Web 应用
	// TODO 主要功能
	// 路由管理：支持动态路由（如 /users/:id）和分组路由。内置对 HTTP 方法（如 GET、POST、PUT、DELETE 等）的支持
	// 中间件支持：提供了强大的中间件机制，便于实现日志记录、身份验证、跨域处理等功能
	// JSON 和数据绑定：自动将请求体解析为结构体（如 JSON、XML、Form 数据）。支持返回 JSON 响应
	// 错误处理：提供了统一的错误处理机制
	// 性能优越：使用 Radix 树算法优化路由匹配，性能远高于 net/http
)

// TODO 2、相册表示关于唱片专辑的数据
// 使用此结构体在内存中存储专辑数据，诸如json:"artist"这样的结构体标签指定了在将结构体的内容序列化为 JSON 时字段的名称应该是什么。如果没有它们，JSON 将使用结构体的大写字段名 —— 这种风格在 JSON 中并不常见
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// TODO 3、在刚刚添加的结构体声明下方，包含用于开始的数据的 album 结构体切片，专辑切片用来生成唱片专辑数据
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// TODO 4、编写一个处理程序来返回所有项目，当客户端在 GET /albums 上发出请求时，将以 JSON 的形式返回所有相册
// TODO 编写以下代码：准备答复的逻辑 getAlbums ，将请求路径 router 映射到 getAlbums 逻辑的代码
// TODO 编写一个 getAlbums 函数，该函数接收一个 gin.Context 参数。请注意，您可以为这个函数取任何名称 —— Gin 和 Go 都不要求特定的函数名称格式
// TODO gin.Context 参数：是 Gin 最重要的部分。它携带请求详细信息、验证和序列化 JSON 等。（尽管名称相似，但这与 Go 的内置context包不同）
func getAlbums(c *gin.Context) {
	// TODO 调用 Context.IndentedJSON 将结构体序列化为 JSON 并将其添加到响应中
	// TODO 该函数的第一个参数是你想要发送给客户端的 HTTP 状态码。在这里，你正在传递来自 net/http 包的 StatusOK 常量以表示 200 OK
	// 请注意，你可以用对 Context.JSON 的调用替换 Context.IndentedJSON 以发送更紧凑的 JSON 实际上，在调试时缩进形式更容易处理，并且大小差异通常很小
	c.IndentedJSON(http.StatusOK, albums)
}

// TODO 5、编写一个处理程序来添加新项，当客户端在 /albums 处发出 POST 请求时，你希望将请求体中描述的专辑添加到现有专辑数据中
// TODO 编写以下代码：将新专辑添加到现有列表的逻辑 postAlbums ，一段将 POST 请求路由 router 到逻辑 postAlbums 的代码
// 添加代码以将相册数据添加到相册列表中
func postAlbums(c *gin.Context) {
	// 创建 album 结构体 保存接收到的 JSON 数据
	var newAlbum album
	// 使用 Context.BindJSON 将请求正文 JSON 绑定到 newAlbum，调用 BindJSON 将接收到的 JSON 绑定到 NewAlbum
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	// 将新专辑添加到切片中，将从 JSON 初始化的 album 结构体追加到 albums 切片中
	albums = append(albums, newAlbum)
	// 向响应中添加状态码为201，同时添加表示你所添加专辑的 JSON
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// TODO 6、编写一个处理程序来返回特定的项，当客户端发出请求到 GET /albums/[id] 时，你需要返回 ID 与id 路径参数匹配的专辑
// TODO 编写以下代码：添加逻辑以检索请求的相册。将路径映射到逻辑
// 这个 getAlbumByID 函数将提取请求路径中的 ID，然后定位一个与之匹配的相册，getAlbumByID 用于查找其 ID 值与客户端发送的 id 参数相匹配的专辑，然后将该专辑作为响应返回
func getAlbumByID(c *gin.Context) {
	// TODO 使用 Context.Param 从 URL 中检索 id 路径参数。当将此处理程序映射到路径时，将在路径中为参数包含一个占位符
	id := c.Param("id")
	// 循环遍历专辑列表，寻找一个 ID 值与参数匹配的专辑
	// TODO 循环遍历切片中的 album 结构体，查找其 ID 字段值与 id 参数值匹配的结构体，如上文所述，在实际应用中，服务很可能会使用数据库查询来执行此查找
	for _, a := range albums {
		if a.ID == id {
			// 如果找到，则将该 album 结构体序列化为 JSON 并作为响应返回，同时返回 200 OK HTTP 状态码
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	// 如果未找到专辑，则返回带有 404 错误代码 和 http.StatusNotFound 的 HTTP 响应
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// TODO 7、将处理函数分配给端点路径，设置一个关联 处理函数 和 终结点路径，将路由器连接到 http 服务器并启动服务器
func Web_service_gin() {
	// TODO 初始化 Gin 路由器，使用 Default.
	router := gin.Default()
	// 这将设置一个关联，getAlbums 在该关联中处理对 /albums 终结点路径
	// TODO 使用 GET 函数将 GET HTTP 方法和 /albums 路径与处理程序函数 getAlbums 相关联，请注意，您正在传递 getAlbums 函数的名称 。这与传递函数的结果不同，后者可以通过传递 getAlbums（）（注意括号）来实现
	router.GET("/albums", getAlbums)
	// TODO 将路径 /albums 处的 POST 方法与 postAlbums 函数相关联
	// TODO 使用 Gin，你可以将处理程序与 HTTP 方法和路径组合相关联。通过这种方式，可以根据客户端使用的方法分别路由发送到单个路径的请求
	router.POST("/albums", postAlbums)
	// TODO 将路径 /albums/：id 路径与 getAlbumByID 函数关联。在 Gin 中，路径中项前面的冒号表示该项是路径参数
	router.GET("/albums/:id", getAlbumByID)

	// TODO 使用 Run 函数将路由器连接到 http 服务器并启动服务器
	router.Run("localhost:8080")

	// TODO 运行 go run . 一旦代码运行，就有了一个正在运行的 HTTP 服务器，可以向其发送请求
	// 使用 curl 向正在运行的 Web 服务发出请求 GET /albums curl http://localhost:8080/albums 或 curl http://localhost:8080/albums --header "Content-Type: application/json" --request "GET"
	// 使用 curl 向正在运行的 Web 服务发出请求 POST /albums curl http://localhost:8080/albums --include --header "Content-Type: application/json" --request "POST" --data '{"id": "4","title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}'
	// 使用 curl 向正在运行的 Web 服务发出请求 GET /albums/:id curl http://localhost:8080/albums/2
}
