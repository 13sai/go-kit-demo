package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// 自定义返回
type JsonRes struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
	TimeStamp int64       `json:"timestmap"`
}

func apiResult(w http.ResponseWriter, code int, data interface{}, msg string) {
	body, _ := json.Marshal(JsonRes{
		Code: code,
		Data: data,
		Msg:  msg,
		// 获取时间戳
		TimeStamp: time.Now().Unix(),
	})
	w.Write(body)
}

func main() {
	srv := http.Server{
		Addr:    ":8080",
		Handler: http.TimeoutHandler(http.HandlerFunc(defaultHttp), 2*time.Second, "Timeout!!!"),
	}
	srv.ListenAndServe()
}

// 默认http处理
func defaultHttp(w http.ResponseWriter, r *http.Request) {
	i := rand.Intn(1000)
	fmt.Println(i)
	time.Sleep(time.Duration(i) * time.Millisecond)
	w.Write([]byte("index"))
	return
}
