package models

import (
	"sync"
	"time"
)

// User はアプリケーション内のユーザーを表す
type User struct {
	ID       string `json:"id"`        // WebSocket接続のユニークID
	Nickname string `json:"nickname"`  // ユーザーのニックネーム
	Location string `json:"location"`  // "lobby" or "room"
	RoomID   string `json:"room_id"`   // 現在いる部屋のID（ロビーの場合は空）
}

// Room はチャットルームを表す
type Room struct {
	ID        string            `json:"id"`
	Password  string            `json:"-"` // 合言葉（JSONには含めない）
	Users     map[string]*User  `json:"users"`
	Messages  []*Message        `json:"messages"`
	CreatedAt time.Time         `json:"created_at"`
	mu        sync.RWMutex
}

// Message はチャットメッセージを表す
type Message struct {
	Type      string    `json:"type"`      // メッセージの種類
	Nickname  string    `json:"nickname"`  // 送信者のニックネーム
	Content   string    `json:"content"`   // メッセージ内容
	Timestamp time.Time `json:"timestamp"` // 送信時刻
	RoomID    string    `json:"room_id"`   // 部屋ID
}

// NewRoom は新しいルームを作成
func NewRoom(id, password string) *Room {
	return &Room{
		ID:        id,
		Password:  password,
		Users:     make(map[string]*User),
		Messages:  make([]*Message, 0),
		CreatedAt: time.Now(),
	}
}

// AddUser はルームにユーザーを追加（最大10人まで）
func (r *Room) AddUser(user *User) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.Users) >= 10 {
		return false // 人数制限
	}

	r.Users[user.ID] = user
	return true
}

// RemoveUser はルームからユーザーを削除
func (r *Room) RemoveUser(userID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.Users, userID)
}

// AddMessage はルームにメッセージを追加
func (r *Room) AddMessage(msg *Message) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Messages = append(r.Messages, msg)
}

// IsEmpty はルームが空かどうかを判定
func (r *Room) IsEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.Users) == 0
}

// GetUserCount はルーム内のユーザー数を取得
func (r *Room) GetUserCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.Users)
}
