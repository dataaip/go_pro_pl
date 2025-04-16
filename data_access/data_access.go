package data_access

/*
TODO 数据库访问
1、设置数据库，创建要使用的数据库。使用数据库管理系统本身的命令行界面（CLI）来创建数据库和表，以及添加数据，使用 MySQL CLI，但大多数 DBMS 都有自己的 CLI，使用 data-tables.sql 的语句创建
2、查找并导入数据库驱动程序，找到并导入一个数据库驱动程序，它将把你通过 database/sql 包中的函数发出的请求转换为数据库能够理解的请求，在你的浏览器中，访问 SQLDrivers 维基页面以确定你可以使用的驱动程序。使用页面上的列表确定你将使用的驱动程序。对于本教程中访问 MySQL，你将使用 Go-MySQL-Driver，注意驱动程序的包名称 —— 这里是 github.com/go-sql-driver/mysql
3、获取数据库句柄并连接 *sql.DB，编写一些 Go 代码，可以使用数据库句柄访问数据库，使用指向 sql.DB 结构体的指针，该结构体表示对特定数据库的访问
*/

/*
下一个主题：
查看数据访问指南，其中包括有关此处仅涉及的主题的更多信息
如果您是 Go 的新手，您可以在 Effective Go https://golang.google.cn/doc/effective_go 和如何编写 Go 代码 https://golang.google.cn/doc/code
Go Tour https://golang.google.cn/tour/ 是对 Go 基础的一步一步的介绍
*/

import (
	"database/sql"
	// TODO "database/sql" 是 Go 语言标准库中的一个包，专门用于与关系型数据库（如 MySQL、PostgreSQL、SQLite 等）进行交互。它提供了一组通用的接口和工具，用于执行 SQL 查询、管理数据库连接池以及处理查询结果
	// 数据库连接管理： 提供了一个线程安全的数据库连接池（*sql.DB），支持并发操作。自动管理数据库连接的打开、关闭、复用等操作
	// 执行 SQL 查询： 支持执行各种 SQL 操作，包括 查询单行数据：QueryRow 查询多行数据：Query 执行非查询语句（如 INSERT、UPDATE、DELETE）：Exec
	// 参数化查询： 使用占位符（如 ? 或 $1）来防止 SQL 注入攻击。不同数据库驱动可能使用不同的占位符格式
	// 错误处理：提供了详细的错误信息，例如：sql.ErrNoRows：表示查询没有找到任何记录。其他数据库相关的错误（如连接失败、语法错误等）
	// 事务支持：提供了对事务的支持，允许你通过 Begin 方法启动事务，并通过 Commit 或 Rollback 提交或回滚事务
	// *sql.DB、*sql.Rows、*sql.Row、*sql.Tx

	"fmt"
	"log"

	"os"
	// TODO "os" 是 Go 语言标准库中的一个包，提供了与操作系统交互的功能。它主要用于处理文件系统、环境变量、进程控制以及输入输出流等操作。"os" 包是 Go 编程中非常重要的工具，尤其是在需要与底层操作系统交互的场景下（如读写文件、管理路径、获取环境变量等）
	// 文件和目录操作：提供了对文件和目录的创建、删除、重命名、读写等操作。支持文件权限和元数据的管理
	// 环境变量：允许读取和设置环境变量，获取当前工作目录或用户主目录
	// 进程管理：启动外部命令或子进程，控制程序的退出状态
	// 输入/输出流：提供了标准输入（os.Stdin）、标准输出（os.Stdout）和标准错误输出（os.Stderr）的访问接口
	// 路径和文件信息：获取文件的元信息（如大小、修改时间等），处理路径相关的操作（如判断路径是否存在、是否为目录等）

	"github.com/go-sql-driver/mysql"
	// TODO MySQL 驱动程序 github.com/go-sql-driver/mysql
)

// TODO 4、*sql.DB 类型的 db 变量。这是您的数据库句柄，将db设为全局变量简化了这个示例。在实际生产环境中，你应避免使用全局变量，可以通过将变量传递给需要它的函数或者将其封装在结构体中来实现
var db *sql.DB

// 将 db 作为参数传递，明确函数的依赖关系，增强代码的可读性和可维护性，更容易进行单元测试，可以传入模拟的 *sql.DB 对象
// func queryUser(db *sql.DB, userID int) (string, error) {
//     var name string
//     err := db.QueryRow("SELECT name FROM users WHERE id = ?", userID).Scan(&name)
//     return name, err
// }
// 将 db 包装在结构体中，面向对象设计风格，便于扩展和组织代码，每个方法都可以通过结构体访问 db，减少重复代码
// type App struct {
//     DB *sql.DB
// }
// func (a *App) QueryUser(userID int) (string, error) {
//     var name string
//     err := a.DB.QueryRow("SELECT name FROM users WHERE id = ?", userID).Scan(&name)
//     return name, err
// }
// 普通函数：func QueryUser(app *App, userID int) (string, error)，在这种情况下，app 是一个显式的参数，调用时需要手动传入
// 方法：func (a *App) QueryUser(userID int) (string, error)，在这种情况下，a 是接收者，表示 QueryUser 方法属于 *App 类型
// 接收者：（receiver）是 Go 中方法的核心概念之一。它的作用是将方法绑定到某个类型上，使得该方法可以访问类型的字段或其他方法
// 指针接收者：func (a *App) MethodName(...) 可以修改接收者的字段。更适合涉及状态变化的场景
// 值接收者：func (a App) MethodName(...) 不会修改接收者的字段（因为是副本）。更适合只读操作

// 定一个 Album 结构体
type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func Data_access() {
	// TODO 5、使用 MySQL 驱动程序的Config—— 以及该类型的FormatDSN—— 来收集连接属性并将它们格式化为连接字符串的数据源名称。Config 结构体使得代码比连接字符串更容易阅读
	cfg := mysql.NewConfig()
	cfg.User = os.Getenv("DBUSER")
	cfg.Passwd = os.Getenv("DBPASS")
	// os.Getenv("DBUSER") 获取系统环境变量，运行前需设置 DBUSER=username，DBPASS=password
	cfg.User = "root"
	cfg.Passwd = "root1234"
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "recordings"

	// TODO 6、调用 sql.Open 来初始化 db 变量，传入 FormatDSN 的返回值
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	// 7、检查 sql.Open 是否有错误。例如，如果你的数据库连接细节格式不正确，它可能会失败
	if err != nil {
		// 为了简化代码，你正在调用 log.Fatal 来结束执行并将错误打印到控制台。在生产代码中，你会希望以更优雅的方式处理错误
		log.Fatal(err)
	}

	// TODO 8、调用 db.Ping 以确认连接数据库是否正常。在运行时，sql.Open 可能不会立即连接，具体取决于驱动程序。在这里，使用 Ping 来确认 database/sql 包在需要时能够连接
	pingErr := db.Ping()
	// 检查来自 Ping 的错误，以防连接失败。如果 Ping 成功连接，则打印一条消息
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	// TODO 9、查询多行数据，调用添加的 albumsByArtist 函数，使用 Query 查询，将其返回值赋给一个新的 albums 变量
	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("albums: %v !\n", albums)

	// TODO 10、查询单行数据，对于知道最多返回一行的 SQL 语句，可以使用 QueryRow，这比使用 Query 循环更简单
	album, err := albumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", album)

	// TODO 11、Add 数据 使用 Go 来执行一条 SQL 语句 ，以便向数据库添加一个新行，你已经看到了如何使用 Query 和 QueryRow 要执行不返回数据的 SQL 语句，你可以使用 Exec
	albID, err := addAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	// 使用一个新相册调用 addAlbum，将您要添加的相册的 ID 赋给一个 albID 变量
	fmt.Printf("ID of added album: %v\n", albID)
}

func albumsByArtist(name string) ([]Album, error) {
	// TODO 声明一个定义的 albums 切片，类型为 Album。这将保存从返回的行中获取的数据。结构体字段名称和类型与数据库列名称和类型相对应
	var albums []Album
	// TODO 使用 DB.Query 执行一个 SELECT 语句来查询具有指定艺术家名称的专辑，Query 的第一个参数是 SQL 语句。在该参数之后，您可以传入零个或多个任意类型的参数。这些参数为您提供了一个位置，用于在 SQL 语句中指定参数的值。通过将 SQL 语句与参数值分开（而不是像使用fmt.Sprintf那样将它们连接起来），您可以使database/sql包将值与 SQL 文本分开发送，从而消除任何 SQL 注入风险
	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	// TODO 延迟关闭rows，以便在函数退出时释放它所持有的任何资源，defer 是 Go 的一个关键字，用于延迟执行某段代码，直到当前函数返回时才执行，defer 语句会在函数返回前（return 之前）按后进先出（LIFO）的顺序执行。它常用于清理资源（如关闭文件、释放锁、关闭数据库连接等），避免忘记手动释放资源
	defer rows.Close()

	// TODO 循环遍历返回的行
	for rows.Next() {
		var alb Album
		// TODO 使用 Rows.Sca 将每行的列值赋给 Album 结构体字段，Scan 接受一个指向 Go 值的指针列表，其中列值将被写入。在这里，你传递指向使用 & 运算符创建的 alb 变量中的字段的指针。Scan 通过指针写入来更新结构体字段
		// 在循环内部，检查将列值扫描到结构体字段中时是否出现错误
		// TODO Go 的 if 语句支持以下形式 if initialization; condition { // 执行代码 } initialization：可以是一个短语句，用于定义变量（如 err := ...） condition：是一个布尔表达式，用于判断是否执行 if 块中的代码，这种写法的好处是，初始化的变量作用域仅限于 if 块内，避免污染外部作用域
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		// 在循环内部，将新的 alb 添加到 albums 切片中
		albums = append(albums, alb)
	}
	// TODO 循环结束后，使用 rows.Err 检查整个查询是否出现错误。请注意，如果查询本身失败，在此处检查错误是发现结果不完整的唯一方法
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

func albumByID(id int64) (Album, error) {
	// 声明 Album 类型用于接收查询结果
	var album Album
	// TODO 使用 DB.QueryRow 执行 SELECT 语句以查询具有指定 ID 的专辑，它返回一个sql.Row，为了简化调用代码（你的代码！），QueryRow 不返回错误。相反，它安排在稍后从Rows.Scan返回任何查询错误（例如sql.ErrNoRows）
	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)

	// TODO 使用 Row.Scan 将列值复制到结构字段中，并检查扫描中的错误
	// TODO if err := row.Scan(...); err != nil { ... } 是 Go 中一种常见的简洁写法，具有以下特点，简洁：将变量定义和条件判断合并在一行，作用域控制：变量仅在 if 块内有效，可读性高：对于熟悉 Go 的开发者来说，这种写法非常直观
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		// TODO 特殊错误 sql.ErrNotch 指示查询未返回任何行。通常，该错误值得用更具体的文本替换，例如此处的 no such album
		if err == sql.ErrNoRows {
			return album, fmt.Errorf("albumsById %d: no such album", id)
		}
		return album, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return album, nil
}

func addAlbum(alb Album) (int64, error) {
	// TODO 使用 DB.Exec 执行一条 SQL 语句 与 Query 类似，Exec 接受 SQL 语句，后跟 SQL 语句的参数值
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	// 检查尝试 INSERT 时是否出现错误
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	// TODO 使用 Result.LastInsertId 检索插入的数据库行的 ID
	id, err := result.LastInsertId()
	// 检查检索 ID 的尝试中是否有错误
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}
