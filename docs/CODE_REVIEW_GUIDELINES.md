# Stockle コードレビューガイドライン

## 目的
効果的なコードレビューを通じて、コード品質の向上、知識共有、バグの早期発見を実現する。

## 🔍 レビューの観点

### 1. 機能性・正確性
- [ ] **要件の実装**: 機能要件を正確に実装しているか
- [ ] **エッジケースの考慮**: 境界値や異常系の処理が適切か
- [ ] **エラーハンドリング**: 例外処理が適切に実装されているか
- [ ] **データ整合性**: データの一貫性が保たれているか

### 2. セキュリティ
- [ ] **認証・認可**: アクセス制御が適切に実装されているか
- [ ] **入力検証**: ユーザー入力の適切なバリデーション
- [ ] **SQLインジェクション対策**: パラメータ化クエリの使用
- [ ] **XSS対策**: 出力時の適切なエスケープ処理
- [ ] **秘密情報の管理**: API キーなどの機密情報がハードコードされていないか

### 3. パフォーマンス
- [ ] **効率的なアルゴリズム**: 時間計算量・空間計算量の最適化
- [ ] **データベースクエリ**: N+1問題の回避、適切なインデックス使用
- [ ] **メモリ管理**: リソースの適切な解放
- [ ] **キャッシュ戦略**: 適切なキャッシュの実装

### 4. 保守性・可読性
- [ ] **命名規則**: 変数名、関数名、クラス名が適切で理解しやすいか
- [ ] **コメント**: 必要十分なコメントがあるか
- [ ] **コード構造**: 適切な関数分割、モジュール化
- [ ] **DRY原則**: 重複コードの排除

### 5. テスト品質
- [ ] **テストカバレッジ**: 新機能に対する適切なテストの存在
- [ ] **テストの質**: 意味のあるテストケースの実装
- [ ] **モック・スタブ**: 外部依存関係の適切なモック化

## 📋 レビュープロセス

### 1. Pull Request 作成者の責任

#### 作成前チェックリスト
- [ ] 自己レビューを実施済み
- [ ] ローカルでのテスト実行・通過確認
- [ ] コンフリクトの解消
- [ ] 適切なブランチ名とコミットメッセージ

#### PR説明の必須項目
```markdown
## 変更概要
- 何を変更したか

## 変更理由
- なぜこの変更が必要か

## 影響範囲
- どの機能に影響するか

## テスト方法
- どのようにテストしたか

## レビューポイント
- 特に注意して見てほしい箇所

## 関連Issue
- Closes #123
```

### 2. レビュアーの責任

#### レビュー優先度
1. **High**: セキュリティ、機能的な問題
2. **Medium**: パフォーマンス、設計上の問題
3. **Low**: スタイル、命名の改善提案

#### コメントの書き方
```markdown
# Good Examples
💡 **提案**: この処理はXXXパターンを使うとより効率的です
🐛 **バグ**: エラーハンドリングが不十分です
⚠️ **重要**: セキュリティ上のリスクがあります
📚 **学習**: この実装方法について詳しく説明します

# Avoid
- 「だめ」「よくない」などの曖昧な表現
- 代替案のない批判
- 感情的な表現
```

### 3. 段階別レビュー手順

#### Phase 1: クイックレビュー（5分以内）
- [ ] PR説明の確認
- [ ] 変更ファイル数・行数の確認
- [ ] CI/CDパイプラインの状態確認

#### Phase 2: 詳細レビュー（15-30分）
- [ ] ロジックの正確性確認
- [ ] セキュリティチェック
- [ ] テストコードの確認
- [ ] パフォーマンス影響の検討

#### Phase 3: 統合確認（必要に応じて）
- [ ] 他機能への影響確認
- [ ] API仕様変更の影響確認
- [ ] データベース変更の影響確認

## 🎯 技術別ガイドライン

### Go (バックエンド)

#### コード品質
```go
// Good: 明確な関数名と適切なエラーハンドリング
func GetUserByID(ctx context.Context, id string) (*User, error) {
    if id == "" {
        return nil, errors.New("user ID cannot be empty")
    }
    
    user, err := repo.FindUserByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return user, nil
}

// Bad: 曖昧な名前と不適切なエラーハンドリング
func Get(id string) *User {
    user, _ := repo.Find(id) // エラーを無視
    return user
}
```

#### セキュリティ
```go
// Good: SQLインジェクション対策
func GetArticlesByUser(ctx context.Context, userID string) ([]Article, error) {
    query := "SELECT * FROM articles WHERE user_id = ? AND deleted_at IS NULL"
    rows, err := db.QueryContext(ctx, query, userID)
    // ...
}

// Bad: SQLインジェクション脆弱性
func GetArticlesByUser(userID string) ([]Article, error) {
    query := fmt.Sprintf("SELECT * FROM articles WHERE user_id = '%s'", userID)
    rows, err := db.Query(query)
    // ...
}
```

### React/TypeScript (フロントエンド)

#### コンポーネント設計
```tsx
// Good: 適切な型定義と Props 分離
interface ArticleCardProps {
  article: Article;
  onSave: (id: string) => void;
  isLoading?: boolean;
}

export const ArticleCard: React.FC<ArticleCardProps> = ({
  article,
  onSave,
  isLoading = false
}) => {
  // ...
};

// Bad: any型の使用、不明確なProps
export const ArticleCard = (props: any) => {
  // ...
};
```

#### Hooks使用方法
```tsx
// Good: カスタムフックの適切な使用
export const useArticles = (userId: string) => {
  const [articles, setArticles] = useState<Article[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchArticles = async () => {
      try {
        setLoading(true);
        const data = await api.getArticles(userId);
        setArticles(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Unknown error');
      } finally {
        setLoading(false);
      }
    };

    if (userId) {
      fetchArticles();
    }
  }, [userId]);

  return { articles, loading, error };
};
```

## ⚡ レビュー効率化のベストプラクティス

### 1. 小さなPRを心がける
- **推奨**: 1PR = 1機能、変更ファイル数 < 10
- **理由**: レビュー精度向上、早期のフィードバック

### 2. 適切なレビュアー選定
- **コア機能**: 2名以上のレビュー必須
- **UI変更**: デザイナーのレビュー含める
- **API変更**: バックエンド・フロントエンド両チームのレビュー

### 3. 自動化ツールの活用
```yaml
# .github/workflows/code-review.yml
name: Code Review
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run ESLint
        run: npm run lint
      - name: Run Go vet
        run: go vet ./...
  
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run security scan
        run: gosec ./...
```

## 📊 品質メトリクス

### レビュー品質指標
- **レビューターンアラウンド時間**: < 24時間
- **承認までの平均ラウンド数**: < 3回
- **レビュー後のバグ発見率**: < 5%

### 追跡すべき指標
- Pull Request サイズ分布
- レビューコメント数
- レビュー時間
- リワーク率

## 🚨 エスカレーション手順

### 意見の相違が生じた場合
1. **議論**: 技術的根拠に基づく建設的な議論
2. **相談**: チームリードに相談
3. **決定**: アーキテクチャ委員会での決定

### 緊急修正の場合
1. **ホットフィックス**: セキュリティ・重大バグの即座修正
2. **事後レビュー**: 修正後24時間以内の詳細レビュー
3. **改善計画**: 再発防止策の策定

## 📚 継続的改善

### 月次振り返り
- レビュープロセスの効果測定
- ガイドラインの更新
- ツール・自動化の改善

### 学習機会
- コードレビュー勉強会
- 良いレビューコメント事例共有
- 外部ツール・手法の調査

---

**Note**: このガイドラインは生きた文書です。チームの成長と共に継続的に更新していきます。