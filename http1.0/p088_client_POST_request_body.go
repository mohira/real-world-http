package main

import (
	"log"
	"net/http"
	"os"
)

// p.88 リスト3-7: POSTメソッドで任意のボディを送信
func main() {
	file, err := os.Open("../README.md")
	if err != nil {
		panic(err)
	}

	response, err := http.Post(
		"http://localhost:18888",
		"text/plain",
		file,
	)
	if err != nil {
		panic(err)
	}

	log.Println("Status:", response.Status)
}
