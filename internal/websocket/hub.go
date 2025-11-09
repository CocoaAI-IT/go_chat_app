package websocket

import (
	"encoding/json"
	"log"
)

// Message はクライアント間でやり取りされるメッセージの構造
type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	Type     string `json:"type"` // "message", "join", "leave"
}

// Hub はアクティブなクライアントとメッセージのブロードキャストを管理
type Hub struct {
	// 登録されたクライアント
	clients map[*Client]bool

	// クライアントからのメッセージ
	broadcast chan []byte

	// クライアントの登録リクエスト
	register chan *Client

	// クライアントの登録解除リクエスト
	unregister chan *Client
}

// NewHub は新しいHubインスタンスを作成
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run はハブのメインループを実行
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("新しいクライアントが接続しました（合計: %d）", len(h.clients))

			// 参加メッセージの送信
			message := Message{
				Username: client.username,
				Content:  "チャットに参加しました",
				Type:     "join",
			}
			msgBytes, _ := json.Marshal(message)
			h.broadcast <- msgBytes

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("クライアントが切断しました（残り: %d）", len(h.clients))

				// 退室メッセージの送信
				message := Message{
					Username: client.username,
					Content:  "チャットから退出しました",
					Type:     "leave",
				}
				msgBytes, _ := json.Marshal(message)
				h.broadcast <- msgBytes
			}

		case message := <-h.broadcast:
			// 全クライアントにメッセージをブロードキャスト
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
