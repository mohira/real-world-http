package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
)

// p.92 リスト3-12: クッキーの送受信
func main() {
	// MEMO: デフォルトの実装では、クッキーはメモリ上にしか保存しない

	// クッキーを保存するcookiejar(クッキーの瓶)を作成
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	// http.Clientにjarを渡せばおk
	client := http.Client{Jar: jar}

	for i := 0; i < 2; i++ {
		response, err := client.Get("http://localhost:18888/cookie")
		if err != nil {
			panic(err)
		}

		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			panic(err)
		}

		log.Println(string(dump))
	}
}
