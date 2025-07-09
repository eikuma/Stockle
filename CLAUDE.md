# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 振る舞い
- **言語**: 日本語で回答してください
- **ファイル編集**: 全てのファイルの末尾に改行を入れてください

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

## チーム開発フロー

### チーム構成
- **PdM**: プロジェクト管理・統合・DevOps（チケット: `00-project-setup/`）
- **Member 1**: フロントエンド開発（チケット: `01-frontend/`）
- **Member 2**: バックエンド基盤開発（チケット: `02-backend-infrastructure/`）
- **Member 3**: バックエンド機能・AI統合開発（チケット: `03-backend-features/`）

### Git Worktree開発フロー

#### ブランチ命名規則
```bash
# フェーズを設定（例: phase1, mvp, v2など）
export PHASE="phase1"

# 各メンバーのブランチ名パターン
feature/${PHASE}-integration        # PdM（統合ブランチ）
feature/${PHASE}-frontend          # Member 1
feature/${PHASE}-backend-infrastructure  # Member 2
feature/${PHASE}-backend-features   # Member 3
```

#### Worktree作成コマンド
```bash
# 最新のmainを取得
git checkout main && git pull origin main

# 各自のworktreeを作成
git worktree add -b feature/${PHASE}-integration worktree-integration          # PdM
git worktree add -b feature/${PHASE}-frontend worktree-frontend               # Member 1
git worktree add -b feature/${PHASE}-backend-infrastructure worktree-backend-infrastructure  # Member 2
git worktree add -b feature/${PHASE}-backend-features worktree-backend-features  # Member 3
```

#### 開発フロー
1. **並列開発**: 各メンバーが自分のworktreeで独立して開発
2. **日次同期**: `git fetch origin && git rebase origin/main`
3. **統合作業**: PdMが各ブランチを統合ブランチにマージ
4. **PR作成**: 統合ブランチからmainへのPR作成

詳細は `TEAM_DEVELOPMENT_FLOW.md` と `docs/git-worktree-guide.md` を参照。

## 📚 プロジェクトドキュメント一覧

### 🔧 開発・運用ガイド
- **`TEAM_DEVELOPMENT_FLOW.md`**: チーム開発フローの詳細（Git Worktree、PR作成手順）
- **`DEVELOPMENT_PROCESS.md`**: 2週間スプリントの詳細プロセス（日次ルーティン、品質保証）
- **`DEVELOPMENT.md`**: 開発環境セットアップとローカル実行手順
- **`docs/git-worktree-guide.md`**: Git Worktreeの使い方とトラブルシューティング
- **`PHASE_NAMING_CONVENTION.md`**: フェーズ命名規則とブランチ管理

### 📋 設計・仕様書（`docs/`フォルダ）
- **`docs/functional-requirements-doc.md`**: 機能要件定義書（Phase 1-4の全機能）
- **`docs/api-design-doc.md`**: REST API設計書（エンドポイント、認証、エラーハンドリング）
- **`docs/database-design-doc.md`**: データベース設計書（ER図、テーブル定義、インデックス）
- **`docs/implementation-guide.md`**: 実装指示書（技術選定、アーキテクチャ、セキュリティ要件）

### 🎫 実装チケット（`tickets/`フォルダ）
- **`tickets/README.md`**: チケット管理システムの概要と並列実行戦略
- **`tickets/IMPLEMENTATION_SUMMARY.md`**: Phase 1実装チケットの完成サマリー
- **`tickets/00-project-setup/`**: PdM担当（プロジェクト基盤、CI/CD、統合管理）
- **`tickets/01-frontend/`**: Member 1担当（Next.js、認証UI、記事管理UI）
- **`tickets/02-backend-infrastructure/`**: Member 2担当（Go基盤、認証システム、記事API）
- **`tickets/03-backend-features/`**: Member 3担当（AI統合、要約生成、非同期処理）

### 🔧 技術ドキュメント（`backend/`フォルダ）
- **`backend/DATABASE_CONNECTION.md`**: MySQL接続設定とトラブルシューティング
- **`backend/JWT_SECRET_GENERATION.md`**: JWT秘密鍵の生成と管理
- **`backend/job_queue_design.md`**: 非同期ジョブキューの設計
- **`backend/http_client_libraries.md`**: HTTP クライアントライブラリの選定
- **`backend/web_scraping_libraries.md`**: Webスクレイピングライブラリの比較

### 🚀 統合・リリース
- **`README.md`**: プロジェクト概要と基本セットアップ
- **`README-INTEGRATION.md`**: チーム統合作業の成果まとめ

### 📝 GitHub テンプレート（`.github/`フォルダ）
- **`.github/ISSUE_TEMPLATE/`**: Issue テンプレート（bug_report, feature_request）
- **`.github/pull_request_template.md`**: Pull Request テンプレート

### 📖 ドキュメント活用方法
1. **開発開始時**: `TEAM_DEVELOPMENT_FLOW.md` でフローを確認
2. **技術仕様確認**: `docs/` フォルダの設計書を参照
3. **実装作業**: `tickets/` フォルダの担当チケットを実行
4. **技術的課題**: `backend/` フォルダの技術ドキュメントを参照
5. **統合・リリース**: `README-INTEGRATION.md` で成果を確認

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
