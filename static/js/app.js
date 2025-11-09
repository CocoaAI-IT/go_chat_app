// ====================================
// グローバル変数
// ====================================
let ws = null;
let currentNickname = '';
let currentScreen = 'nickname'; // 'nickname', 'lobby', 'room'
let laneIndex = 0; // 次に使用するレーン

// ====================================
// DOM要素の取得
// ====================================

// 画面要素
const nicknameScreen = document.getElementById('nickname-screen');
const lobbyScreen = document.getElementById('lobby-screen');
const roomScreen = document.getElementById('room-screen');

// ニックネーム登録
const nicknameInput = document.getElementById('nickname-input');
const nicknameSubmitBtn = document.getElementById('nickname-submit-btn');
const nicknameError = document.getElementById('nickname-error');

// ロビー
const currentNicknameSpan = document.getElementById('current-nickname');
const lobbyUsersDiv = document.getElementById('lobby-users');
const allUsersDiv = document.getElementById('all-users');
const lobbyCountSpan = document.getElementById('lobby-count');
const totalCountSpan = document.getElementById('total-count');
const passwordInput = document.getElementById('password-input');
const roomJoinBtn = document.getElementById('room-join-btn');
const roomJoinError = document.getElementById('room-join-error');
const lobbyLeaveBtn = document.getElementById('lobby-leave-btn');

// 部屋
const roomUserCountSpan = document.getElementById('room-user-count');
const roomUsersDiv = document.getElementById('room-users');
const commentDisplayDiv = document.getElementById('comment-display');
const logMessagesDiv = document.getElementById('log-messages');
const commentInput = document.getElementById('comment-input');
const commentSendBtn = document.getElementById('comment-send-btn');
const roomLeaveBtn = document.getElementById('room-leave-btn');

// ====================================
// 初期化
// ====================================
document.addEventListener('DOMContentLoaded', () => {
    setupEventListeners();
    connectWebSocket();
});

// ====================================
// イベントリスナーの設定
// ====================================
function setupEventListeners() {
    // ニックネーム登録
    nicknameSubmitBtn.addEventListener('click', registerNickname);
    nicknameInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') registerNickname();
    });

    // ロビー
    roomJoinBtn.addEventListener('click', joinRoom);
    passwordInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') joinRoom();
    });
    lobbyLeaveBtn.addEventListener('click', leaveToNicknameScreen);

    // 部屋
    commentSendBtn.addEventListener('click', sendComment);
    commentInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') sendComment();
    });
    roomLeaveBtn.addEventListener('click', leaveRoom);
}

// ====================================
// WebSocket接続
// ====================================
function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        console.log('WebSocket接続が確立されました');
    };

    ws.onmessage = (event) => {
        try {
            const message = JSON.parse(event.data);
            handleWebSocketMessage(message);
        } catch (e) {
            console.error('メッセージのパースエラー:', e);
        }
    };

    ws.onerror = (error) => {
        console.error('WebSocketエラー:', error);
    };

    ws.onclose = () => {
        console.log('WebSocket接続が切断されました');
        // 再接続を試みる
        setTimeout(() => {
            console.log('再接続を試みています...');
            connectWebSocket();
        }, 3000);
    };
}

// WebSocketメッセージハンドラー
function handleWebSocketMessage(message) {
    console.log('受信:', message);

    switch (message.type) {
        case 'nickname_register_success':
            onNicknameRegisterSuccess(message.data);
            break;
        case 'lobby_update':
            onLobbyUpdate(message.data);
            break;
        case 'room_join_success':
            onRoomJoinSuccess(message.data);
            break;
        case 'room_leave_success':
            onRoomLeaveSuccess();
            break;
        case 'room_update':
            onRoomUpdate(message.data);
            break;
        case 'chat':
            onChatMessage(message.data);
            break;
        case 'error':
            onError(message.error);
            break;
        default:
            console.log('不明なメッセージタイプ:', message.type);
    }
}

// WebSocketメッセージ送信
function sendWebSocketMessage(type, data = {}) {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ type, data }));
    } else {
        console.error('WebSocket未接続');
    }
}

// ====================================
// ニックネーム登録
// ====================================
function registerNickname() {
    const nickname = nicknameInput.value.trim();

    // バリデーション
    if (!nickname) {
        showError(nicknameError, 'ニックネームを入力してください');
        return;
    }

    if (nickname.length > 20) {
        showError(nicknameError, 'ニックネームは20文字以内にしてください');
        return;
    }

    // サーバーに送信
    sendWebSocketMessage('nickname_register', { nickname });
    nicknameSubmitBtn.disabled = true;
}

function onNicknameRegisterSuccess(data) {
    currentNickname = data.nickname;
    currentNicknameSpan.textContent = currentNickname;
    nicknameError.textContent = '';
    nicknameSubmitBtn.disabled = false;

    // ロビー画面へ
    switchScreen('lobby');
}

// ====================================
// ロビー
// ====================================
function onLobbyUpdate(data) {
    const lobbyUsers = data.lobby_users || [];
    const allUsers = data.all_users || [];

    // ロビーユーザー表示
    lobbyCountSpan.textContent = lobbyUsers.length;
    lobbyUsersDiv.innerHTML = '';
    lobbyUsers.forEach(user => {
        const userEl = createUserElement(user);
        lobbyUsersDiv.appendChild(userEl);
    });

    // 全ユーザー表示
    totalCountSpan.textContent = allUsers.length;
    allUsersDiv.innerHTML = '';
    allUsers.forEach(user => {
        const userEl = createUserElement(user);
        if (user.location === 'room') {
            userEl.classList.add('in-room');
        }
        allUsersDiv.appendChild(userEl);
    });
}

function createUserElement(user) {
    const div = document.createElement('div');
    div.className = 'user-item';
    div.textContent = user.nickname;
    return div;
}

// ====================================
// 部屋参加
// ====================================
function joinRoom() {
    const password = passwordInput.value.trim();

    if (!password) {
        showError(roomJoinError, '合言葉を入力してください');
        return;
    }

    sendWebSocketMessage('room_join', { password });
    roomJoinBtn.disabled = true;
}

function onRoomJoinSuccess(data) {
    roomJoinBtn.disabled = false;
    roomJoinError.textContent = '';
    passwordInput.value = '';

    // メッセージ履歴を表示
    const messages = data.messages || [];
    logMessagesDiv.innerHTML = '';
    messages.forEach(msg => {
        addLogMessage(msg);
    });

    // 部屋画面へ
    switchScreen('room');
}

// ====================================
// 部屋退出
// ====================================
function leaveRoom() {
    sendWebSocketMessage('room_leave');
}

function onRoomLeaveSuccess() {
    // ログをクリア
    logMessagesDiv.innerHTML = '';

    // コメント表示エリアをクリア
    const lanes = commentDisplayDiv.querySelectorAll('.comment-lane');
    lanes.forEach(lane => {
        lane.innerHTML = '';
    });

    commentInput.value = '';

    // ロビー画面へ
    switchScreen('lobby');
}

// ====================================
// 部屋の状態更新
// ====================================
function onRoomUpdate(data) {
    const users = data.users || [];
    const count = data.count || 0;

    roomUserCountSpan.textContent = count;
    roomUsersDiv.innerHTML = '';
    users.forEach(user => {
        const userEl = createUserElement(user);
        roomUsersDiv.appendChild(userEl);
    });
}

// ====================================
// チャットメッセージ
// ====================================
function sendComment() {
    const content = commentInput.value.trim();

    if (!content) return;

    sendWebSocketMessage('chat', { content });
    commentInput.value = '';
}

function onChatMessage(data) {
    // ログに追加
    addLogMessage(data);

    // ニコニコ風に流す
    displayFlowingComment(data);
}

// ====================================
// ニコニコ風コメント表示
// ====================================
function displayFlowingComment(message) {
    // 5つのレーンを順番に使用
    const lane = laneIndex;
    laneIndex = (laneIndex + 1) % 5;

    const laneElement = commentDisplayDiv.querySelector(`[data-lane="${lane}"]`);
    if (!laneElement) return;

    // コメント要素を作成
    const commentEl = document.createElement('div');
    commentEl.className = 'flowing-comment';
    commentEl.textContent = `${message.nickname}: ${message.content}`;

    // ランダムな縦位置（レーン内での微調整）
    const randomTop = Math.random() * 60 + 20; // 20-80%
    commentEl.style.top = `${randomTop}%`;

    // レーンに追加
    laneElement.appendChild(commentEl);

    // 5秒後に削除
    setTimeout(() => {
        if (commentEl.parentNode) {
            commentEl.parentNode.removeChild(commentEl);
        }
    }, 5000);
}

// ====================================
// ログ表示
// ====================================
function addLogMessage(message) {
    const logItem = document.createElement('div');
    logItem.className = 'log-item';

    if (message.type === 'join') {
        logItem.classList.add('join');
    } else if (message.type === 'leave') {
        logItem.classList.add('leave');
    }

    const header = document.createElement('div');
    header.className = 'log-item-header';

    const nickname = document.createElement('span');
    nickname.textContent = message.nickname;

    const time = document.createElement('span');
    time.className = 'log-item-time';
    time.textContent = formatTime(message.timestamp);

    header.appendChild(nickname);
    header.appendChild(time);

    const content = document.createElement('div');
    content.className = 'log-item-content';
    content.textContent = message.content;

    logItem.appendChild(header);
    logItem.appendChild(content);

    logMessagesDiv.appendChild(logItem);

    // 自動スクロール
    logMessagesDiv.scrollTop = logMessagesDiv.scrollHeight;

    // 最大100件まで保持
    while (logMessagesDiv.children.length > 100) {
        logMessagesDiv.removeChild(logMessagesDiv.firstChild);
    }
}

// ====================================
// エラー処理
// ====================================
function onError(errorMessage) {
    console.error('エラー:', errorMessage);

    if (currentScreen === 'nickname') {
        showError(nicknameError, errorMessage);
        nicknameSubmitBtn.disabled = false;
    } else if (currentScreen === 'lobby') {
        showError(roomJoinError, errorMessage);
        roomJoinBtn.disabled = false;
    }
}

function showError(element, message) {
    element.textContent = message;
    element.style.display = 'block';

    // 5秒後に消す
    setTimeout(() => {
        element.textContent = '';
    }, 5000);
}

// ====================================
// 画面切り替え
// ====================================
function switchScreen(screenName) {
    // 全画面を非表示
    nicknameScreen.classList.remove('active');
    lobbyScreen.classList.remove('active');
    roomScreen.classList.remove('active');

    // 指定された画面を表示
    switch (screenName) {
        case 'nickname':
            nicknameScreen.classList.add('active');
            nicknameInput.focus();
            break;
        case 'lobby':
            lobbyScreen.classList.add('active');
            passwordInput.focus();
            break;
        case 'room':
            roomScreen.classList.add('active');
            commentInput.focus();
            break;
    }

    currentScreen = screenName;
}

// ====================================
// ロビーへ戻る（ニックネーム登録画面から）
// ====================================
function leaveToNicknameScreen() {
    // WebSocket切断して再接続
    if (ws) {
        ws.close();
    }

    currentNickname = '';
    nicknameInput.value = '';
    passwordInput.value = '';
    nicknameError.textContent = '';
    roomJoinError.textContent = '';

    // 再接続
    setTimeout(() => {
        connectWebSocket();
    }, 100);

    switchScreen('nickname');
}

// ====================================
// ユーティリティ関数
// ====================================

// 時刻フォーマット
function formatTime(timestamp) {
    const date = new Date(timestamp);
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    const seconds = date.getSeconds().toString().padStart(2, '0');
    return `${hours}:${minutes}:${seconds}`;
}

// HTMLエスケープ（XSS対策は必要に応じて）
function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return text.replace(/[&<>"']/g, (m) => map[m]);
}

// ページを離れる時の処理
window.addEventListener('beforeunload', () => {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.close();
    }
});
