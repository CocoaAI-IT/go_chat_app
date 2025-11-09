.PHONY: help run build clean test deps

help: ## ヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## アプリケーションを実行
	go run cmd/server/main.go

build: ## アプリケーションをビルド
	go build -o bin/chat-server cmd/server/main.go

clean: ## ビルド成果物を削除
	rm -rf bin/
	go clean

test: ## テストを実行
	go test -v ./...

deps: ## 依存関係を取得
	go mod download
	go mod tidy

fmt: ## コードをフォーマット
	go fmt ./...

vet: ## コードを静的解析
	go vet ./...

dev: deps ## 開発環境のセットアップ
	@echo "依存関係のインストールが完了しました"
	@echo "make run でアプリケーションを起動できます"
