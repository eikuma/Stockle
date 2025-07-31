# P1-103 & P1-203 並列実行指示書

## 🎯 概要

記事管理機能のフロントエンド（P1-103）とバックエンドAPI（P1-203）を効率的に並列開発するためのガイドです。

## 👥 担当者

| チケット | 担当者 | 主な実装内容 |
|---------|--------|-------------|
| **P1-103** | Member 1 (Frontend) | 記事管理UI（保存フォーム、一覧、詳細、検索） |
| **P1-203** | Member 2 (Backend Infrastructure) | 記事管理API（CRUD、検索、スクレイピング） |

## 🚀 並列実行戦略

### Phase 1: 独立開発フェーズ（Day 1-5）

両メンバーが**完全に独立**して作業可能な部分を先行実装：

#### Member 1 (Frontend) - 優先実装順序
1. **記事データ型定義** (`types/article.ts`)
   - Article, Category, Tag型の定義
   - API レスポンス型の定義

2. **UIコンポーネント実装**
   - `ArticleCard` コンポーネント
   - `SaveArticleDialog` コンポーネント
   - `SearchBar` コンポーネント
   - `ArticleFilters` コンポーネント

3. **モックデータでのUI確認**
   - 静的データでのUI動作確認
   - レスポンシブデザイン調整

#### Member 2 (Backend) - 優先実装順序
1. **データモデル実装**
   - Article, Category, Tag モデルの拡張
   - データベースマイグレーション

2. **Webスクレイピングサービス**
   - `ScraperService` 実装
   - メタデータ抽出機能

3. **リポジトリ層実装**
   - `ArticleRepository` 完全実装
   - 検索・フィルタリング機能

### Phase 2: API仕様調整フェーズ（Day 6-7）

#### 共同作業事項
1. **API仕様の確定**
   - エンドポイント定義の最終確認
   - レスポンス形式の統一
   - エラーハンドリングの統一

2. **型定義の同期**
   - TypeScript型とGo構造体の整合性確認
   - フィールド名の統一

### Phase 3: 統合開発フェーズ（Day 8-14）

#### Member 1 (Frontend)
- `services/articles.ts` APIクライアント実装
- `hooks/useArticles.ts` カスタムフック実装
- 実際のAPIとの接続テスト

#### Member 2 (Backend)
- `ArticleController` 完全実装
- エラーハンドリングの充実
- パフォーマンス最適化

## 📋 並列実行スケジュール

### Week 1: 独立開発期間

| 日 | Member 1 (Frontend) | Member 2 (Backend) | 同期事項 |
|----|--------------------|--------------------|----------|
| **Day 1** | 型定義・ArticleCard実装 | データモデル・マイグレーション | API仕様書レビュー |
| **Day 2** | SaveArticleDialog実装 | ScraperService実装 | - |
| **Day 3** | SearchBar・Filters実装 | ArticleRepository実装 | - |
| **Day 4** | モックデータでのUI確認 | 検索・フィルタリング実装 | - |
| **Day 5** | レスポンシブ調整 | コントローラー基盤実装 | API仕様調整会議 |

### Week 2: 統合開発期間

| 日 | Member 1 (Frontend) | Member 2 (Backend) | 統合作業 |
|----|--------------------|--------------------|----------|
| **Day 8** | APIクライアント実装 | 記事保存API完成 | 保存機能テスト |
| **Day 9** | 記事一覧機能完成 | 記事取得API完成 | 一覧機能テスト |
| **Day 10** | 検索機能完成 | 検索API完成 | 検索機能テスト |
| **Day 11** | 記事詳細機能完成 | 更新・削除API完成 | CRUD機能テスト |
| **Day 12** | エラーハンドリング | パフォーマンス最適化 | 統合テスト |
| **Day 13** | UI/UX調整 | API最適化 | E2Eテスト |
| **Day 14** | 最終調整 | 最終調整 | **完成** |

## 🔄 Git Worktree 並列開発設定

### 初期セットアップ
```bash
# 最新のmainブランチを取得
git checkout main && git pull origin main

# 各メンバーのworktreeを作成
git worktree add -b feature/phase1-frontend worktree-frontend          # Member 1
git worktree add -b feature/phase1-backend-infrastructure worktree-backend-infrastructure  # Member 2
```

### 日次同期コマンド
```bash
# 各メンバーのworktreeで実行
git fetch origin
git rebase origin/main

# 進捗の共有
git log --oneline -5
```

## 📡 API仕様書（確定版）

### 1. 記事保存 API
```
POST /api/articles
Authorization: Bearer <token>
Content-Type: application/json

{
  "url": "https://example.com/article",
  "categoryId": "uuid-string", // optional
  "tags": ["AI", "技術"] // optional
}

Response 201:
{
  "message": "Article saved successfully",
  "article": {
    "id": "uuid",
    "title": "記事タイトル",
    "url": "https://example.com/article",
    "thumbnailUrl": "https://example.com/thumb.jpg",
    "summary": "記事の要約...",
    "author": "著者名",
    "siteName": "サイト名",
    "status": "unread",
    "savedAt": "2024-01-15T10:00:00Z",
    "category": { "id": "uuid", "name": "カテゴリ名", "color": "#6B7280" },
    "tags": [{ "id": "uuid", "name": "AI" }]
  }
}
```

### 2. 記事一覧取得 API
```
GET /api/articles?page=1&limit=20&status=unread&category_id=uuid&search=keyword
Authorization: Bearer <token>

Response 200:
{
  "articles": [Article...],
  "total": 150,
  "page": 1,
  "limit": 20
}
```

### 3. 記事詳細取得 API
```
GET /api/articles/:id
Authorization: Bearer <token>

Response 200:
{
  "article": {
    "id": "uuid",
    "title": "記事タイトル",
    "content": "記事の全文...",
    "summary": "記事の要約...",
    // ... その他のフィールド
  }
}
```

### 4. 記事更新 API
```
PATCH /api/articles/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "read",        // optional
  "isFavorite": true,     // optional
  "categoryId": "uuid",   // optional
  "readingProgress": 0.75 // optional
}

Response 200:
{
  "message": "Article updated successfully"
}
```

### 5. 記事削除 API
```
DELETE /api/articles/:id
Authorization: Bearer <token>

Response 200:
{
  "message": "Article deleted successfully"
}
```

## 📊 TypeScript型定義（共有）

```typescript
// types/article.ts
export interface Article {
  id: string;
  userId: string;
  categoryId?: string;
  url: string;
  title: string;
  content?: string;
  summary?: string;
  summaryShort?: string;
  summaryLong?: string;
  thumbnailUrl?: string;
  author?: string;
  siteName?: string;
  publishedAt?: string;
  savedAt: string;
  lastAccessedAt?: string;
  status: 'unread' | 'read' | 'archived';
  isFavorite: boolean;
  readingProgress: number;
  readingTimeSeconds: number;
  wordCount?: number;
  language: string;
  category?: Category;
  tags: Tag[];
}

export interface Category {
  id: string;
  userId: string;
  name: string;
  color: string;
  displayOrder: number;
  isDefault: boolean;
  articleCount: number;
  articles?: Article[];
}

export interface Tag {
  id: string;
  userId: string;
  name: string;
  usageCount: number;
}

export interface SaveArticleForm {
  url: string;
  categoryId?: string;
  tags?: string[];
}

export interface ArticleFilters {
  status?: 'unread' | 'read' | 'archived';
  categoryId?: string;
  search?: string;
  page?: number;
  limit?: number;
  favorite?: boolean;
}
```

## 🛠 開発環境セットアップ

### Member 1 (Frontend)
```bash
cd worktree-frontend/frontend
npm install
npm run dev  # http://localhost:3000

# 開発時のモックAPI利用
export NEXT_PUBLIC_API_URL=http://localhost:3001  # モックサーバー
```

### Member 2 (Backend)
```bash
cd worktree-backend-infrastructure/backend
go mod download
air  # ホットリロード http://localhost:8080

# データベース準備
docker-compose up -d mysql
migrate -path migrations -database "mysql://user:password@tcp(localhost:3306)/stockle_db" up
```

## 🧪 テスト戦略

### Member 1 (Frontend) テスト
- **Unit Tests**: コンポーネントテスト（Vitest + Testing Library）
- **Integration Tests**: APIクライアント・フック統合テスト
- **E2E Tests**: 主要ユーザーフロー（Playwright）

### Member 2 (Backend) テスト
- **Unit Tests**: サービス・リポジトリテスト
- **Integration Tests**: API エンドポイントテスト
- **Performance Tests**: スクレイピング・検索パフォーマンス

### 統合テスト（共同）
- **API Contract Tests**: 型安全性の確認
- **End-to-End Tests**: 記事保存 → 一覧 → 詳細フロー
- **Error Handling Tests**: エラーケースの確認

## 🚨 並列実行時の注意点

### 1. API仕様の変更管理
- **ルール**: API仕様変更は必ず両メンバーに事前共有
- **ツール**: OpenAPI仕様書を共有リポジトリで管理
- **更新頻度**: 仕様変更は1日1回までに制限

### 2. 型定義の同期
- **共有場所**: `types/article.ts` をソースオブトゥルース
- **更新手順**: TypeScript型を更新 → Go構造体に反映
- **検証**: 統合テストで型の整合性を自動確認

### 3. モックデータの活用
- **Frontend**: 開発初期はモックデータでUI確認
- **Backend**: PostmanコレクションでAPI確認
- **同期**: モックデータは実際のAPI仕様に合わせて作成

### 4. エラーハンドリングの統一
- **HTTPステータスコード**: 統一された使い方
- **エラーレスポンス形式**: 共通のエラー形式
- **ユーザー向けメッセージ**: フロントエンドで適切な表示

## 📈 成功指標

### 開発効率
- [ ] **並列開発期間**: 単独開発比50%短縮（14日 → 7日）
- [ ] **統合時間**: 2日以内で完全統合
- [ ] **バグ発生率**: 統合時の致命的バグ0件

### 品質指標
- [ ] **API応答時間**: 95%タイル 200ms以内
- [ ] **フロントエンド性能**: TTI 3秒以内
- [ ] **テストカバレッジ**: 80%以上

### ユーザビリティ
- [ ] **記事保存**: 3ステップ以内で完了
- [ ] **検索レスポンス**: 1秒以内で結果表示
- [ ] **モバイル対応**: 完全レスポンシブ

## 🎯 Phase 3 統合後の次ステップ

1. **パフォーマンス最適化**
   - API キャッシュ戦略
   - フロントエンド最適化
   - データベースインデックス最適化

2. **ユーザビリティ向上**
   - アクセシビリティ対応
   - ローディング状態の改善
   - エラー通知の最適化

3. **Phase 2 機能への準備**
   - 音声機能インターフェース設計
   - ポッドキャスト生成API設計

---

このガイドに従うことで、P1-103とP1-203を効率的に並列開発し、高品質な記事管理機能を短期間で完成させることができるっしょ！✨