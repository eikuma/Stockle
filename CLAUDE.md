# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

「Stockle」は記事保存アプリケーションで、ユーザーが後で読みたい記事を保存し、AIによる要約生成機能を持つWebアプリケーションです。現在は設計段階で、`docs/`フォルダに設計書が格納されています。

## 技術スタック

### フロントエンド
- **フレームワーク**: Next.js 14 (App Router)
- **UIライブラリ**: Tailwind CSS + shadcn/ui + Radix UI
- **状態管理**: Zustand
- **データフェッチング**: TanStack Query (React Query)
- **フォーム管理**: React Hook Form + Zod
- **認証**: NextAuth.js

### バックエンド
- **言語**: Go 1.21+
- **フレームワーク**: Gin
- **ORM**: GORM
- **認証**: JWT (golang-jwt)
- **バリデーション**: go-playground/validator
- **開発支援**: Air（ホットリロード）

### データベース・インフラ
- **RDBMS**: MySQL 8.0
- **マイグレーション**: golang-migrate
- **AI/ML**: Groq API（第1選択）、Anthropic Claude API（第2選択）

## プロジェクト構造

```
Stockle/
├── docs/                      # 設計ドキュメント（現在の状態）
│   ├── api-design-doc.md     # API設計書
│   ├── database-design-doc.md # データベース設計書
│   ├── functional-requirements-doc.md # 機能要件定義書
│   └── implementation-guide.md # 実装指示書（重要）
│
# 将来実装予定の構造:
├── frontend/                  # Next.js アプリケーション
│   ├── src/
│   │   ├── app/              # App Router
│   │   ├── components/       # UIコンポーネント
│   │   ├── features/         # 機能別モジュール
│   │   ├── hooks/            # カスタムフック
│   │   ├── lib/              # ユーティリティ
│   │   ├── services/         # API クライアント
│   │   ├── stores/           # Zustand ストア
│   │   └── types/            # TypeScript 型定義
│   └── package.json
│
├── backend/                   # Go API サーバー
│   ├── cmd/api/              # メインエントリーポイント
│   ├── internal/
│   │   ├── config/           # 設定
│   │   ├── controllers/      # HTTPハンドラー
│   │   ├── middleware/       # ミドルウェア
│   │   ├── models/           # データモデル
│   │   ├── repositories/     # データアクセス層
│   │   ├── services/         # ビジネスロジック
│   │   └── validators/       # バリデーション
│   ├── migrations/           # DBマイグレーション
│   ├── pkg/                  # 共有パッケージ
│   └── go.mod
│
└── docker-compose.yml        # ローカル開発環境
```

## 開発コマンド（実装後）

### 初期セットアップ
```bash
# Dockerでの開発環境起動
docker-compose up -d

# フロントエンドのセットアップ
cd frontend
npm install
npm run dev

# バックエンドのセットアップ
cd backend
go mod download
air  # ホットリロード開発サーバー
```

### データベースマイグレーション
```bash
cd backend
migrate -path migrations -database "mysql://user:password@tcp(localhost:3306)/readlater_db" up
```

### ビルド・テストコマンド
```bash
# フロントエンド
cd frontend
npm run build    # ビルド
npm test        # テスト（Vitest）
npm run lint    # ESLint

# バックエンド
cd backend
go build -o api cmd/api/main.go  # ビルド
go test ./...                     # テスト
```

## 実装アーキテクチャ

### 認証フロー
- JWT認証（httpOnly Cookie）
- Google OAuth 2.0統合
- リフレッシュトークンローテーション

### AI統合
- 要約生成: Groq API → Anthropic Claude API（フォールバック）
- 非同期処理とリトライ機構
- 自動カテゴリ分類機能

### セキュリティ対策
- CSRF対策（Double Submit Cookie）
- Rate Limiting（1分あたり100リクエスト）
- SQLインジェクション・XSS対策

## 重要な実装指針

1. **設計書を必読**: `docs/implementation-guide.md`に詳細な実装指示が記載されています
2. **段階的実装**: Phase 1（基盤構築 → 認証 → 基本CRUD → AI統合）の順序で実装
3. **エラーハンドリング**: 統一フォーマットとユーザーフレンドリーなメッセージ
4. **パフォーマンス**: API応答時間 < 200ms、フロントエンドTTI < 3秒
5. **テスト戦略**: フロントエンド（Vitest + Testing Library + Playwright）、バックエンド（Go標準テスト + testcontainers-go）

## 環境変数

### フロントエンド (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-secret-key
GOOGLE_CLIENT_ID=your-google-client-id
```

### バックエンド (.env)
```env
PORT=8080
ENV=development
DB_HOST=localhost
DB_PORT=3306
JWT_SECRET=your-jwt-secret
GROQ_API_KEY=your-groq-key
ANTHROPIC_API_KEY=your-anthropic-key
```