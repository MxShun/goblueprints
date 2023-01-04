package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ErrNoAvatarURL is the error that is returned when the
// Avatar instance is unable to provide an avatar URL.
// このエラーオブジェクトが生成されるのは一度だけで、問題が発生した場合にはこのオブジェクトのポインタが渡される。
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents types capable of representing
// user profile pictures.
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client,
	// or returns an error if something goes wrong.
	// ErrNoAvatarURL is returned if the object is unable to get
	// a URL for the specified client.
	GetAvatarURL(ChatUser) (string, error)
}

type TryAvatars []Avatar

// GetAvatarURL chatUser ではなく ChatUser なのは各 GetAvatarURL が受け取った引数を柔軟に扱えるため。
// TryAvatars も Avatar の実装クラスになることで利用側が全部か別個かを呼び分けられる
func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

// GetAvatarURL ユーザがアップロードした画像を（あるか確認した後に）URL として取得する
func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	files, err := ioutil.ReadDir("avatars")
	if err != nil {
		return "", ErrNoAvatarURL
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fname := file.Name()
		if u.UniqueID() == strings.TrimSuffix(fname, filepath.Ext(fname)) {
			return "/avatars/" + fname, nil
		}
	}
	return "", ErrNoAvatarURL
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

// GetAvatarURL 認証サービスに登録された画像を（あるか確認した後に）URL として取得する
// レシーバオブジェクト AuthAvatar にフィールドがないのでレシーバを参照する必要がない
// func (_ AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {... も可
func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return u.AvatarURL(), nil
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}
