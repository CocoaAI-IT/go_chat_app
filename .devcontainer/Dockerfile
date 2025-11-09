# Go言語の最新版を使用
FROM golang:1.21-bullseye

# 非rootユーザーの作成
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

# 必要なパッケージのインストール
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    git \
    curl \
    wget \
    vim \
    nano \
    build-essential \
    ca-certificates \
    apt-transport-https \
    gnupg \
    lsb-release \
    sudo \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# ユーザーの作成
RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

# Go開発ツールのインストール
RUN go install -v golang.org/x/tools/gopls@latest \
    && go install -v github.com/ramya-rao-a/go-outline@latest \
    && go install -v github.com/go-delve/delve/cmd/dlv@latest \
    && go install -v honnef.co/go/tools/cmd/staticcheck@latest \
    && go install -v github.com/fatih/gomodifytags@latest \
    && go install -v github.com/josharian/impl@latest \
    && go install -v github.com/cweill/gotests/gotests@latest

# Cloudflare Wrangler CLIのインストール（Node.jsは別途Features経由でインストール）
RUN npm install -g wrangler || true

# 作業ディレクトリの設定
WORKDIR /workspaces/go_chat_app

# ユーザーの切り替え
USER $USERNAME

# Go環境変数の設定
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
ENV GO111MODULE=on

# エントリーポイント
CMD ["/bin/bash"]
