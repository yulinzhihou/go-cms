package main

import "fmt"
import "net/http"

func handleFunc(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "请求路径为："+r.URL.Path+"\n")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.URL.Path == "/about" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "<h1>About</h1>")
	} else if r.URL.Path == "/login" {
		fmt.Fprint(w, "<h1>Login</h1>\n")
	} else if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello My Go Cms</h1>\n")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "<h1>页面未找到，请联系我们</h1>\n"+"<a href='mailto:yulinzhihou@163.com'>yulinzhihou@163.com</a>\n")
	}
}

func main() {
	http.HandleFunc("/", handleFunc)
	http.ListenAndServe(":8081", nil)
}
