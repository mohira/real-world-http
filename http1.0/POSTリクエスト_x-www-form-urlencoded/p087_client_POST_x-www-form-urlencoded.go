package main

import (
	"log"
	"net/http"
	"net/url"
)

// p.87 リスト3-6: x-www-form-urlencoded形式のフォームを送信する
func main() {
	values := url.Values{
		"title":  {"Real World HTTP 第2版"},
		"author": {"Yoshiki Shibukawa"},
	}

	response, err := http.PostForm("http://localhost:5000", values)
	if err != nil {
		panic(err)
	}

	log.Println("Status:", response.Status)
}
