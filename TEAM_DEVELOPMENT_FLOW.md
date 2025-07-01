# Stockle チーム開発フロー

## 📋 概要

本ドキュメントは、PdM（プロダクトマネージャー）1名とメンバー3名の計4名でStockleプロジェクトを開発する際の、Git Worktreeを活用した並列開発フローを定義します。

## 👥 チーム構成

| 役割 | 担当者 | 責任範囲 | ブランチ名パターン |
|------|--------|----------|------------------|
| **PdM** | Claude Code 1 | プロジェクト管理・統合・最終PR作成 | `feature/<phase>-integration` |
| **Member 1** | Claude Code 2 | フロントエンド開発 | `feature/<phase>-frontend` |
| **Member 2** | Claude Code 3 | バックエンド基盤開発 | `feature/<phase>-backend-infra` |
| **Member 3** | Claude Code 4 | バックエンド機能・AI統合開発 | `feature/<phase>-backend-features` |

※ `<phase>`は開発フェーズに応じて変更（例: phase1, phase2, mvp, v2など）

## 🔄 開発フロー

### 1. 初期セットアップ（全員）

```bash
# メインリポジトリで最新の状態を取得
git checkout main
git pull origin main

# 現在のフェーズを定義（例: phase1, phase2, mvp, etc.）
PHASE="phase1"  # または環境変数として export PHASE=phase1

# 各メンバーが自分のworktreeを作成
# PdM
git worktree add -b feature/${PHASE}-integration worktree-integration

# Member 1 (Frontend)
git worktree add -b feature/${PHASE}-frontend worktree-frontend

# Member 2 (Backend Infrastructure)
git worktree add -b feature/${PHASE}-backend-infra worktree-backend-infra

# Member 3 (Backend Features)
git worktree add -b feature/${PHASE}-backend-features worktree-backend-features
```

### 2. 並列開発フェーズ

各メンバーは自分のworktreeディレクトリで独立して作業します：

```bash
# 例：Member 1の作業
cd worktree-frontend
# フロントエンド開発を実施
npm install
npm run dev
# コミット
git add .
git commit -m "feat(frontend): 実装内容の説明"
git push origin feature/${PHASE}-frontend
```

### 3. 日次同期（推奨）

```bash
# 各メンバーのworktreeで
git fetch origin
git rebase origin/main  # コンフリクトがある場合は解決
```

### 4. 統合フェーズ（PdM主導）

各メンバーの作業が完了したら、PdMが統合作業を実施：

```bash
# PdMのworktreeで作業
cd worktree-integration

# 最新のmainを取得
git fetch origin
git rebase origin/main

# 各メンバーのブランチをマージ
# 1. フロントエンドの統合
git fetch origin feature/${PHASE}-frontend
git merge origin/feature/${PHASE}-frontend -m "feat: フロントエンド実装を統合"

# 2. バックエンド基盤の統合
git fetch origin feature/${PHASE}-backend-infra
git merge origin/feature/${PHASE}-backend-infra -m "feat: バックエンド基盤実装を統合"

# 3. バックエンド機能の統合
git fetch origin feature/${PHASE}-backend-features
git merge origin/feature/${PHASE}-backend-features -m "feat: バックエンド機能・AI統合を統合"

# 統合テストの実施
docker-compose up -d
npm run test:integration
```

### 5. Pull Request作成（PdM）

統合が完了し、全てのテストが通ったら：

```bash
# 統合ブランチをプッシュ
git push origin feature/${PHASE}-integration

# GitHub CLIでPR作成
gh pr create \
  --base main \
  --head feature/${PHASE}-integration \
  --title "feat: ${PHASE} 実装完了" \
  --body "$(cat <<'EOF'
## 概要
${PHASE}の実装を完了しました。

## 実装内容
- フロントエンド: 認証UI、記事管理UI
- バックエンド基盤: Go基盤、認証システム、記事管理API
- バックエンド機能: AI統合、要約生成機能

## テスト結果
- ✅ 単体テスト: 全て合格
- ✅ 統合テスト: 全て合格
- ✅ E2Eテスト: 全て合格

## チェックリスト
- [x] コードレビュー済み
- [x] テスト実施済み
- [x] ドキュメント更新済み
- [x] セキュリティチェック済み

🤖 Generated with Claude Code
EOF
)"
```

## 📝 コミットメッセージ規約

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type
- `feat`: 新機能
- `fix`: バグ修正
- `docs`: ドキュメント
- `style`: フォーマット修正
- `refactor`: リファクタリング
- `test`: テスト追加・修正
- `chore`: ビルドプロセスやツールの変更

### Scope
- `frontend`: フロントエンド関連
- `backend`: バックエンド関連
- `infra`: インフラ・DevOps関連
- `ai`: AI機能関連

### 例
```
feat(frontend): 記事一覧画面のUI実装

- グリッド/リスト表示切り替え機能を追加
- 無限スクロールを実装
- レスポンシブデザイン対応

Resolves: #123
```

## 🚨 コンフリクト解決

### 予防策
1. **定期的な同期**: 日次でorigin/mainをrebase
2. **小さなコミット**: 機能単位で細かくコミット
3. **早期統合**: 完成した機能から順次統合

### 解決手順
```bash
# コンフリクトが発生した場合
git status  # コンフリクトファイルを確認
# エディタでコンフリクトを解決
git add <resolved-files>
git rebase --continue
```

## 🔧 Worktree管理コマンド

```bash
# worktree一覧表示
git worktree list

# worktree削除
git worktree remove worktree-frontend

# 不要なworktreeを整理
git worktree prune
```

## 📊 進捗管理

### 日次スタンドアップ（仮想）
各メンバーは以下を共有：
1. 昨日完了したタスク
2. 今日取り組むタスク
3. ブロッカーや課題

### 週次統合
- 金曜日: 各メンバーの成果物を統合
- 土曜日: 統合テスト実施
- 日曜日: 次週の計画策定

## ⚡ ベストプラクティス

1. **独立性の維持**
   - 他メンバーの作業領域に直接変更を加えない
   - 共通コンポーネントは事前に調整

2. **テストの重要性**
   - 機能実装と同時にテストを書く
   - 統合前に自分の担当範囲のテストを完了

3. **ドキュメント更新**
   - 実装と同時にドキュメントを更新
   - APIの変更は必ずOpenAPI仕様を更新

4. **コミュニケーション**
   - 大きな設計変更は事前に共有
   - ブロッカーは早めに報告

## 🎯 成功の指標

- **開発速度**: 並列作業により開発期間を50%短縮
- **品質**: テストカバレッジ80%以上を維持
- **統合効率**: 統合時のコンフリクトを最小化

---

このフローに従うことで、4人のチームが効率的に並列開発を進め、高品質なプロダクトを迅速に開発できます。