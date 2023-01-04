package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path"
)

// uploaderHandler expects two fields to be posted, userid and avatarFile.
func uploaderHandler(w http.ResponseWriter, req *http.Request) {
	userID := req.FormValue("userid")
	// file: io.Reader 型であると共に multipart.File を実装している。
	// header: multipart.FileHeader でファイルに関するメタデータを持っている
	file, header, err := req.FormFile("avatarFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	filename := path.Join("avatars", userID+path.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, "Successful")
}
