package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

// リスト6-17: サーバーからチャンク形式で送信
func main() {
	httpServer := http.Server{Addr: "localhost:18888"}
	log.Printf("Start http listening %s\n", httpServer.Addr)

	http.HandleFunc("/chunked", handlerChunkedResponse)

	log.Println(httpServer.ListenAndServe())
}

// $ curl http://localhost:18888/chunked
// $ http http://localhost:18888/chunked # 一気に読み込むっぽい
func handlerChunkedResponse(w http.ResponseWriter, r *http.Request) {
	// 検証に便利なのでリクエストをサーバーのログとして表示しておく
	dumpRequestMessage(r)

	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("expected http.ResponseWriter to be an http.flusher")
	}

	for i := 1; i <= 10; i++ {
		// 現在時刻も表示するとわかりやすいと思う
		fmt.Fprintf(w, "%s: Chuncked #%d\n", time.Now().Format("2006-01-02T15:04:05.000"), i)

		flusher.Flush() // ここのFlush()をコメントアウトするとわかりやすい
		time.Sleep(500 * time.Millisecond)
	}

	flusher.Flush()
}

func dumpRequestMessage(r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(dump))
}
