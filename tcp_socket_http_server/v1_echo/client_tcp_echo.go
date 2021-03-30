package main

import (
	"fmt"
	"io/ioutil"
	"net"
)

// 『Goならわかるシステムプログラミング』第4版 p.109 6.5.2 TCP ソケットを使ったHTTPクライアント
func main() {
	fmt.Println("Run main")
	conn, err := net.Dial("tcp_socket_http_server", "localhost:8888")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Connect %s\n", conn.RemoteAddr())
	requestMessage := "Hello from client_tcp_echo.go"
	if _, err := conn.Write([]byte(requestMessage)); err != nil {
		panic(err)
	}

	responseMessageBuffer, err := ioutil.ReadAll(conn)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(responseMessageBuffer))
}
