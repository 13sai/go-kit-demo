package main

import (
	"fmt"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
)

func main() {
	http.HandleFunc("/", handle)
	fmt.Println("8081")
	http.ListenAndServe(":8081", nil)
}

// 默认http处理
func handle(w http.ResponseWriter, r *http.Request) {
	hystrix.ConfigureCommand("hystrix", hystrix.CommandConfig{
		Timeout:                500,
		RequestVolumeThreshold: 3,
		SleepWindow:            10000,
	})

	hystrix.Do("hystrix", func() error {
		// talk to other services
		_, err := http.Get("http://localhost:8080/")
		if err != nil {
			fmt.Println("get error")
			return err
		}
		w.Write([]byte("index"))
		return nil
	}, func(err error) error {
		w.Write([]byte("get an error, handle it"))
		return nil
	})
	return
}
