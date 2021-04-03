package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// 『Goならわかるシステムプログラミング』第4版 p.110 6.6.1 Keep-Alive対応のHTTP サーバー
func main() {
	fmt.Println("Run webapi")
	listener, err := net.Listen("tcp_socket_http_server", "localhost:8888")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Server is running at %s\n", listener.Addr())
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	// Accept後のソケットを使いまわして(==closeせずに)何度も応答を変えすためのforループ
	for {
		// タイムアウトの上限を設定
		_ = conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		request, err := http.ReadRequest(bufio.NewReader(conn))
		if err != nil {
			// タイムアウトもしくはソケットクローズ時は終了
			neterr, ok := err.(net.Error)
			if ok && neterr.Timeout() {
				fmt.Println("Timeout")
				break
			} else if err == io.EOF {
				break
			}
			panic(err)
		}

		// リクエストメッセージをサーバーのログに吐く(イメージ)
		dump, err := httputil.DumpRequest(request, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))

		// レスポンスメッセージの生成
		message := "Hello from server.go\n"
		response := http.Response{
			StatusCode:    200,
			ProtoMajor:    1,
			ProtoMinor:    1,
			ContentLength: int64(len(message)),
			Body:          ioutil.NopCloser(strings.NewReader(message)),
		}

		// レスポンス
		response.Write(conn)
	}
}
