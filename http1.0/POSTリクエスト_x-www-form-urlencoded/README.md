# 2.1 シンプルなフォームの送信(x-www-form-urlencoded)

## 参考

### [POST - HTTP | MDN](https://developer.mozilla.org/ja/docs/Web/HTTP/Methods/POST)

このへんの違いが書いてあるのでみるとGood

- `application/x-www-form-urlencoded`
- `multipart/form-data`

## 準備: サーバーの用意とPOSTフォームを用意

```python
from flask import Flask, render_template, request

app = Flask(__name__)


@app.route('/', methods=['GET', 'POST'])
def index():
    if request.method == 'POST':
        dump_request()

    return render_template('index.html')


def dump_request():
    content_type = request.headers.get('Content-Type')
    request_body = request.get_data()

    print(f'POST {request.path} {request.headers.environ.get("SERVER_PROTOCOL")}')
    print(f'Host: {request.host}')
    print(f'Content-Type: {content_type}')
    print()
    print(request_body.decode('utf-8'))


def main():
    app.run(debug=True)


if __name__ == '__main__':
    main()

```

```html
<!-- templates/index.html -->
<form action="" method="POST">
    Title: <input type="text" name="title"><br>
    Author: <input type="text" name="author"><br>
    <input type="submit">
</form>
```

## 実験

次の4つを比較する

1. `curl -d`
2. `curl --data-urlencode`
3. ブラウザのPOSTフォームからリクエスト
4. Goで実装したクライアント

### 1. `curl -d`

- `-d` は `--data`

```
$ man curl
       -d, --data <data>
              (HTTP) Sends the specified data in a POST request to the HTTP server, in the same way that  a
              browser  does when a user has filled in an HTML form and presses the submit button. This will
              cause curl to pass the data to the  server  using  the  content-type  application/x-www-form-
              urlencoded.  Compare to -F, --form.
```

```sh
$ curl --http1.0 \
    -d title="Real World HTTP 第2版" \
    -d author="Yoshiki Shibukawa" \
    http://localhost:5000
```

```
POST / HTTP/1.0
Host: localhost:5000
Content-Type: application/x-www-form-urlencoded

title=Real World HTTP 第2版&author=Yoshiki Shibukawa
```

### 2. `curl --data-urlencode`

- `URL-encoding`するかどうかが違うってちゃんと書いてあるよね

```
$ man curl
       --data-urlencode <data>
              (HTTP)  This posts data, similar to the other -d, --data options with the exception that this
              performs URL-encoding.
```

```sh
$ curl --http1.0 \
    --data-urlencode title="Real World HTTP 第2版" \
    --data-urlencode author="Yoshiki Shibukawa" \
    http://localhost:5000
```

```
POST / HTTP/1.0
Host: localhost:5000
Content-Type: application/x-www-form-urlencoded

title=Real%20World%20HTTP%20%E7%AC%AC2%E7%89%88&author=Yoshiki%20Shibukawa
```

### 3. ブラウザのPOSTフォームからリクエスト

```
POST / HTTP/1.1
Host: 127.0.0.1:5000
Content-Type: application/x-www-form-urlencoded

title=Real+World+HTTP+%E7%AC%AC2%E7%89%88&author=Yoshiki+Shibukawa
```

### 4. Goで実装したクライアント

```go
package main

import (
	"log"
	"net/http"
	"net/url"
)

// p.87 リスト3-6: x-www-form-urlencoded形式のフォームを送信する
func main() {
	values := url.Values{
		"title":  {"Real World HTTP 第2版"},
		"author": {"Yoshiki Shibukawa"},
	}

	response, err := http.PostForm("http://localhost:5000", values)
	if err != nil {
		panic(err)
	}

	log.Println("Status:", response.Status)
}
```

```
POST / HTTP/1.1
Host: localhost:5000
Content-Type: application/x-www-form-urlencoded

author=Yoshiki+Shibukawa&title=Real+World+HTTP+%E7%AC%AC2%E7%89%88
```