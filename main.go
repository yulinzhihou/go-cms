package main

import "fmt"
import "net/http"

func handleFunc(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "请求路径为："+r.URL.Path+"\n")
	if r.URL.Path == "/about" {
		fmt.Fprint(w, "<h1>About</h1>")
	} else if r.URL.Path == "/login" {
		fmt.Fprint(w, "<h1>Login</h1>\n")
	} else if r.URL.Path == "/" {
		fmt.Fprint(w, "<h1>Hello My Go Cms</h1>\n")
	} else {
		fmt.Fprint(w, "<h1>页面未找到，请联系我们not found</h1>\n")
	}
}

func main() {
	http.HandleFunc("/", handleFunc)
	http.ListenAndServe(":8081", nil)
}
