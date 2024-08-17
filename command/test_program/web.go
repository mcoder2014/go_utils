package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 启动一个 ping 的 HTTP server
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	fmt.Printf("http server started")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
