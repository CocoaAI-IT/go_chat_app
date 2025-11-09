# Go Chat App

Go言語とWebSocketを使用したリアルタイムチャットアプリケーションです。

## 主な機能

- リアルタイムメッセージング（WebSocket使用）
- ユーザー参加/退出通知
- クリーンでモダンなUI
- DevContainer対応（開発環境の即座セットアップ）
- Cloudflareデプロイ対応

## 技術スタック

- **バックエンド**: Go 1.21+
- **WebSocket**: gorilla/websocket
- **ルーター**: gorilla/mux
- **フロントエンド**: HTML/CSS/JavaScript
- **開発環境**: Docker + VSCode DevContainer

## 環境構築手順（Windows 11 + WSL2）

### 前提条件

1. **WSL2のインストール**
   ```powershell
   wsl --install
   ```

2. **Dockerのインストール**
   - Docker Desktop for Windowsをインストール
   - Settings → General → "Use the WSL2 based engine" を有効化
   - Settings → Resources → WSL Integration で使用するディストリビューションを有効化

3. **VSCodeのインストール**
   - VSCodeをインストール
   - 拡張機能「Remote - Containers」をインストール
   - 拡張機能「WSL」をインストール

### DevContainerでの起動方法

1. **プロジェクトをクローン**
   ```bash
   git clone https://github.com/CocoaAI-IT/go_chat_app.git
   cd go_chat_app
   ```

2. **VSCodeでプロジェクトを開く**
   ```bash
   code .
   ```

3. **DevContainerで開く**
   - VSCodeで `Ctrl+Shift+P` を押す
   - 「Dev Containers: Reopen in Container」を選択
   - コンテナのビルドと起動を待つ（初回は数分かかる場合があります）

4. **依存関係のインストール**
   コンテナ内のターミナルで：
   ```bash
   make deps
   ```

5. **アプリケーションの起動**
   ```bash
   make run
   ```

6. **ブラウザでアクセス**
   - http://localhost:8080 にアクセス

## ローカル環境での起動方法（DevContainer不使用）

### 必要なもの

- Go 1.21以上
- Git

### 手順

1. **リポジトリのクローン**
   ```bash
   git clone https://github.com/CocoaAI-IT/go_chat_app.git
   cd go_chat_app
   ```

2. **依存関係のインストール**
   ```bash
   go mod download
   ```

3. **アプリケーションの起動**
   ```bash
   go run cmd/server/main.go
   ```
   または
   ```bash
   make run
   ```

4. **ブラウザでアクセス**
   - http://localhost:8080

## プロジェクト構造

```
go_chat_app/
├── .devcontainer/          # DevContainer設定
│   ├── devcontainer.json   # VSCode DevContainer設定
│   └── Dockerfile          # 開発環境のDockerfile
├── cmd/
│   └── server/
│       └── main.go         # アプリケーションエントリーポイント
├── internal/
│   ├── handler/            # HTTPハンドラー
│   │   └── home.go
│   └── websocket/          # WebSocket関連
│       ├── client.go       # クライアント管理
│       └── hub.go          # メッセージハブ
├── static/                 # 静的ファイル
│   ├── css/
│   │   └── style.css
│   └── js/
│       └── app.js
├── templates/              # HTMLテンプレート
│   └── index.html
├── go.mod                  # Go modules
├── Makefile               # ビルドタスク
└── README.md
```

## 利用可能なMakeコマンド

```bash
make help       # ヘルプを表示
make run        # アプリケーションを実行
make build      # アプリケーションをビルド
make clean      # ビルド成果物を削除
make test       # テストを実行
make deps       # 依存関係を取得
make fmt        # コードをフォーマット
make vet        # コードを静的解析
make dev        # 開発環境のセットアップ
```

## 使い方

1. ブラウザでアプリケーションにアクセス
2. ユーザー名を入力して「参加する」をクリック
3. メッセージを入力して送信
4. 複数のブラウザタブまたはウィンドウで開くとリアルタイムチャットを体験できます

## Cloudflareへのデプロイ

### 方法1: Cloudflare Pages（静的サイト + Functions）

GoアプリケーションをCloudflare Workersで実行するには、TinyGoを使用してWASMにコンパイルするか、別のホスティングサービスと組み合わせる必要があります。

### 方法2: 外部ホスティング + Cloudflare Tunnel

より実用的な方法として、Goアプリケーションを任意のサーバーでホストし、Cloudflare Tunnelを使用して公開することができます。

```bash
# Cloudflare Tunnelのインストール
cloudflared tunnel login
cloudflared tunnel create go-chat-app
cloudflared tunnel route dns go-chat-app chat.yourdomain.com
cloudflared tunnel run --url http://localhost:8080 go-chat-app
```

## トラブルシューティング

### DevContainerが起動しない

- Docker Desktopが起動していることを確認
- WSL2とDockerの統合が有効になっていることを確認
- VSCodeを再起動してみる

### ポートが既に使用されている

```bash
# ポート8080を使用しているプロセスを確認（WSL内）
lsof -i :8080

# または環境変数でポートを変更
PORT=3000 go run cmd/server/main.go
```

### WebSocket接続エラー

- ブラウザのコンソールでエラーを確認
- ファイアウォールの設定を確認
- HTTPSを使用する場合はWSS接続に変更

## ライセンス

MIT License

## 貢献

プルリクエストを歓迎します！

1. このリポジトリをフォーク
2. 機能ブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## サポート

問題や質問がある場合は、GitHubのIssuesページで報告してください。
