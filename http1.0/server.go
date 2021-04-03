package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

// p.5 1.1.1 テストエコーサーバーの実行
// p.45 2.5 クッキー 初訪問かどうかを判定するハンドラ
func main() {
	var httpServer http.Server

	http.HandleFunc("/", echoHandler)
	http.HandleFunc("/cookie", cookieHandler)
	log.Println("Start http listening :18888")

	httpServer.Addr = ":18888"
	log.Println(httpServer.ListenAndServe())
}

func cookieHandler(w http.ResponseWriter, r *http.Request) {
	// 実験がわかりやすいように、リクエスト情報をdumpする
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(dump))

	w.Header().Add("Set-Cookie", "VISIT=TRUE")

	if _, ok := r.Header["Cookie"]; ok {
		fmt.Fprintf(w, "<html><body>2回目以降の訪問ですね！！</body></html>")
	} else {
		fmt.Fprintf(w, "<html><body>1回目の訪問かCookieがないですね！</body></html>")
	}
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	fmt.Println(string(dump))
	fmt.Fprintf(w, "<html><body>hello</body></html>\n")
}
