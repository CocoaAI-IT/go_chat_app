# Go Chat App

Go言語とWebSocketを使用したニコニコ動画風リアルタイムチャットアプリケーションです。

## 主な機能

### コア機能
- **ニックネーム登録**: 重複チェック付き（最大20文字、絵文字・特殊文字対応）
- **ロビーシステム**:
  - 現在ロビーにいるユーザー一覧表示
  - 全オンラインユーザー一覧表示（ロビー/部屋の区別）
- **部屋システム**:
  - 合言葉による完全非公開ルーム
  - 最大10人までの人数制限
  - 誰もいなくなると自動削除
- **ニコニコ動画風コメント表示**:
  - 右から左に流れるコメント（5秒で画面横断）
  - 5つのレーンに自動振り分け（コメント重複防止）
- **ログパネル**: 入室時からのチャット履歴を常時表示
- **レスポンシブデザイン**: PC・タブレット・スマートフォン対応

### 技術的特徴
- リアルタイム通信（WebSocket）
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
│   ├── manager/            # 状態管理
│   │   └── manager.go      # ユーザー・ロビー・部屋管理
│   └── websocket/          # WebSocket関連
│       ├── client.go       # クライアント管理
│       └── hub.go          # メッセージハブ
├── pkg/
│   └── models/             # データモデル
│       └── models.go       # User, Room, Message
├── static/                 # 静的ファイル
│   ├── css/
│   │   └── style.css       # レスポンシブデザイン
│   └── js/
│       └── app.js          # フロントエンドロジック
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

### 基本的な流れ

1. **ニックネーム登録**
   - アプリケーションにアクセス（http://localhost:8080）
   - ニックネームを入力（最大20文字、絵文字も使えます）
   - 「登録してロビーへ」をクリック
   - ※ 既に使われているニックネームは使用できません

2. **ロビー**
   - 現在ロビーにいるユーザーと、全オンラインユーザーが表示されます
   - 部屋に入るには「合言葉」を入力して「部屋に入る」をクリック
   - 同じ合言葉を入力した人が同じ部屋に集まります
   - 「ニックネーム登録画面へ戻る」で最初に戻れます

3. **部屋でチャット**
   - **ニコニコ風コメント**: 入力したコメントが画面上を右から左へ流れます
   - **ログパネル**: 右側（スマホは下側）にログが表示され、過去のメッセージを確認できます
   - **参加者リスト**: 部屋にいる人が表示されます（最大10人）
   - 「ロビーに戻る」で退出できます

### 複数人でテストする方法

- 同じPCで複数のブラウザタブを開く
- 異なるニックネームで登録
- 同じ合言葉で部屋に入る
- リアルタイムでコメントが流れることを確認

### モバイルでの利用

- スマートフォンのブラウザでもアクセス可能
- レスポンシブデザインで最適化されています

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
