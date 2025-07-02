# P1-001: プロジェクト基盤セットアップ

## 概要
プロジェクト全体の基盤構築と開発環境の統一

## 担当者
**PdM (Claude Code)**

## 優先度
**最高** - 全チームメンバーの作業開始に必要

## 前提条件
- GitHubリポジトリが作成済み
- 各メンバーのGitHubアクセス権限が設定済み

## 作業内容

### 1. プロジェクト構造の作成
- [x] ルートディレクトリにモノレポ構造を作成
- [x] `frontend/` フォルダを作成
- [x] `backend/` フォルダを作成
- [x] `docs/` フォルダを整理
- [x] `.gitignore` ファイルを作成（Node.js + Go対応）
- [x] `docker-compose.yml` を作成
- [x] 環境変数テンプレートファイルを作成

### 2. Docker開発環境の構築
- [x] MySQL 8.0のDockerコンテナ設定
- [x] Redis（キャッシュ用）のDockerコンテナ設定
- [x] 開発用のネットワーク設定
- [x] ボリュームマウント設定
- [x] ポート設定（Frontend: 3000, Backend: 8080, MySQL: 3306）

### 3. 環境変数管理
- [x] `.env.example` ファイルを作成
- [x] 開発環境用の環境変数を設定
- [x] 各メンバーに環境変数ファイルを配布
- [x] シークレット管理方針の策定

### 4. GitHub設定
- [x] ブランチ保護ルールの設定
- [x] Issue/PR テンプレートの作成
- [x] ラベルの設定
- [x] マイルストーンの作成

### 5. プロジェクト文書の作成
- [x] 開発者向けREADME.mdの作成
- [x] API仕様書のテンプレート作成
- [x] 開発ガイドラインの作成
- [x] Git運用ルールの作成

## 実装詳細

### docker-compose.yml
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

volumes:
  mysql_data:
  redis_data:
```

### .env.example
```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=readlater_app
DB_PASSWORD=secure_password
DB_NAME=readlater_db

# JWT
JWT_SECRET=your-jwt-secret-here
JWT_EXPIRY=7d

# AI Services
GROQ_API_KEY=your-groq-api-key
ANTHROPIC_API_KEY=your-anthropic-api-key

# OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

# Frontend
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-nextauth-secret
```

## 受入条件
### 必須条件
- [x] `docker-compose up` でMySQL + Redisが起動する
- [x] 各メンバーが同じ開発環境を構築できる
- [x] 環境変数ファイルが正しく配布されている
- [x] GitHubのブランチ保護が有効になっている

### 品質条件
- [x] README.mdに明確なセットアップ手順が記載されている
- [x] Docker環境でのホットリロードが動作する
- [x] 各サービスのヘルスチェックが通る

## 推定時間
**8時間** (1-2日)

## 依存関係
- なし（最初に実行する必要があります）

## 完了後の次ステップ
1. Member 1: Frontend基盤セットアップ開始
2. Member 2: Backend基盤セットアップ開始
3. Member 3: AI統合準備開始

## 備考
- 各メンバーの作業開始前に必ず完了させること
- 問題があれば即座にチーム全体に共有すること
- 環境変数は絶対にコミットしないこと