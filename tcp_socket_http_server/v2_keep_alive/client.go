package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
)

// 『Goならわかるシステムプログラミング』第4版 p.110 6.6.2 Keep-Alive対応のHTTPクライアント
func main() {
	fmt.Println("Run main")

	// 3つのリクエストを送るイメージ
	// わかりやすくするためにリクエストボディに載せる
	sendMessages := []string{
		"ASCII",
		"PROGRAMMING",
		"PLUS",
	}
	var current int
	var conn net.Conn = nil

	for {
		var err error

		// 接続してない場合は、connectする
		// Keep-Aliveに関係あるのはこの部分の記述
		if conn == nil {
			conn, err = net.Dial("tcp_socket_http_server", "localhost:8888")
			if err != nil {
				panic(err)
			}
			fmt.Printf("Connected %s Access #%d\n", conn.RemoteAddr(), current)
		}

		// リクエストボディでメッセージを表現する
		request, err := http.NewRequest(
			"POST",
			"http://localhost:8888",
			strings.NewReader(sendMessages[current]),
		)
		if err != nil {
			panic(err)
		}

		// リクエストする
		if err := request.Write(conn); err != nil {
			panic(err)
		}

		// レスポンスメッセージを読む
		response, err := http.ReadResponse(bufio.NewReader(conn), request)
		if err != nil {
			fmt.Println("Retry")
			conn = nil
			continue
		}

		// けっか
		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))

		// 3つのメッセージ(==3つのリクエスト)をすべて送ったら終了する
		current++
		if current == len(sendMessages) {
			break
		}
	}
	conn.Close()
}
