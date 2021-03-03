package main

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", defaultHttp)
	http.ListenAndServe(":8080", middlewareLimit(mux))
}

// 限流桶，每2s一个请求
var limiter = rate.NewLimiter(rate.Every(2*time.Second), 1)

func middlewareLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() == false {
			w.Write([]byte("so fast!!!"))
			fmt.Println("limit")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// 默认http处理
func defaultHttp(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "/" {
		w.Write([]byte("index"))
		fmt.Println("index")
		return
	}

	// 自定义404
	http.Error(w, "you lost???", http.StatusNotFound)
}
