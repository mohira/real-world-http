package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// p.85 リスト3-4: GETメソッドでクエリーを送信する
/*
$ curl localhost:18888 -G --data-urlencode "query=hello world"
$ curl localhost:18888 -G --data-urlencode "query=hello world&name=ぼぶ"
*/
func main() {
	// https://golang.org/pkg/net/url/#Values
	values := url.Values{
		"query": {"hello world"},
		"name":  {"ぼぶ"},
	}

	// https://golang.org/pkg/net/url/#Values.Encode
	response, err := http.Get("http://localhost:18888" + "?" + values.Encode())
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body))
}
