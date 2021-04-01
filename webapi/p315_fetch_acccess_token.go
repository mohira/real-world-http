package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

var clientID = os.Getenv("CLIENT_ID")
var clientSecret = os.Getenv("CLIENT_SECRET")
var redirectURL = "http://localhost:18888"
var state = "your state" // OAuth2.0 State の役割  https://qiita.com/naoya_matsuda/items/67a5a0fb4f50ac1e30c1

// リスト10-1: アクセストークンを取得する
// [OAuth 2.0 全フローの図解と動画 - Qiita](https://qiita.com/TakahikoKawasaki/items/200951e5b5929f840a1f)
// ここを追っていけば RFC6749 の Authorization Gode Grant がわかる！
// 	1. 最初の図解をみる
// 	2. Oauthライブラリ(Real-World-HTTPのサンプルコードでおk)を使って流れをざっくり理解する
// 	3. デバッグ実行して、どの情報が(ex: URLとか認可コード)どう関連しているかを図と対応させる
// 4. リクエストとレスポンスはデバッグ実行して、リクエストする直前のRequestオブジェクトの中をみればいける
// 重要な感想: Goのコードだと、ライブラリ中身も読みやすい感じがある
// 	io.Reader/io.Writer と net/http まわりの書き方を知っていたから、全然読めた
//  もちろん、HTTPやOAuth2の知識も大事だけどね。

func main() {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user:email", "gist"},
		Endpoint:     github.Endpoint,
	}
	var token *oauth2.Token // これは最終的なアクセストークンか？

	file, err := os.Open("./access_token.json")

	if os.IsNotExist(err) {
		// 初回アクセス
		// まず、認可画面のURLを取得する
		// 「認可エンドポイント」 https://qiita.com/TakahikoKawasaki/items/200951e5b5929f840a1f#11-%E8%AA%8D%E5%8F%AF%E3%82%A8%E3%83%B3%E3%83%89%E3%83%9D%E3%82%A4%E3%83%B3%E3%83%88%E3%81%B8%E3%81%AE%E3%83%AA%E3%82%AF%E3%82%A8%E3%82%B9%E3%83%88
		// ex: https://github.com/login/oauth/authorize?access_type=online&client_id=1d502a402511caeabe40&response_type=code&scope=user%3Aemail+gist&state=your+state
		url := conf.AuthCodeURL(state, oauth2.AccessTypeOnline)
		fmt.Println(url)

		// コールバックを受け取るWebサーバーを用意する
		code := make(chan string) // アクセストークンの引換券コードを格納するチャンネル

		var server *http.Server
		server = &http.Server{
			Addr: "127.0.0.1:18888", // 参考: https://qiita.com/terabyte/items/ec4c58bec9425fc02e86
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// 認可画面でOKを押すと、設定済みのリダイレクトURI(つまり、http://localhost:18888 にリダイレクトする)
				// リダイレクトURLでのサーバーを自分で立てておいて、そこから認可コードを取得する作戦

				// ブラウザ閉じたほうが行儀がいいけど、わかりにくいので閉じないことにした
				//_, _ = io.WriteString(w, "<html><script>window.open('about:blank', '_self').close()</script></html>")
				w.Header().Set("Content-Type", "text/html")
				message := fmt.Sprintf("URL: %s\n", r.URL.String())
				_, _ = io.WriteString(w, message)
				w.(http.Flusher).Flush()

				// 引き換えコード(認可コード)は、リダイレクトURLのクエリパラメータに仕込まれている！
				// ex: http://localhost
				// :18888/?code=79a361a22491547f21bc&state=your+state
				code <- r.URL.Query().Get("code")

				// 引き換えコードがほしいだけ(リダイレクトの受け皿があればいいだけ)なので、サーバーはすぐおさらばする
				_ = server.Shutdown(context.Background())
			}),
		}
		go server.ListenAndServe() // サーバー起動

		// ブラウザで認可画面を開く。ここで人間がアプリに対して許可を与える感じ
		_ = open.Start(url)

		// 手に入れた引き換えコードを、最終的なアクセストークンと交換する
		// 1.3. トークンエンドポイントへのリクエスト
		// https://qiita.com/TakahikoKawasaki/items/200951e5b5929f840a1f#13-%E3%83%88%E3%83%BC%E3%82%AF%E3%83%B3%E3%82%A8%E3%83%B3%E3%83%89%E3%83%9D%E3%82%A4%E3%83%B3%E3%83%88%E3%81%B8%E3%81%AE%E3%83%AA%E3%82%AF%E3%82%A8%E3%82%B9%E3%83%88
		// MEMO: ここだけ自作できるんじゃね？ → curlコマンドからできた！
		/*
			こんな感じ
			```
			$  curl -X POST https://github.com/login/oauth/access_token \
			       -d 'code=452bbee01b0325cd0e4d&grant_type=authorization_code' \
			       -H "Content-Type: application/x-www-form-urlencoded" \
			       -H "Authorization: Basic {{認証情報}}"
			```

			2回めでやるとアウト！

			```
			$ curl -X POST https://github.com/login/oauth/access_token \
			       -d 'code=452bbee01b0325cd0e4d&grant_type=authorization_code' \
			       -H "Content-Type: application/x-www-form-urlencoded" \
			       -H "Authorization: Basic {{認証情報}}"
			error=bad_verification_code&error_description=The+code+passed+is+incorrect+or+expired.&error_uri=https%3A%2F%2Fdocs.github.com%2Fapps%2Fmanaging-oauth-apps%2Ftroubleshooting-oauth-app-access-token-request-errors%2F%23bad-verification-code%
			```
		*/
		token, err = conf.Exchange(oauth2.NoContext, <-code)
		if err != nil {
			panic(err)
		}

		// アクセストークンをファイルに保存する(これは単なるファイル書き込みなのでOAuth2は関係ない)
		file, err := os.Create("./access_token.json")
		if err != nil {
			panic(err)
		}
		_ = json.NewEncoder(file).Encode(token)

	} else if err == nil {
		// 一度認可をしてローカルに保存済みの場合は、JSONからアクセストークンを読みコム
		token = &oauth2.Token{}
		_ = json.NewDecoder(file).Decode(token)
	} else {
		panic(err)
	}

	// いろいろやる
	fmt.Println(token)

	// conf.TokenSource()関数を使うと、
	// リフレッシュトークンを使ってアクセストークンの再取得を自動化する。
	// GoogleのサービスもGitHubも、アクセストークンは1時間で期限が切れるが、 ← これまずくね？
	// ログイン状態をリフレッシュすることで、
	// 一度取得したアクセストークンが、
	// サーバー側で無効化するまで自由に使えるようになります。
	client := oauth2.NewClient(oauth2.NoContext, conf.TokenSource(oauth2.NoContext, token))

	// リスト10-2: メールアドレスを取得する
	fetchEmail(client)

	// gist投稿
	createNewGist(client)
}

// リスト10-2: メールアドレスを取得する
func fetchEmail(client *http.Client) {
	response, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	emails, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Email: %v", string(emails))
}

// リスト名なし: 新規Gistの作成
func createNewGist(client *http.Client) {
	gist := `{
		"description": "API Example",
		"public": false,
		"files": {
			"README.md": {
				"content": "Hello World"
			},
			"users.txt": {
				"content": "Bob,Tom,Ken"
			}
		}
	}`

	response, err := client.Post(
		"https://api.github.com/gists",
		"application/json",
		strings.NewReader(gist),
	)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 201 {
		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			panic(err)
		}
		log.Println(string(dump))
		os.Exit(1)
	}

	type GistResult struct {
		Url string `json:"html_url"`
	}
	gistResult := &GistResult{}

	err = json.NewDecoder(response.Body).Decode(&gistResult)
	if err != nil {
		panic(err)
	}

	// せっかくなので、作成したGistをブラウザで開く。このほうが分かりやすい。
	if gistResult.Url != "" {
		_ = open.Start(gistResult.Url)
	}

}
