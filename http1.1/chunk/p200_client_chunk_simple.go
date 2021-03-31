package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

// リスト6-18: クライアントでの逐次受信
// 区切り文字が明確ならこのくらいの記述でいける
func main() {
	response, err := http.Get("http://localhost:18888/chunked")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// わかりやすいのでヘッダーを表示
	dumpHeader(response, os.Stdout)

	reader := bufio.NewReader(response.Body)

	for {
		// MEMO: ReadBytes() はレスポンスが来るまでブロックする
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		log.Println(string(bytes.TrimSpace(line)))
	}

}

func dumpHeader(response *http.Response, w io.Writer) {
	dump, err := httputil.DumpResponse(response, false)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = fmt.Fprint(w, string(dump))
}
