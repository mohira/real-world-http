package main

import (
	"log"
	"net/http"
	"net/url"
)

// p.87 リスト3-6: x-www-form-urlencoded形式のフォームを送信する
func main() {
	values := url.Values{
		"name": {"Bob"},
		"age":  {"25"},
	}
	response, err := http.PostForm("http://localhost:18888", values)
	if err != nil {
		panic(err)
	}

	log.Println("Status:", response.Status)

}
