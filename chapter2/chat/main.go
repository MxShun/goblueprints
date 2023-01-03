package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/matryer/goblueprints/chapter1/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	// map[string]interface{} 型を意味する
	// interface{} はどんな型も格納できる特殊な型
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

var host = flag.String("host", ":8080", "The host of the application.")

// ユーザから見た処理の流れ
//  0. chat などにアクセスすると login にリダイレクトする。
//  1. auth/login から認証プラバイダにリダイレクトする。リダイレクト先の URL に含まれるクライアントIDが含まれる。
//     こんな感じ https://accounts.google.com/o/oauth2/auth/oauthchooseaccount&client_id=238242881503-hf7rp46ojbuhl732n94vj56n0pcnrabe.apps.googleusercontent.com&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fauth%2Fcallback%2Fgoogle?...
//  2. 認証プロバイダ画面で、認証プラバイダに対しチャットアプリケーションがアクセスを要求する。
//  3. 認証プロバイダにログインし、チャットアプリケーションのアクセスを許可します。
//  4. 認可コードデータと共に auth/callback にリダイレクトする
//  5. チャットアプリケーションから認証プロバイダに対し認可コードを送信し、アクセストークンを払い出します。
//  6. アクセストークンを含んだ認証済みのリクエストを行い、ユーザ情報（や現在の状況など）を取得します。
func main() {

	flag.Parse() // parse the flags

	// setup gomniauth
	// https://github.com/stretchr/gomniauth
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	gomniauth.WithProviders(
		//github.New("3d1e6ba69036e0624b61", "7e8938928d802e7582908a5eadaaaf22d64babf1", "http://localhost:8080/auth/callback/github"),
		google.New("238242881503-hf7rp46ojbuhl732n94vj56n0pcnrabe.apps.googleusercontent.com", "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", "http://localhost:8080/auth/callback/google"),
		//facebook.New("537611606322077", "f9f4d77b3d3f4f5775369f5c9f88f65e", "http://localhost:8080/auth/callback/facebook"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	// get the room going
	go r.run()

	// start the web server
	log.Println("Starting web server on", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
