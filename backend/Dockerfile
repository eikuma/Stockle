# ビルドステージ
FROM golang:1.23-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git make

WORKDIR /app

# 依存関係をコピー＆インストール
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# ビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/api/main.go

# 実行ステージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# ビルド済みバイナリをコピー
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

# ポート公開
EXPOSE 8080

# 実行
CMD ["./main"]