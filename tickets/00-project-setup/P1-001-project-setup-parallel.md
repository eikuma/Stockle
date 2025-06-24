# P1-001: プロジェクト基盤セットアップ（並列実行版）

## 概要
プロジェクト全体の基盤構築と開発環境の統一（4人並列実行）

## 担当者と作業分担
**全チームメンバー（並列実行）**

### 🎯 PdM (Claude Code) 担当作業
**推定時間**: 6時間

#### 1. プロジェクト構造の作成
- [ ] ルートディレクトリにモノレポ構造を作成
- [ ] `frontend/` フォルダを作成
- [ ] `backend/` フォルダを作成
- [ ] `.gitignore` ファイルを作成（Node.js + Go対応）
- [ ] 環境変数テンプレートファイル（`.env.example`）を作成

#### 2. Docker開発環境の構築
- [ ] `docker-compose.yml` を作成
- [ ] MySQL 8.0のDockerコンテナ設定
- [ ] Redis（キャッシュ用）のDockerコンテナ設定
- [ ] 開発用のネットワーク設定
- [ ] ボリュームマウント設定
- [ ] ポート設定（Frontend: 3000, Backend: 8080, MySQL: 3306）

#### 3. GitHub設定
- [ ] ブランチ保護ルールの設定
- [ ] Issue/PR テンプレートの作成
- [ ] ラベルの設定（bug, enhancement, documentation など）
- [ ] マイルストーンの作成（Phase 1, Phase 2, etc）

#### 4. プロジェクト文書の作成
- [ ] 開発者向けREADME.mdの作成
- [ ] Git運用ルールの作成
- [ ] 開発ガイドラインの作成

### 🎯 Member 1 (Frontend) 担当作業
**推定時間**: 4時間

#### 1. フロントエンド環境確認
- [ ] Node.js 18以上がインストールされていることを確認
- [ ] npm/yarn が最新版であることを確認
- [ ] VSCode拡張機能の確認（ES7+ React/Redux/React-Native snippets, Tailwind CSS IntelliSense, Prettier）

#### 2. フロントエンド基本セットアップ準備
- [ ] `frontend/` フォルダ内でNext.js 14の動作確認
- [ ] TypeScript設定の確認
- [ ] ESLint + Prettier設定の確認
- [ ] shadcn/ui のセットアップ準備

#### 3. フロントエンド環境変数の設定
- [ ] `.env.local.example` ファイルの作成
- [ ] NextAuth.js用環境変数の準備
- [ ] API URL設定の準備

#### 4. 依存関係管理の準備
- [ ] `package.json` のテンプレート作成
- [ ] フロントエンド専用の `.gitignore` 追加項目の確認

### 🎯 Member 2 (Backend Infrastructure) 担当作業
**推定時間**: 4時間

#### 1. バックエンド環境確認
- [ ] Go 1.21以上がインストールされていることを確認
- [ ] Docker & Docker Composeがインストールされていることを確認
- [ ] MySQLクライアントツールの確認

#### 2. バックエンド基本セットアップ準備
- [ ] `backend/` フォルダ内でGo moduleの初期化準備
- [ ] 基本的なフォルダ構造の設計確認
- [ ] Makefile のテンプレート作成

#### 3. データベース設定の準備
- [ ] MySQL接続設定の確認
- [ ] データベース名・ユーザー設定の確認
- [ ] マイグレーション用フォルダの準備

#### 4. バックエンド環境変数の設定
- [ ] `.env.example` のバックエンド部分作成
- [ ] JWT秘密鍵の生成方法確認
- [ ] データベース接続文字列の準備

### 🎯 Member 3 (Backend Features) 担当作業
**推定時間**: 3時間

#### 1. AI API接続確認
- [ ] Groq APIキーの取得・動作確認
- [ ] Anthropic Claude APIキーの取得・動作確認
- [ ] API制限・レート制限の確認

#### 2. AI統合環境の準備
- [ ] AI API用の環境変数設定
- [ ] HTTPクライアントの動作確認（curl/Postmanなど）
- [ ] APIレスポンス形式の確認

#### 3. 非同期処理環境の準備
- [ ] Redisの動作確認（Dockerコンテナ経由）
- [ ] ジョブキュー機能の設計確認

#### 4. 外部依存関係の確認
- [ ] Webスクレイピング用ライブラリの確認
- [ ] Go用HTTPクライアントライブラリの確認

## 🔄 並列実行スケジュール

### Day 1 (各自2-3時間)
**午前**: 環境確認・基本セットアップ
- PdM: プロジェクト構造作成 + Docker基本設定
- Member 1: 開発環境確認 + フロントエンド準備
- Member 2: 開発環境確認 + バックエンド準備  
- Member 3: AI API確認 + 外部依存関係確認

**午後**: 統合確認・問題解決
- 全員: 各自の作業結果をプッシュ
- 全員: `docker-compose up` で統合環境の動作確認
- 問題があれば全員でトラブルシューティング

### Day 2 (各自1-2時間)
**午前**: 残作業完了・文書化
- PdM: GitHub設定 + 文書作成完了
- Member 1: フロントエンド設定完了
- Member 2: バックエンド設定完了
- Member 3: AI統合準備完了

**午後**: 最終確認・次フェーズ準備
- 全員: 環境の最終確認
- 各自: 次のチケット（P1-101, P1-201, P1-301）の準備開始

## 🤝 コミュニケーション・同期ポイント

### 必須同期タイミング
1. **Day 1 午前終了時**: 各自の進捗確認（15分）
2. **Day 1 午後開始時**: Docker統合環境の動作確認（30分）
3. **Day 2 午前終了時**: 最終確認・問題がないかチェック（15分）

### 共有が必要な情報
- **PdM → 全員**: Docker設定、ポート番号、環境変数名
- **Member 1 → PdM**: フロントエンド用の環境変数要件
- **Member 2 → PdM**: バックエンド用の環境変数要件
- **Member 3 → 全員**: AI API制限・レート制限情報

## 実装詳細

### PdM作成: docker-compose.yml
```yaml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: readlater_db
      MYSQL_USER: readlater_app
      MYSQL_PASSWORD: secure_password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./backend/migrations:/docker-entrypoint-initdb.d
    command: --default-authentication-plugin=mysql_native_password

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes

volumes:
  mysql_data:
  redis_data:
```

### PdM作成: .env.example
```env
# === Database Configuration ===
DB_HOST=localhost
DB_PORT=3306
DB_USER=readlater_app
DB_PASSWORD=secure_password
DB_NAME=readlater_db

# === Backend Configuration ===
PORT=8080
ENV=development

# === JWT Configuration ===
JWT_SECRET=your-jwt-secret-here-change-in-production
JWT_EXPIRY=7d

# === AI Services Configuration ===
# Groq API (Primary choice)
GROQ_API_KEY=your-groq-api-key-here

# Anthropic Claude API (Fallback)
ANTHROPIC_API_KEY=your-anthropic-api-key-here

# === OAuth Configuration ===
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

# === Frontend Configuration ===
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-nextauth-secret-here

# === Redis Configuration ===
REDIS_URL=redis://localhost:6379
```

### Member 1作成: frontend/package.json (テンプレート)
```json
{
  "name": "stockle-frontend",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "lint:fix": "next lint --fix",
    "type-check": "tsc --noEmit",
    "format": "prettier --write .",
    "format:check": "prettier --check ."
  },
  "dependencies": {
    "next": "^14.0.0",
    "react": "^18.0.0",
    "react-dom": "^18.0.0"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "@types/react": "^18.0.0",
    "@types/react-dom": "^18.0.0",
    "typescript": "^5.0.0",
    "eslint": "^8.0.0",
    "eslint-config-next": "^14.0.0",
    "prettier": "^3.0.0"
  }
}
```

### Member 2作成: backend/Makefile
```makefile
.PHONY: help dev build test clean deps

# Variables
BINARY_NAME=stockle-api
MAIN_PACKAGE=./cmd/api

# Default target
help:
	@echo "Available commands:"
	@echo "  dev           - Run development server with hot reload"
	@echo "  build         - Build the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  deps          - Install dependencies"
	@echo "  clean         - Clean build artifacts"

# Development
dev:
	air

# Build
build:
	go build -o $(BINARY_NAME) $(MAIN_PACKAGE)

# Test
test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Dependencies
deps:
	go mod tidy
	go mod download

# Clean
clean:
	go clean
	rm -f $(BINARY_NAME) coverage.out coverage.html
```

## 受入条件

### 必須条件（全員達成必要）
- [ ] `docker-compose up` でMySQL + Redisが起動する
- [ ] 各メンバーが同じ開発環境を構築できる
- [ ] 環境変数ファイルが正しく設定されている
- [ ] GitHubのブランチ保護が有効になっている
- [ ] 各自の開発環境で基本的な動作確認ができる

### 品質条件
- [ ] README.mdに明確なセットアップ手順が記載されている
- [ ] Docker環境でのヘルスチェックが通る
- [ ] 各サービスのポート衝突がない
- [ ] 環境変数の命名規則が統一されている

## 推定時間
**合計17時間** (並列実行により2日で完了)
- PdM: 6時間
- Member 1: 4時間  
- Member 2: 4時間
- Member 3: 3時間

## 依存関係
- なし（最初に実行する必要があります）

## 完了後の次ステップ
1. **即座に並列開始可能**:
   - Member 1: P1-101 (Frontend基盤セットアップ)
   - Member 2: P1-201 (Backend基盤セットアップ)  
   - Member 3: P1-301 (AI統合基盤構築)
2. **PdM**: P1-002 (CI/CD構築) 開始

## 🚨 注意事項・トラブルシューティング

### よくある問題と解決策
1. **ポート衝突**
   - MySQL(3306), Redis(6379), Backend(8080), Frontend(3000)が他のサービスと衝突する場合
   - 解決策: docker-compose.ymlのポート番号を変更

2. **環境変数の不一致**
   - 各メンバーが異なる環境変数名を使用
   - 解決策: .env.example をベースに統一

3. **Docker権限問題**
   - ボリュームマウント時の権限エラー
   - 解決策: Docker設定確認、必要に応じてsudoを使用

4. **API キー関連**
   - Groq/Claude APIキーの取得に時間がかかる
   - 解決策: Member 3は事前にAPIキーを準備、テスト用のダミーキーで進行

### エスカレーション基準
- **1時間以上**同じ問題で進捗が止まった場合 → チーム全体に相談
- **3時間以上**環境構築が完了しない場合 → スケジュール調整を検討

---

**🎯 この並列実行版により、4人チームが効率的に同時作業を行い、2日でプロジェクト基盤を完成させることができます！**