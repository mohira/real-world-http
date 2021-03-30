package main

import (
	"log"
	"net/http"
	"strings"
)

// p.89 リスト3-8: 文字列を io.Reader インタフェース化する
func main() {
	reader := strings.NewReader("りくえすとぼでぃ")

	response, err := http.Post(
		"http://localhost:18888",
		"text/plain",
		reader,
	)

	if err != nil {
		panic(err)
	}

	log.Println("Status:", response.Status)
}
