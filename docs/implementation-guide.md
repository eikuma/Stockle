# 記事保存アプリ - 実装指示書

## 1. 概要

本ドキュメントは、Claude Codeが記事保存アプリを実装する際の指示書です。以下の設計書に基づいて実装を行ってください：

- 機能要件定義書
- API設計書
- データベース設計書

## 2. 技術スタック

### 2.1 フロントエンド
- **フレームワーク**: Next.js 14 (App Router)
- **UI ライブラリ**: 
  - Tailwind CSS
  - shadcn/ui
  - Radix UI
- **状態管理**: Zustand
- **データフェッチング**: TanStack Query (React Query)
- **フォーム管理**: React Hook Form + Zod
- **認証**: NextAuth.js
- **その他**:
  - TypeScript
  - ESLint + Prettier
  - Framer Motion（アニメーション）

### 2.2 バックエンド
- **言語**: Go 1.21+
- **フレームワーク**: Gin
- **ORM**: GORM
- **認証**: JWT (golang-jwt)
- **バリデーション**: go-playground/validator
- **その他**:
  - Air（ホットリロード）
  - Swagger（API ドキュメント）
  - Zap（ロギング）

### 2.3 データベース
- **RDBMS**: MySQL 8.0
- **マイグレーション**: golang-migrate
- **接続プール**: database/sql

### 2.4 AI/ML（コスト最適化）
- **要約生成**: 
  - 第1選択: Groq API（無料枠）
  - 第2選択: Anthropic Claude API（従量課金）
- **音声合成**: Edge-TTS（無料）

### 2.5 インフラ
- **フロントエンド**: Vercel（無料枠）
- **バックエンド**: Railway または Fly.io
- **データベース**: PlanetScale（無料枠）または Neon
- **ファイルストレージ**: Cloudflare R2（無料枠）

## 3. プロジェクト構造

### 3.1 モノレポ構成

```
readlater-app/
├── frontend/               # Next.js アプリケーション
│   ├── src/
│   │   ├── app/           # App Router
│   │   ├── components/    # UIコンポーネント
│   │   ├── features/      # 機能別モジュール
│   │   ├── hooks/         # カスタムフック
│   │   ├── lib/          # ユーティリティ
│   │   ├── services/      # API クライアント
│   │   ├── stores/        # Zustand ストア
│   │   └── types/         # TypeScript 型定義
│   ├── public/
│   └── package.json
│
├── backend/                # Go API サーバー
│   ├── cmd/
│   │   └── api/          # メインエントリーポイント
│   ├── internal/
│   │   ├── config/       # 設定
│   │   ├── controllers/  # HTTPハンドラー
│   │   ├── middleware/   # ミドルウェア
│   │   ├── models/       # データモデル
│   │   ├── repositories/ # データアクセス層
│   │   ├── services/     # ビジネスロジック
│   │   └── validators/   # バリデーション
│   ├── migrations/        # DBマイグレーション
│   ├── pkg/              # 共有パッケージ
│   └── go.mod
│
├── docker-compose.yml     # ローカル開発環境
└── README.md
```

## 4. 実装順序（Phase 1）

### Step 1: 基盤構築
1. **プロジェクトセットアップ**
   - モノレポの初期化
   - 開発環境の構築（Docker Compose）
   - 環境変数の設定

2. **データベース構築**
   - MySQLのセットアップ
   - マイグレーションファイルの作成
   - 初期データの投入

### Step 2: バックエンド基本実装
1. **認証システム**
   - JWT認証の実装
   - メール/パスワード認証
   - Google OAuth 2.0統合
   - セッション管理

2. **基本CRUD API**
   - ユーザー管理
   - カテゴリ管理
   - 記事保存・取得
   - タグ管理

### Step 3: フロントエンド基本実装
1. **認証画面**
   - ログイン/サインアップ
   - Google認証ボタン
   - パスワードリセット

2. **メイン画面**
   - 記事一覧表示
   - 記事保存フォーム
   - カテゴリフィルター
   - 検索機能

### Step 4: AI機能統合
1. **要約生成**
   - Groq API統合
   - 非同期処理の実装
   - リトライ機構

2. **自動カテゴリ分類**
   - 分類ロジックの実装
   - 信頼度スコアの計算

## 5. UI/UXデザイン指針

### 5.1 デザインシステム

**カラーパレット**
```css
:root {
  /* Primary */
  --primary-50: #eff6ff;
  --primary-500: #3b82f6;
  --primary-600: #2563eb;
  
  /* Neutral */
  --gray-50: #f9fafb;
  --gray-100: #f3f4f6;
  --gray-900: #111827;
  
  /* Success/Error */
  --success: #10b981;
  --error: #ef4444;
}
```

**タイポグラフィ**
- フォント: Inter（英数字）、Noto Sans JP（日本語）
- 見出し: font-bold
- 本文: font-normal

### 5.2 主要画面のレイアウト

#### ダッシュボード（記事一覧）
```
┌─────────────────────────────────────────────────┐
│ [Logo] ReadLater    [Search...]    [+] [User]   │ ← ヘッダー
├─────────────────────────────────────────────────┤
│ ┌───────────┐ ┌─────────────────────────────┐  │
│ │           │ │ □ 未読のみ  ⚡ 並び順 ▼     │  │
│ │ すべて    │ ├─────────────────────────────┤  │
│ │ テクノロジー │ │ ┌─────────────────────┐   │  │
│ │ ビジネス   │ │ │ [サムネイル]         │   │  │
│ │ ライフ     │ │ │ タイトル             │   │  │
│ │           │ │ │ 要約テキスト...      │   │  │
│ │ + カテゴリ │ │ │ 🏷️ AI 📅 2時間前    │   │  │
│ │           │ │ └─────────────────────┘   │  │
│ └───────────┘ │ （記事カードの繰り返し）    │  │
│  サイドバー    │                            │  │
│               └─────────────────────────────┘  │
└─────────────────────────────────────────────────┘
```

#### 記事保存モーダル
```
┌─────────────────────────────────────┐
│ 記事を保存              [×]         │
├─────────────────────────────────────┤
│ URL                                 │
│ ┌─────────────────────────────────┐ │
│ │ https://...                     │ │
│ └─────────────────────────────────┘ │
│                                     │
│ カテゴリ                            │
│ [テクノロジー ▼]                    │
│                                     │
│ タグ（任意）                        │
│ ┌─────────────────────────────────┐ │
│ │ AI, 機械学習                    │ │
│ └─────────────────────────────────┘ │
│                                     │
│ [キャンセル]         [保存]         │
└─────────────────────────────────────┘
```

#### 記事詳細画面
```
┌─────────────────────────────────────────────────┐
│ [← 戻る]                          [編集] [削除] │
├─────────────────────────────────────────────────┤
│ # 記事タイトル                                  │
│                                                 │
│ 📅 2024年1月1日 | 👤 著者名 | 🏷️ テクノロジー    │
│                                                 │
│ ## 要約                                         │
│ ┌─────────────────────────────────────────────┐ │
│ │ AIが生成した要約テキスト...                  │ │
│ └─────────────────────────────────────────────┘ │
│                                                 │
│ [元記事を開く] [既読にする]                     │
└─────────────────────────────────────────────────┘
```

### 5.3 レスポンシブデザイン

**ブレークポイント**
- Mobile: < 640px
- Tablet: 640px - 1024px
- Desktop: > 1024px

**モバイル対応**
- サイドバーはハンバーガーメニューに
- カードレイアウトは1カラムに
- タッチフレンドリーなUI（最小44pxのタップエリア）

## 6. セキュリティ実装要件

### 6.1 認証・認可
- JWTトークンの安全な保存（httpOnly Cookie）
- リフレッシュトークンのローテーション
- CSRF対策（Double Submit Cookie）
- Rate Limiting（1分あたり100リクエスト）

### 6.2 入力検証
- SQLインジェクション対策（プリペアドステートメント）
- XSS対策（出力エスケープ）
- ファイルアップロード制限（将来実装時）

### 6.3 通信セキュリティ
- HTTPS必須
- CORS設定の適切な実装
- セキュリティヘッダーの設定

## 7. パフォーマンス要件

### 7.1 フロントエンド
- First Contentful Paint: < 1.5秒
- Time to Interactive: < 3秒
- 画像の遅延読み込み
- コード分割とダイナミックインポート

### 7.2 バックエンド
- API応答時間: < 200ms（95パーセンタイル）
- データベースクエリ最適化
- キャッシング戦略（Redis検討）

## 8. テスト戦略

### 8.1 フロントエンド
- 単体テスト: Vitest
- 統合テスト: Testing Library
- E2Eテスト: Playwright

### 8.2 バックエンド
- 単体テスト: Go標準のtestingパッケージ
- 統合テスト: testcontainers-go
- APIテスト: Postman/Newman

## 9. CI/CD設定

### 9.1 GitHub Actions
```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Node.js
        uses: actions/setup-node@v3
      - name: Install dependencies
        run: cd frontend && npm ci
      - name: Run tests
        run: cd frontend && npm test
      - name: Build
        run: cd frontend && npm run build

  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v4
      - name: Run tests
        run: cd backend && go test ./...
      - name: Build
        run: cd backend && go build -o api cmd/api/main.go
```

## 10. 環境変数設定

### 10.1 フロントエンド (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-secret-key
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

### 10.2 バックエンド (.env)
```env
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=readlater_app
DB_PASSWORD=secure_password
DB_NAME=readlater_db

# JWT
JWT_SECRET=your-jwt-secret
JWT_EXPIRY=7d

# AI Services
GROQ_API_KEY=your-groq-key
ANTHROPIC_API_KEY=your-anthropic-key

# OAuth
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

## 11. 実装時の注意事項

### 11.1 エラーハンドリング
- すべてのAPIエラーは統一フォーマットで返す
- ユーザーフレンドリーなエラーメッセージ
- 適切なログ記録

### 11.2 国際化対応（将来）
- 日付フォーマットの考慮
- タイムゾーン対応
- 多言語対応の準備

### 11.3 アクセシビリティ
- セマンティックHTML
- 適切なARIA属性
- キーボードナビゲーション対応

## 12. 開発開始時のコマンド

### 12.1 初期セットアップ
```bash
# リポジトリのクローン
git clone https://github.com/your-username/readlater-app.git
cd readlater-app

# Dockerでの開発環境起動
docker-compose up -d

# フロントエンドのセットアップ
cd frontend
npm install
npm run dev

# バックエンドのセットアップ
cd ../backend
go mod download
air
```

### 12.2 データベースマイグレーション
```bash
cd backend
migrate -path migrations -database "mysql://user:password@tcp(localhost:3306)/readlater_db" up
```

## 13. デプロイ手順

### 13.1 フロントエンド（Vercel）
1. GitHubリポジトリをVercelに接続
2. ビルド設定：
   - Framework Preset: Next.js
   - Build Command: `cd frontend && npm run build`
   - Output Directory: `frontend/.next`

### 13.2 バックエンド（Railway/Fly.io）
1. Dockerfileの作成
2. 環境変数の設定
3. データベース接続の設定

## 14. 追加実装ガイドライン

### 14.1 コーディング規約
- Go: Effective Goに準拠
- TypeScript: Airbnb Style Guide
- コミットメッセージ: Conventional Commits

### 14.2 ドキュメント
- API: OpenAPI 3.0仕様
- コード: JSDoc/GoDoc
- README: セットアップ手順を明記

この実装指示書に基づいて、Phase 1のMVP機能から順次実装を進めてください。