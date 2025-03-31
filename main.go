package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"strings"
)
import "net/http"

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
	fmt.Fprint(w, "<h1>文章创建页</h1>\n")
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
	fmt.Fprint(w, "<h1>文件保存</h1>\n")
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

func main() {
	//router := http.NewServeMux()
	router := mux.NewRouter().StrictSlash(true)
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
