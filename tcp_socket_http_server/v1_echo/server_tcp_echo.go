package main

import (
	"fmt"
	"net"
)

// 『Goならわかるシステムプログラミング』第4版 p.107 6.5.1 TCP ソケットを使ったHTTPサーバー
func main() {
	fmt.Println("Run webapi")

	// net.Listen() は socket(2) -> bind(2) -> listen(2) まで一気にやっている
	// なので net.Bind() や net.Socket()という命令はない
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

		requestMessageBuffer := make([]byte, 1024)
		if _, err := conn.Read(requestMessageBuffer); err != nil {
			panic(err)
		}
		fmt.Println(string(requestMessageBuffer))

		responseMessage := "Hello from server_tcp_echo.go"
		if _, err := conn.Write([]byte(responseMessage)); err != nil {
			panic(err)
		}

		conn.Close()
	}

}
