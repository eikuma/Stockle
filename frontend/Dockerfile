# ビルドステージ
FROM node:20-alpine AS builder

WORKDIR /app

# 依存関係をコピー＆インストール
COPY package*.json ./
RUN npm ci

# ソースコードをコピー
COPY . .

# ビルド
RUN npm run build

# 実行ステージ
FROM node:20-alpine

WORKDIR /app

# 必要なファイルのみコピー
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/node_modules ./node_modules

# ポート公開
EXPOSE 3000

# 本番モードで起動
CMD ["npm", "start"]