# Stockle - 記事保存アプリケーション

後で読みたい記事を保存し、AI による要約生成機能を持つ Web アプリケーションです。

## 🚀 クイックスタート

### 前提条件

- **Node.js** 18.0.0 以上
- **Go** 1.21 以上
- **Docker** & **Docker Compose**
- **Git**

### 開発環境のセットアップ

1. **リポジトリのクローン**
   ```bash
   git clone <repository-url>
   cd stockle
   ```

2. **環境変数の設定**
   ```bash
   cp .env.example .env
   # .env ファイルを編集して適切な値を設定
   ```

3. **Docker コンテナの起動**
   ```bash
   docker-compose up -d
   ```

4. **フロントエンドの起動**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

5. **バックエンドの起動**
   ```bash
   cd backend
   go mod download
   air  # ホットリロード開発サーバー
   ```

6. **アプリケーションへのアクセス**
   - フロントエンド: http://localhost:3000
   - バックエンド API: http://localhost:8080
   - MySQL: localhost:3306
   - Redis: localhost:6379

## 🏗️ プロジェクト構造

```
stockle/
├── frontend/                  # Next.js フロントエンドアプリケーション
│   ├── src/
│   │   ├── app/              # Next.js App Router
│   │   ├── components/       # 再利用可能なUIコンポーネント
│   │   ├── features/         # 機能別モジュール
│   │   ├── hooks/            # カスタムReactフック
│   │   ├── lib/              # ユーティリティ関数
│   │   ├── services/         # API クライアント
│   │   ├── stores/           # Zustand 状態管理
│   │   └── types/            # TypeScript 型定義
│   └── package.json
│
├── backend/                   # Go API サーバー
│   ├── cmd/api/              # アプリケーションエントリーポイント
│   ├── internal/
│   │   ├── config/           # 設定管理
│   │   ├── controllers/      # HTTP リクエストハンドラー
│   │   ├── middleware/       # HTTP ミドルウェア
│   │   ├── models/           # データベースモデル
│   │   ├── repositories/     # データアクセス層
│   │   ├── services/         # ビジネスロジック
│   │   └── validators/       # リクエストバリデーション
│   ├── migrations/           # データベースマイグレーション
│   ├── pkg/                  # 共有パッケージ
│   └── go.mod
│
├── docs/                     # プロジェクトドキュメント
├── docker-compose.yml        # Docker 開発環境設定
├── .env.example             # 環境変数テンプレート
└── README.md                # このファイル
```

## 🛠️ 開発コマンド

### フロントエンド

```bash
cd frontend

# 開発サーバー起動
npm run dev

# ビルド
npm run build

# 本番サーバー起動
npm start

# リント
npm run lint
npm run lint:fix

# 型チェック
npm run type-check

# フォーマット
npm run format
npm run format:check

# テスト
npm test
```

### バックエンド

```bash
cd backend

# 開発サーバー起動（ホットリロード）
air

# ビルド
go build -o api cmd/api/main.go
# または
make build

# テスト
go test ./...
# または
make test

# テストカバレッジ
make test-coverage

# 依存関係更新
go mod tidy
# または
make deps
```

### データベース

```bash
# マイグレーション実行
migrate -path backend/migrations -database "mysql://readlater_app:secure_password@tcp(localhost:3306)/readlater_db" up

# マイグレーション取り消し
migrate -path backend/migrations -database "mysql://readlater_app:secure_password@tcp(localhost:3306)/readlater_db" down
```

### Docker

```bash
# サービス起動
docker-compose up -d

# ログ確認
docker-compose logs -f

# サービス停止
docker-compose down

# ボリューム含めて削除
docker-compose down -v
```

## 🔧 技術スタック

### フロントエンド
- **Next.js 14** (App Router)
- **TypeScript**
- **Tailwind CSS** + **shadcn/ui** + **Radix UI**
- **Zustand** (状態管理)
- **TanStack Query** (データフェッチング)
- **React Hook Form** + **Zod** (フォーム管理)
- **NextAuth.js** (認証)

### バックエンド
- **Go 1.21+**
- **Gin** (Web フレームワーク)
- **GORM** (ORM)
- **golang-jwt** (JWT認証)
- **go-playground/validator** (バリデーション)
- **Air** (ホットリロード)

### データベース・インフラ
- **MySQL 8.0**
- **Redis**
- **Docker & Docker Compose**
- **golang-migrate** (マイグレーション)

### AI/ML
- **Groq API** (第1選択)
- **Anthropic Claude API** (第2選択)

## 🔐 環境変数

環境変数の設定については `.env.example` ファイルを参照してください。

主要な環境変数：

| 変数名 | 説明 | 例 |
|--------|------|-----|
| `DB_HOST` | データベースホスト | `localhost` |
| `DB_PORT` | データベースポート | `3306` |
| `JWT_SECRET` | JWT署名用秘密鍵 | `your-secret-key` |
| `GROQ_API_KEY` | Groq API キー | `gsk_xxx` |
| `NEXT_PUBLIC_API_URL` | フロントエンド用API URL | `http://localhost:8080` |

## 📊 API ドキュメント

API の詳細仕様は `docs/api-design-doc.md` を参照してください。

## 🧪 テスト

### フロントエンド
- **Vitest** + **Testing Library** (単体・統合テスト)
- **Playwright** (E2Eテスト)

### バックエンド
- **Go 標準テスト** (単体・統合テスト)
- **testcontainers-go** (統合テスト)

## 🚀 デプロイ

デプロイに関する詳細は `docs/deployment.md` を参照してください。

## 🤝 貢献

1. このリポジトリをフォーク
2. 機能ブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## 📝 ライセンス

このプロジェクトは MIT ライセンスの下で公開されています。

## 📞 お問い合わせ

プロジェクトに関する質問や提案は GitHub Issues でお願いします。

## 🔗 関連リンク

- [API設計書](docs/api-design-doc.md)
- [データベース設計書](docs/database-design-doc.md)
- [機能要件定義書](docs/functional-requirements-doc.md)
- [実装指示書](docs/implementation-guide.md)

---

_Last updated: 2025-01-15 - Test PR created_