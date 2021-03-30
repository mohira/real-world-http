package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

// p.81 リスト3-1: GETメソッドを送信して、レスポンスのボディを画面に出力する
func main() {
	response, err := http.Get("http://localhost:18888")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	log.Println(string(body))
}
