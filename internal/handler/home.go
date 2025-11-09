package handler

import (
	"html/template"
	"log"
	"net/http"
)

// HomeHandler はチャットアプリのホームページを表示
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("テンプレートの読み込みエラー: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("テンプレートの実行エラー: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
