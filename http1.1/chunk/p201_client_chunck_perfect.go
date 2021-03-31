package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

// リスト6-19: 低レベルソケットでチャンクを直接読みこむ
func main() {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	conn, err := dialer.Dial("tcp", "localhost:18888")
	if err != nil {
		panic(err)
	}

	// サーバーへリクエストする
	request, err := http.NewRequest("GET", "http://localhost:18888/chunked", nil)
	if err != nil {
		panic(err)
	}
	if err := request.Write(conn); err != nil {
		panic(err)
	}

	// ???
	reader := bufio.NewReader(conn)
	response, err := http.ReadResponse(reader, request)
	if err != nil {
		panic(err)
	}

	if response.TransferEncoding[0] != "chunked" {
		panic("wrong transfer encoding")
	}

	for {
		// サイズを取得
		// MEMO: ReadBytes() はブロックする
		// MEMO: 区切り文字を 'X' などにすれば、ブロックされていることがわかる
		sizeStr, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		// 16進数のサイズをパース。サイズがゼロならクローズ
		hoge := string(sizeStr[:len(sizeStr)-2])
		size, err := strconv.ParseInt(hoge, 16, 64)
		if size == 0 {
			break
		}
		if err != nil {
			panic(err)
		}

		// サイズの数だけバッファを確保して読み込み
		line := make([]byte, int(size))
		_, _ = reader.Read(line)
		_, _ = reader.Discard(2)
		log.Println("  ", string(line))
	}

}
