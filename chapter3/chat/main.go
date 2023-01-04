package main

import (
	"flag"
	"github.com/stretchr/signature"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/matryer/goblueprints/chapter1/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// set the active Avatar implementation
// 1. ユーザがアップロードした画像を（あるか確認した後に）URL として取得する
// 2. 認証サービスに登録された画像を（あるか確認した後に）URL として取得する
// 3. Gravatar に登録された画像を（あるか確認した後に）URL として取得する
var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar}

// templ represents a single template
type templateHandler struct {
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if t.templ == nil {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	}

	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	t.templ.Execute(w, data)
}

var host = flag.String("host", ":8080", "The host of the application.")

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
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "", // クッキーを削除しないブラウザは空白で上書きする
			Path:   "/",
			MaxAge: -1, // ブラウザ上のクッキーは即座に削除される
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)

	http.Handle("/avatars/",
		http.StripPrefix("/avatars/", // http.Handle を受け取ってパスの接頭辞部分を削除する
			http.FileServer(http.Dir("./avatars")))) // 静的ファイルの提供やファイル一覧の作成、404 エラーの生成などの機能を提供する

	// get the room going
	go r.run()

	// start the web server
	log.Println("Starting web server on", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
