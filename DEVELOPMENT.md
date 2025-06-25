# 開発ガイドライン

## Git 運用ルール

### ブランチ戦略

当プロジェクトでは **Git Flow** を採用しています。

#### ブランチの種類

1. **main** - 本番環境へのリリース用ブランチ
2. **develop** - 開発統合ブランチ
3. **feature/** - 機能開発ブランチ
4. **release/** - リリース準備ブランチ
5. **hotfix/** - 緊急修正ブランチ

#### ブランチ命名規則

```
feature/P1-001-project-setup     # チケット番号を含む
feature/user-authentication      # 機能名
bugfix/fix-login-error          # バグ修正
hotfix/security-patch           # 緊急修正
```

### コミットメッセージ規約

[Conventional Commits](https://www.conventionalcommits.org/) 形式を採用。

#### 基本フォーマット

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

#### Type の種類

- **feat**: 新機能
- **fix**: バグ修正
- **docs**: ドキュメント更新
- **style**: コードフォーマット（機能に影響しない変更）
- **refactor**: リファクタリング
- **perf**: パフォーマンス改善
- **test**: テスト追加・修正
- **chore**: ビルドプロセスやツール変更

#### 例

```bash
feat(auth): Google OAuth2認証機能を追加

- NextAuth.jsを使用したGoogle認証を実装
- ユーザー情報をデータベースに保存
- セッション管理機能を追加

Closes #123
```

### プルリクエスト運用

#### 作成前チェックリスト

- [ ] 最新の develop ブランチから branch を作成
- [ ] コミットメッセージが規約に準拠
- [ ] テストが通ることを確認
- [ ] コードレビュー用の説明を記載

#### レビュー基準

1. **機能要件** - 要求された機能が正しく実装されている
2. **コード品質** - 読みやすく保守しやすいコード
3. **テスト** - 適切なテストカバレッジ
4. **セキュリティ** - セキュリティ上の問題がない
5. **パフォーマンス** - パフォーマンスに悪影響がない

#### マージ条件

- [ ] 2名以上のレビュー承認
- [ ] すべてのCIチェックが通過
- [ ] コンフリクトが解決済み
- [ ] 関連するIssueがクローズされる

## コーディング規約

### TypeScript/JavaScript (Frontend)

#### 命名規則

```typescript
// Variables & Functions: camelCase
const userName = 'john';
const getUserData = () => {};

// Components: PascalCase
const UserProfile = () => {};

// Constants: UPPER_SNAKE_CASE
const API_BASE_URL = 'https://api.example.com';

// Types & Interfaces: PascalCase
interface UserData {
  id: number;
  name: string;
}
```

#### ファイル命名規則

```
components/UserProfile.tsx       # React Component
hooks/useUserData.ts            # Custom Hook
services/userApi.ts             # API Service
types/user.types.ts             # Type Definitions
utils/dateHelper.ts             # Utility Functions
```

#### Import順序

```typescript
// 1. Node modules
import React from 'react';
import { NextPage } from 'next';

// 2. Internal modules (absolute paths)
import { Button } from '@/components/ui/Button';
import { useAuth } from '@/hooks/useAuth';

// 3. Relative imports
import './UserProfile.styles.css';
```

### Go (Backend)

#### 命名規則

```go
// Packages: lowercase
package userservice

// Functions & Variables: camelCase (exported: PascalCase)
func getUserByID(id int) (*User, error) {}
func GetAllUsers() ([]User, error) {}

// Constants: PascalCase or UPPER_SNAKE_CASE
const DefaultTimeout = 30 * time.Second
const MAX_RETRY_COUNT = 3

// Structs: PascalCase
type UserService struct {
    db *gorm.DB
}
```

#### ファイル命名規則

```
user_service.go              # Service
user_controller.go           # Controller
user_repository.go           # Repository
user_model.go               # Model
user_service_test.go        # Test
```

#### エラーハンドリング

```go
// Error wrapping
func (s *UserService) GetUser(id int) (*User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return user, nil
}

// Custom errors
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
)
```

## テスト戦略

### フロントエンド

#### ユニットテスト（Vitest + Testing Library）

```typescript
// components/__tests__/Button.test.tsx
import { render, screen } from '@testing-library/react';
import { Button } from '../Button';

describe('Button', () => {
  it('renders with correct text', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByText('Click me')).toBeInTheDocument();
  });
});
```

#### E2Eテスト（Playwright）

```typescript
// e2e/auth.spec.ts
import { test, expect } from '@playwright/test';

test('user can login with Google', async ({ page }) => {
  await page.goto('/login');
  await page.click('[data-testid="google-login"]');
  await expect(page).toHaveURL('/dashboard');
});
```

### バックエンド

#### ユニットテスト

```go
// services/user_service_test.go
func TestUserService_GetUser(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    service := NewUserService(db)
    
    // Test
    user, err := service.GetUser(1)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "John Doe", user.Name)
}
```

#### 統合テスト（testcontainers）

```go
func TestUserController_Integration(t *testing.T) {
    // Setup test container
    ctx := context.Background()
    mysqlContainer, err := mysql.RunContainer(ctx)
    require.NoError(t, err)
    defer mysqlContainer.Terminate(ctx)
    
    // Test with real database
    // ...
}
```

## セキュリティガイドライン

### 1. 認証・認可

- JWT トークンは httpOnly Cookie で管理
- リフレッシュトークンローテーション実装
- RBAC（Role-Based Access Control）採用

### 2. 入力検証

```go
// バックエンド
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2,max=100"`
    Email string `json:"email" validate:"required,email"`
}
```

```typescript
// フロントエンド
const schema = z.object({
  name: z.string().min(2).max(100),
  email: z.string().email(),
});
```

### 3. SQL インジェクション対策

```go
// Good: GORM使用でSQLインジェクション対策済み
user := &User{}
db.Where("email = ?", email).First(user)

// Bad: 生のSQL文字列結合
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
```

### 4. XSS対策

```typescript
// Good: エスケープ処理
<div>{escapeHtml(userInput)}</div>

// Bad: 生のHTML挿入
<div dangerouslySetInnerHTML={{__html: userInput}} />
```

### 5. CSRF対策

- Double Submit Cookie パターン実装
- SameSite Cookie属性設定

## パフォーマンス指標

### フロントエンド目標値

- **First Contentful Paint (FCP)**: < 1.5秒
- **Time to Interactive (TTI)**: < 3秒
- **Cumulative Layout Shift (CLS)**: < 0.1

### バックエンド目標値

- **API応答時間**: < 200ms (95パーセンタイル)
- **データベースクエリ**: < 100ms
- **メモリ使用量**: < 512MB

## CI/CD パイプライン

### ブランチ別実行内容

#### feature/* ブランチ

```yaml
- Lint & Format Check
- Unit Tests
- Type Check (Frontend)
- Security Scan
```

#### develop ブランチ

```yaml
- All feature branch checks
- Integration Tests
- E2E Tests (subset)
- Build & Deploy to Staging
```

#### main ブランチ

```yaml
- All develop branch checks
- Full E2E Test Suite
- Performance Tests
- Build & Deploy to Production
```

## 開発ツール設定

### VSCode 推奨拡張機能

```json
{
  "recommendations": [
    "golang.go",
    "bradlc.vscode-tailwindcss",
    "esbenp.prettier-vscode",
    "ms-vscode.vscode-typescript-next",
    "ms-playwright.playwright"
  ]
}
```

### 共通設定

```json
{
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  },
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint"
}
```

## トラブルシューティング

### よくある問題

#### 1. Docker コンテナが起動しない

```bash
# ポート衝突確認
lsof -i :3306
lsof -i :6379

# ボリューム削除
docker-compose down -v
docker-compose up -d
```

#### 2. Node.js依存関係エラー

```bash
# キャッシュクリア
rm -rf node_modules package-lock.json
npm install
```

#### 3. Go モジュールエラー

```bash
# モジュール整理
go mod tidy
go mod verify
```

## 参考資料

- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [React Best Practices](https://react.dev/learn)
- [Next.js Documentation](https://nextjs.org/docs)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)