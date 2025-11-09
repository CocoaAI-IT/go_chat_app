package manager

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/CocoaAI-IT/go_chat_app/pkg/models"
)

var (
	ErrNicknameTaken    = errors.New("ニックネームは既に使用されています")
	ErrRoomFull         = errors.New("部屋が満員です")
	ErrInvalidPassword  = errors.New("合言葉が違います")
	ErrRoomNotFound     = errors.New("部屋が見つかりません")
	ErrUserNotFound     = errors.New("ユーザーが見つかりません")
)

// Manager はアプリケーション全体の状態を管理
type Manager struct {
	users     map[string]*models.User  // userID -> User
	nicknames map[string]string        // nickname -> userID
	rooms     map[string]*models.Room  // roomID -> Room
	lobby     map[string]*models.User  // userID -> User (ロビーにいるユーザー)
	mu        sync.RWMutex
}

// NewManager は新しいManagerを作成
func NewManager() *Manager {
	return &Manager{
		users:     make(map[string]*models.User),
		nicknames: make(map[string]string),
		rooms:     make(map[string]*models.Room),
		lobby:     make(map[string]*models.User),
	}
}

// RegisterUser はユーザーを登録（ニックネーム重複チェック）
func (m *Manager) RegisterUser(userID, nickname string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// ニックネームの重複チェック
	if _, exists := m.nicknames[nickname]; exists {
		return ErrNicknameTaken
	}

	user := &models.User{
		ID:       userID,
		Nickname: nickname,
		Location: "lobby",
	}

	m.users[userID] = user
	m.nicknames[nickname] = userID
	m.lobby[userID] = user

	return nil
}

// UnregisterUser はユーザーを登録解除
func (m *Manager) UnregisterUser(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	user, exists := m.users[userID]
	if !exists {
		return
	}

	// ニックネームマップから削除
	delete(m.nicknames, user.Nickname)

	// ロビーから削除
	delete(m.lobby, userID)

	// 部屋から削除
	if user.RoomID != "" {
		if room, exists := m.rooms[user.RoomID]; exists {
			room.RemoveUser(userID)
			// 部屋が空になったら削除
			if room.IsEmpty() {
				delete(m.rooms, user.RoomID)
			}
		}
	}

	// ユーザー削除
	delete(m.users, userID)
}

// JoinRoom はユーザーを部屋に参加させる
func (m *Manager) JoinRoom(userID, password string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	user, exists := m.users[userID]
	if !exists {
		return "", ErrUserNotFound
	}

	// 合言葉で部屋を検索、なければ作成
	var targetRoom *models.Room
	var roomID string

	for id, room := range m.rooms {
		if room.Password == password {
			targetRoom = room
			roomID = id
			break
		}
	}

	// 部屋が見つからない場合は新規作成
	if targetRoom == nil {
		roomID = generateID()
		targetRoom = models.NewRoom(roomID, password)
		m.rooms[roomID] = targetRoom
	}

	// 人数チェック
	if !targetRoom.AddUser(user) {
		return "", ErrRoomFull
	}

	// ユーザーの状態を更新
	user.Location = "room"
	user.RoomID = roomID

	// ロビーから削除
	delete(m.lobby, userID)

	return roomID, nil
}

// LeaveRoom はユーザーを部屋から退出させる
func (m *Manager) LeaveRoom(userID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	user, exists := m.users[userID]
	if !exists {
		return ErrUserNotFound
	}

	if user.RoomID == "" {
		return nil // 既にロビーにいる
	}

	// 部屋から削除
	if room, exists := m.rooms[user.RoomID]; exists {
		room.RemoveUser(userID)
		// 部屋が空になったら削除
		if room.IsEmpty() {
			delete(m.rooms, user.RoomID)
		}
	}

	// ユーザーの状態を更新
	user.Location = "lobby"
	user.RoomID = ""

	// ロビーに追加
	m.lobby[userID] = user

	return nil
}

// GetLobbyUsers はロビーにいるユーザー一覧を取得
func (m *Manager) GetLobbyUsers() []*models.User {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]*models.User, 0, len(m.lobby))
	for _, user := range m.lobby {
		users = append(users, user)
	}
	return users
}

// GetAllUsers は全オンラインユーザーを取得
func (m *Manager) GetAllUsers() []*models.User {
	m.mu.RLock()
	defer m.mu.RUnlock()

	users := make([]*models.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users
}

// GetRoom は部屋を取得
func (m *Manager) GetRoom(roomID string) (*models.Room, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	room, exists := m.rooms[roomID]
	if !exists {
		return nil, ErrRoomNotFound
	}
	return room, nil
}

// GetUser はユーザーを取得
func (m *Manager) GetUser(userID string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[userID]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// AddMessageToRoom は部屋にメッセージを追加
func (m *Manager) AddMessageToRoom(roomID string, msg *models.Message) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	room, exists := m.rooms[roomID]
	if !exists {
		return ErrRoomNotFound
	}

	room.AddMessage(msg)
	return nil
}

// GetRoomMessages は部屋のメッセージ履歴を取得
func (m *Manager) GetRoomMessages(roomID string) ([]*models.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	room, exists := m.rooms[roomID]
	if !exists {
		return nil, ErrRoomNotFound
	}

	return room.Messages, nil
}

// generateID はランダムなIDを生成
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
