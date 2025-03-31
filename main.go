package main

import (
	"fmt"
	"strings"
)
import "net/http"

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello My Go Cms</h1>\n")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>页面未找到，请联系我们</h1>\n"+"<a href='mailto:yulinzhihou@163.com'>yulinzhihou@163.com</a>\n")
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/about" {
		fmt.Fprint(w, "<h1>About</h1>\n")
	}
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", defaultHandler)
	router.HandleFunc("/about", aboutHandler)

	// 文章详情
	router.HandleFunc("/articles/", func(writer http.ResponseWriter, request *http.Request) {
		id := strings.SplitN(request.URL.Path, "/", 3)[2]
		fmt.Fprint(writer, "<h1>文章详情id：</h1>\n", id)
	})

	http.ListenAndServe(":8081", router)
}
