package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// クライアントへの書き込み待機時間
	writeWait = 10 * time.Second

	// クライアントからの次のpongメッセージ待機時間
	pongWait = 60 * time.Second

	// この期間内にクライアントにpingを送信（pongWaitより短くする必要がある）
	pingPeriod = (pongWait * 9) / 10

	// クライアントから許可される最大メッセージサイズ
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 本番環境では適切なオリジンチェックを実装してください
	},
}

// Client はWebSocket接続とハブの間の仲介役
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
}

// readPump は WebSocket接続からハブへメッセージを送る
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("エラー: %v", err)
			}
			break
		}

		c.hub.broadcast <- message
	}
}

// writePump はハブからWebSocket接続へメッセージを送る
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 待機中のメッセージを現在のWebSocketメッセージに追加
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs はHTTPサーバーからのWebSocket接続を処理
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// クエリパラメータからユーザー名を取得
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "匿名ユーザー"
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
	}
	client.hub.register <- client

	// goroutineで読み書きを実行
	go client.writePump()
	go client.readPump()
}
