package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// p.89 リスト3-9: multipart/form-dataを使ったファイル送信
func main() {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	writer.WriteField("title", "Real World HTTP 第2版")
	writer.WriteField("author", "Yoshiki Shibukawa")

	fileWriter, err := writer.CreateFormFile("thumbnail", "photo.jpg")
	if err != nil {
		panic(err)
	}

	readFile, err := os.Open("photo.jpg")
	if err != nil {
		panic(err)
	}
	defer readFile.Close()

	_, _ = io.Copy(fileWriter, readFile)
	_ = writer.Close()

	response, err := http.Post(
		"http://localhost:5000",
		writer.FormDataContentType(),
		&buffer,
	)
	if err != nil {
		panic(err)
	}

	log.Println("Status:", response.Status)
}
