package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/CocoaAI-IT/go_chat_app/internal/manager"
	"github.com/CocoaAI-IT/go_chat_app/pkg/models"
)

// WSMessage はクライアント間でやり取りされるメッセージの構造
type WSMessage struct {
	Type     string      `json:"type"`
	Data     interface{} `json:"data,omitempty"`
	Error    string      `json:"error,omitempty"`
}

// NicknameRegisterData はニックネーム登録時のデータ
type NicknameRegisterData struct {
	Nickname string `json:"nickname"`
}

// RoomJoinData は部屋参加時のデータ
type RoomJoinData struct {
	Password string `json:"password"`
}

// ChatMessageData はチャットメッセージのデータ
type ChatMessageData struct {
	Content string `json:"content"`
}

// Hub はアクティブなクライアントとメッセージのブロードキャストを管理
type Hub struct {
	manager    *manager.Manager
	clients    map[string]*Client // userID -> Client
	register   chan *Client
	unregister chan *Client
	message    chan *ClientMessage
}

// ClientMessage はクライアントからのメッセージ
type ClientMessage struct {
	client  *Client
	message []byte
}

// NewHub は新しいHubインスタンスを作成
func NewHub() *Hub {
	return &Hub{
		manager:    manager.NewManager(),
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		message:    make(chan *ClientMessage),
	}
}

// Run はハブのメインループを実行
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client
			log.Printf("新しいクライアントが接続しました（合計: %d）", len(h.clients))

		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				// マネージャーからユーザーを削除
				user, err := h.manager.GetUser(client.id)
				if err == nil {
					// ユーザーが部屋にいた場合、退出通知
					if user.RoomID != "" {
						h.notifyRoomUserLeft(user.RoomID, user.Nickname)
					}
				}
				h.manager.UnregisterUser(client.id)

				// クライアントを削除
				delete(h.clients, client.id)
				close(client.send)
				log.Printf("クライアントが切断しました（残り: %d）", len(h.clients))

				// ロビー状態を更新
				h.broadcastLobbyUpdate()
			}

		case clientMsg := <-h.message:
			h.handleMessage(clientMsg.client, clientMsg.message)
		}
	}
}

// handleMessage はクライアントからのメッセージを処理
func (h *Hub) handleMessage(client *Client, message []byte) {
	var wsMsg WSMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		log.Printf("メッセージのパースエラー: %v", err)
		return
	}

	switch wsMsg.Type {
	case "nickname_register":
		h.handleNicknameRegister(client, wsMsg.Data)
	case "room_join":
		h.handleRoomJoin(client, wsMsg.Data)
	case "room_leave":
		h.handleRoomLeave(client)
	case "chat":
		h.handleChat(client, wsMsg.Data)
	default:
		log.Printf("不明なメッセージタイプ: %s", wsMsg.Type)
	}
}

// handleNicknameRegister はニックネーム登録を処理
func (h *Hub) handleNicknameRegister(client *Client, data interface{}) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		h.sendError(client, "無効なデータ形式")
		return
	}

	nickname, ok := dataMap["nickname"].(string)
	if !ok || nickname == "" {
		h.sendError(client, "ニックネームが必要です")
		return
	}

	// 文字数チェック（20文字まで）
	if len([]rune(nickname)) > 20 {
		h.sendError(client, "ニックネームは20文字以内にしてください")
		return
	}

	// ユーザー登録
	err := h.manager.RegisterUser(client.id, nickname)
	if err != nil {
		h.sendError(client, err.Error())
		return
	}

	// 成功レスポンス
	h.sendToClient(client, WSMessage{
		Type: "nickname_register_success",
		Data: map[string]string{"nickname": nickname},
	})

	// ロビー状態を全員に配信
	h.broadcastLobbyUpdate()
}

// handleRoomJoin は部屋参加を処理
func (h *Hub) handleRoomJoin(client *Client, data interface{}) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		h.sendError(client, "無効なデータ形式")
		return
	}

	password, ok := dataMap["password"].(string)
	if !ok || password == "" {
		h.sendError(client, "合言葉が必要です")
		return
	}

	// 部屋に参加
	roomID, err := h.manager.JoinRoom(client.id, password)
	if err != nil {
		h.sendError(client, err.Error())
		return
	}

	user, _ := h.manager.GetUser(client.id)

	// 部屋のメッセージ履歴を取得
	messages, _ := h.manager.GetRoomMessages(roomID)

	// 成功レスポンス（メッセージ履歴付き）
	h.sendToClient(client, WSMessage{
		Type: "room_join_success",
		Data: map[string]interface{}{
			"room_id":  roomID,
			"messages": messages,
		},
	})

	// ロビー状態を更新
	h.broadcastLobbyUpdate()

	// 部屋のメンバーに参加通知
	h.notifyRoomUserJoined(roomID, user.Nickname)

	// 部屋の状態を更新
	h.broadcastRoomUpdate(roomID)
}

// handleRoomLeave は部屋退出を処理
func (h *Hub) handleRoomLeave(client *Client) {
	user, err := h.manager.GetUser(client.id)
	if err != nil {
		h.sendError(client, err.Error())
		return
	}

	roomID := user.RoomID
	if roomID == "" {
		return // 既にロビーにいる
	}

	// 部屋から退出
	err = h.manager.LeaveRoom(client.id)
	if err != nil {
		h.sendError(client, err.Error())
		return
	}

	// 成功レスポンス
	h.sendToClient(client, WSMessage{
		Type: "room_leave_success",
	})

	// 部屋のメンバーに退出通知
	h.notifyRoomUserLeft(roomID, user.Nickname)

	// 部屋の状態を更新
	h.broadcastRoomUpdate(roomID)

	// ロビー状態を更新
	h.broadcastLobbyUpdate()
}

// handleChat はチャットメッセージを処理
func (h *Hub) handleChat(client *Client, data interface{}) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		h.sendError(client, "無効なデータ形式")
		return
	}

	content, ok := dataMap["content"].(string)
	if !ok || content == "" {
		return
	}

	user, err := h.manager.GetUser(client.id)
	if err != nil || user.RoomID == "" {
		return
	}

	// メッセージを作成
	msg := &models.Message{
		Type:      "chat",
		Nickname:  user.Nickname,
		Content:   content,
		Timestamp: time.Now(),
		RoomID:    user.RoomID,
	}

	// 部屋にメッセージを追加
	h.manager.AddMessageToRoom(user.RoomID, msg)

	// 部屋のメンバーにブロードキャスト
	h.broadcastToRoom(user.RoomID, WSMessage{
		Type: "chat",
		Data: msg,
	})
}

// broadcastLobbyUpdate はロビーの状態を全員に配信
func (h *Hub) broadcastLobbyUpdate() {
	lobbyUsers := h.manager.GetLobbyUsers()
	allUsers := h.manager.GetAllUsers()

	msg := WSMessage{
		Type: "lobby_update",
		Data: map[string]interface{}{
			"lobby_users": lobbyUsers,
			"all_users":   allUsers,
		},
	}

	msgBytes, _ := json.Marshal(msg)

	// ロビーにいる全ユーザーに送信
	for _, user := range lobbyUsers {
		if client, ok := h.clients[user.ID]; ok {
			select {
			case client.send <- msgBytes:
			default:
				close(client.send)
				delete(h.clients, user.ID)
			}
		}
	}
}

// broadcastRoomUpdate は部屋の状態を部屋のメンバーに配信
func (h *Hub) broadcastRoomUpdate(roomID string) {
	room, err := h.manager.GetRoom(roomID)
	if err != nil {
		return
	}

	users := make([]*models.User, 0, len(room.Users))
	for _, user := range room.Users {
		users = append(users, user)
	}

	msg := WSMessage{
		Type: "room_update",
		Data: map[string]interface{}{
			"users": users,
			"count": len(users),
		},
	}

	h.broadcastToRoom(roomID, msg)
}

// broadcastToRoom は部屋のメンバーにメッセージを配信
func (h *Hub) broadcastToRoom(roomID string, msg WSMessage) {
	room, err := h.manager.GetRoom(roomID)
	if err != nil {
		return
	}

	msgBytes, _ := json.Marshal(msg)

	for userID := range room.Users {
		if client, ok := h.clients[userID]; ok {
			select {
			case client.send <- msgBytes:
			default:
				close(client.send)
				delete(h.clients, userID)
			}
		}
	}
}

// notifyRoomUserJoined は部屋に参加したことを通知
func (h *Hub) notifyRoomUserJoined(roomID, nickname string) {
	msg := &models.Message{
		Type:      "join",
		Nickname:  nickname,
		Content:   "が入室しました",
		Timestamp: time.Now(),
		RoomID:    roomID,
	}

	h.manager.AddMessageToRoom(roomID, msg)

	h.broadcastToRoom(roomID, WSMessage{
		Type: "chat",
		Data: msg,
	})
}

// notifyRoomUserLeft は部屋から退出したことを通知
func (h *Hub) notifyRoomUserLeft(roomID, nickname string) {
	msg := &models.Message{
		Type:      "leave",
		Nickname:  nickname,
		Content:   "が退室しました",
		Timestamp: time.Now(),
		RoomID:    roomID,
	}

	h.manager.AddMessageToRoom(roomID, msg)

	h.broadcastToRoom(roomID, WSMessage{
		Type: "chat",
		Data: msg,
	})
}

// sendToClient は特定のクライアントにメッセージを送信
func (h *Hub) sendToClient(client *Client, msg WSMessage) {
	msgBytes, _ := json.Marshal(msg)
	select {
	case client.send <- msgBytes:
	default:
		close(client.send)
		delete(h.clients, client.id)
	}
}

// sendError はエラーメッセージを送信
func (h *Hub) sendError(client *Client, errMsg string) {
	h.sendToClient(client, WSMessage{
		Type:  "error",
		Error: errMsg,
	})
}
