package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var router = mux.NewRouter()
var db *sql.DB

func initDB() {
	var err error
	config := mysql.Config{
		User:                 "root",
		Passwd:               "Sinall0.123",
		Addr:                 "127.0.0.1:3306",
		Net:                  "tcp",
		DBName:               "go-cms",
		AllowNativePasswords: true,
	}
	// 准备数据库连接池
	db, err = sql.Open("mysql", config.FormatDSN())
	//fmt.Println(config.FormatDSN())
	checkError(err)
	// 设置最大连接数
	db.SetMaxOpenConns(25)
	// 设置最大空闲连接数
	db.SetMaxIdleConns(25)
	// 设置每个连接的过期时间
	db.SetConnMaxLifetime(8 * time.Minute)
	// 尝试连接，失败会报错
	err = db.Ping()
	checkError(err)
}

// 创建数据表
func createTables() {
	createArticlesSql := `CREATE TABLE IF NOT EXISTS articles (
    	id bigint(20) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    	title VARCHAR(255) NOT NULL,
    	body longtext NOT NULL COLLATE utf8mb4_general_ci NOT NULL
	);`

	_, err := db.Exec(createArticlesSql)
	checkError(err)
}

// 文章保存进数据库
func saveArticleToDB(title string, body string) (int64, error) {
	// 变量初始化
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)
	// 1，获取一个 prepare 声明语句
	stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES (?, ?)")
	// 2.例行的错误检测
	if err != nil {
		return 0, err
	}
	// 2,在此函数运行结束后关闭此语句，防止占用 SQL 连接
	defer stmt.Close()
	// 3.执行请求，传参进入绑定的内容
	rs, err = stmt.Exec(title, body)
	if err != nil {
		return 0, err
	}
	// 4.插入成功的话，会返回自增 ID
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}

	return 0, err
}

// 记录错误日志
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// 获取路由参数表
func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

// 使用 ID 获取文章
func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := `SELECT * FROM articles WHERE id=?`
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

// ArticlesFormData 创建文章表单数据
type ArticleFormData struct {
	Title  string
	Body   string
	URL    *url.URL
	Errors map[string]string
}

// 文章结构
type Article struct {
	ID    int64
	Title string
	Body  string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello! 欢迎来到 go-cms</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>about页</h1>"+"<a href='mailto:yulinzhihou@163.com'>mail</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>not found</h1>")
}

// 文章详情
func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	// 1。获取 URL 参数
	id := getRouteVariable("id", r)
	// 2。读取对应的文章数据
	article, err := getArticleByID(id)

	// 3.如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "<h1>文章未找到</h1>")
		} else {
			// 3.2 数据库有错误
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500服务器内部错误")
		}
	} else {
		// 4。读取成功
		tmpl, err := template.ParseFiles("resource/views/articles/show.gohtml")
		checkError(err)
		err = tmpl.Execute(w, article)
		checkError(err)
		fmt.Fprint(w, "读取成功，文章标题 -- "+article.Title)
	}
}

// 文章列表页
func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>文章列表页</h1>\n")
}

// 文章创建页
func articlesCreateHandler(w http.ResponseWriter, r *http.Request) {

	storeURL, _ := router.Get("articles.store").URL()
	data := ArticleFormData{
		Title:  "",
		Body:   "",
		URL:    storeURL,
		Errors: nil,
	}

	tmpl, err := template.ParseFiles("resource/views/articles/create.gohtml")
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

// 文章编辑页
func articlesEditHandler(w http.ResponseWriter, r *http.Request) {
	// 1.获取 URL 参数
	id := getRouteVariable("id", r)
	// 2.读取对应的文章数据
	article, err := getArticleByID(id)

	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 数据库错误
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 读取成功，显示表彰
		updateURL, _ := router.Get("articles.update").URL("id", id)
		data := ArticleFormData{
			Title:  article.Title,
			Body:   article.Body,
			URL:    updateURL,
			Errors: nil,
		}
		tmpl, err := template.ParseFiles("resource/views/articles/edit.gohtml")
		checkError(err)
		err = tmpl.Execute(w, data)
		checkError(err)

	}
}

// 文章更新
func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// 1.获取 URL 参数
	id := getRouteVariable("id", r)
	// 2。读取对应的文章数据
	_, err := getArticleByID(id)
	// 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 未出现错误
		// 表单验证
		title := r.FormValue("title")
		body := r.FormValue("body")
		// 验证字段
		errors := make(map[string]string)
		if title == "" {
			errors["title"] = "标题不能为空"
		} else if utf8.RuneCountInString(title) < 3 || utf8.RuneCountInString(title) > 10 {
			errors["title"] = "标题长度需介于 3 ~ 10 字符之间 "
		}

		if body == "" {
			errors["body"] = "内容不能为空"
		} else if utf8.RuneCountInString(body) < 10 || utf8.RuneCountInString(body) > 20 {
			errors["body"] = "内容长度需介于 10 ~ 20 字符之间"
		}

		// 验证内容
		if len(errors) == 0 {
			query := `UPDATE articles SET title=?, body=? WHERE id=?`
			rs, err := db.Exec(query, title, body, id)

			if err != nil {
				checkError(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "500 服务器内部错误")
			}

			// 更新成功。跳转到文章详情页面
			if n, _ := rs.RowsAffected(); n > 0 {
				showURL, _ := router.Get("articles.show").URL("id", id)
				http.Redirect(w, r, showURL.String(), http.StatusFound)
			} else {
				fmt.Fprint(w, "您没做任何更改")
			}

		} else {
			// 表单验证不通过
			updateURL, _ := router.Get("articles.update").URL("id", id)
			data := ArticleFormData{
				Title:  title,
				Body:   body,
				URL:    updateURL,
				Errors: errors,
			}

			tmpl, err := template.ParseFiles("resource/views/articles/edit.gohtml")
			checkError(err)
			err = tmpl.Execute(w, data)
			checkError(err)
		}
	}
}

// 文章删除
func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>文章删除</h1>\n")
}

// 文件保存
func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// 解析错误
		fmt.Fprint(w, "请提供正确的数据")
		return
	}
	// 定义错误日志
	errors := make(map[string]string)

	title := r.FormValue("title")
	body := r.FormValue("body")

	// 验证表单数据
	if title == "" {
		errors["title"] = "标题不能为空"
	} else if len(title) < 3 || len(title) > 100 {
		errors["title"] = "标题长度只能介于3-100个字符之间"
	}

	if body == "" {
		errors["body"] = "内容不能为空"
	} else if len(body) < 10 || len(body) > 255 {
		errors["body"] = "内容长度只能在10-255之间"
	}

	if len(errors) == 0 {
		fmt.Fprint(w, "验证成功！<br>")
		fmt.Fprintf(w, "title 的值为：%v <br>", title)
		fmt.Fprintf(w, "title 的长度为: %d <br> ", utf8.RuneCountInString(title))
		fmt.Fprintf(w, "body 的值为：%v <br> ", body)
		fmt.Fprintf(w, "body 的长度为：%d <br> ", utf8.RuneCountInString(body))
		lastInsertID, err := saveArticleToDB(title, body)
		if lastInsertID > 0 {
			fmt.Fprint(w, "插入成功，ID 为 "+strconv.FormatInt(lastInsertID, 10))
		} else {
			checkError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		fmt.Fprintf(w, "有错误发生：errors 的值为 %v<br>\n", errors)

		storeURL, _ := router.Get("articles.store").URL()

		data := ArticleFormData{
			Title:  title,
			Body:   body,
			URL:    storeURL,
			Errors: errors,
		}

		//tmpl, err := template.New("create-form").Parse(html)
		tmpl, err := template.ParseFiles("resource/views/articles/create.gohtml")

		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	}

}

// 中间件
func forceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 继续请求处理
		next.ServeHTTP(w, r)
	})
}

// 去除恪的中间件
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}

// 核心方法
func main() {
	initDB()
	createTables()

	router.StrictSlash(true)
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")

	router.HandleFunc("/articles/create", articlesCreateHandler).Methods("GET").Name("article.create")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles/{id:[0-9]+}/edit", articlesEditHandler).Methods("GET").Name("articles.edit")
	router.HandleFunc("/articles/{id:[0-9]+}", articlesUpdateHandler).Methods("POST").Name("articles.update")

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 通过命名路由来获取 URL 示例
	homeUrl, _ := router.Get("home").URL()
	fmt.Println("homeUrl = ", homeUrl)
	articleUrl, _ := router.Get("articles.show").URL("id", "123123")
	fmt.Println("articleUrl = ", articleUrl)

	// 使用中间件
	router.Use(forceMiddleware)

	http.ListenAndServe(":8081", removeTrailingSlash(router))
}
