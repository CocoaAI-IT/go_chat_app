// グローバル変数
let ws = null;
let username = '';

// DOM要素の取得
const loginArea = document.getElementById('login-area');
const chatArea = document.getElementById('chat-area');
const usernameInput = document.getElementById('username-input');
const joinBtn = document.getElementById('join-btn');
const currentUsername = document.getElementById('current-username');
const messagesDiv = document.getElementById('messages');
const messageInput = document.getElementById('message-input');
const sendBtn = document.getElementById('send-btn');

// ログイン処理
joinBtn.addEventListener('click', () => {
    const inputUsername = usernameInput.value.trim();

    if (inputUsername === '') {
        alert('ユーザー名を入力してください');
        return;
    }

    username = inputUsername;
    currentUsername.textContent = username;

    // WebSocket接続の確立
    connectWebSocket();

    // UIの切り替え
    loginArea.style.display = 'none';
    chatArea.style.display = 'flex';
    messageInput.focus();
});

// Enterキーでログイン
usernameInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        joinBtn.click();
    }
});

// WebSocket接続の確立
function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws?username=${encodeURIComponent(username)}`;

    ws = new WebSocket(wsUrl);

    ws.onopen = () => {
        console.log('WebSocket接続が確立されました');
        sendBtn.disabled = false;
        addSystemMessage('チャットに接続しました');
    };

    ws.onmessage = (event) => {
        try {
            const message = JSON.parse(event.data);
            displayMessage(message);
        } catch (e) {
            console.error('メッセージのパースエラー:', e);
        }
    };

    ws.onerror = (error) => {
        console.error('WebSocketエラー:', error);
        addSystemMessage('接続エラーが発生しました', 'error');
    };

    ws.onclose = () => {
        console.log('WebSocket接続が切断されました');
        sendBtn.disabled = true;
        addSystemMessage('接続が切断されました', 'error');

        // 再接続を試みる
        setTimeout(() => {
            if (username) {
                addSystemMessage('再接続を試みています...', 'info');
                connectWebSocket();
            }
        }, 3000);
    };
}

// メッセージ送信
function sendMessage() {
    const content = messageInput.value.trim();

    if (content === '' || !ws || ws.readyState !== WebSocket.OPEN) {
        return;
    }

    const message = {
        username: username,
        content: content,
        type: 'message'
    };

    ws.send(JSON.stringify(message));
    messageInput.value = '';
    messageInput.focus();
}

// 送信ボタンのイベントリスナー
sendBtn.addEventListener('click', sendMessage);

// Enterキーでメッセージ送信
messageInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        sendMessage();
    }
});

// メッセージの表示
function displayMessage(message) {
    const messageDiv = document.createElement('div');
    messageDiv.className = 'message';

    // メッセージタイプに応じてクラスを追加
    if (message.type === 'join') {
        messageDiv.classList.add('join-message');
    } else if (message.type === 'leave') {
        messageDiv.classList.add('leave-message');
    } else {
        messageDiv.classList.add('user-message');
    }

    // タイムスタンプの生成
    const now = new Date();
    const timeStr = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`;

    // メッセージHTML
    messageDiv.innerHTML = `
        <span class="username">${escapeHtml(message.username)}</span>
        <span class="content">${escapeHtml(message.content)}</span>
        <span class="timestamp">${timeStr}</span>
    `;

    messagesDiv.appendChild(messageDiv);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

// システムメッセージの追加
function addSystemMessage(content, type = 'info') {
    const messageDiv = document.createElement('div');
    messageDiv.className = 'message';

    if (type === 'error') {
        messageDiv.classList.add('leave-message');
    } else {
        messageDiv.classList.add('join-message');
    }

    const now = new Date();
    const timeStr = `${now.getHours().toString().padStart(2, '0')}:${now.getMinutes().toString().padStart(2, '0')}`;

    messageDiv.innerHTML = `
        <span class="username">システム</span>
        <span class="content">${escapeHtml(content)}</span>
        <span class="timestamp">${timeStr}</span>
    `;

    messagesDiv.appendChild(messageDiv);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

// HTMLエスケープ関数（XSS対策）
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

// ページを離れる時の警告
window.addEventListener('beforeunload', (e) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
        ws.close();
    }
});
