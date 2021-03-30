package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

// p.5 1.1.1 テストエコーサーバーの実行
func main() {
	var httpServer http.Server

	http.HandleFunc("/", handler)
	log.Println("Start http listening :18888")

	httpServer.Addr = ":18888"
	log.Println(httpServer.ListenAndServe())
}

func handler(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)

	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	fmt.Println(string(dump))
	fmt.Fprintf(w, "<html><body>hello</body></html>\n")
}
