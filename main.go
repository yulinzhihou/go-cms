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

// 记录错误日志
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// ArticlesFormData 创建文章表单数据
type ArticleFormData struct {
	Title  string
	Body   string
	URL    *url.URL
	Errors map[string]string
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
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprintf(w, "<h1>文章ID ： %s</h1>", id)
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

// 文章更新
func articlesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>文章更新</h1>\n")
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
