package main

import (
	"log"
	"net/http"
)

// p.086 リスト3-5: HEADメソッドを送信してヘッダーを取得する
func main() {
	response, err := http.Head("http://localhost:18888")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	log.Println("Status:", response.Status)
	for k, v := range response.Header {
		log.Printf("%s: %v\n", k, v)
	}

	// HEADメソッドだからBodyは関係ないね
	// HEADの定義の通り、ボディを取得することはできません。
	// Get と同じように読み込んでみてもエラーにはなりませんが、
	// 長さゼロのバイト配列が返ってきます。
}
