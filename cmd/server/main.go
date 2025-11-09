package main

import (
	"log"
	"net/http"
	"os"

	"github.com/CocoaAI-IT/go_chat_app/internal/handler"
	"github.com/CocoaAI-IT/go_chat_app/internal/websocket"
	"github.com/gorilla/mux"
)

func main() {
	// WebSocketハブの作成と起動
	hub := websocket.NewHub()
	go hub.Run()

	// ルーターの設定
	router := mux.NewRouter()

	// 静的ファイルの提供
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// ルートの設定
	router.HandleFunc("/", handler.HomeHandler).Methods("GET")
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	// ポートの設定（環境変数から取得、デフォルトは8080）
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("サーバーを起動します: http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("サーバーの起動に失敗しました:", err)
	}
}
