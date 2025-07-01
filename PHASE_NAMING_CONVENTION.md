# フェーズ命名規則

## 📋 概要

Stockleプロジェクトでは、開発フェーズごとに異なるブランチ名を使用します。このドキュメントでは、フェーズの命名規則と使用例を定義します。

## 🏷️ フェーズ命名パターン

### 基本形式
```
feature/<phase>-<component>
```

- `<phase>`: 開発フェーズを表す識別子
- `<component>`: チームメンバーの担当領域

### コンポーネント名（固定）
- `integration`: PdM統合ブランチ
- `frontend`: フロントエンド開発
- `backend-infrastructure`: バックエンド基盤開発
- `backend-features`: バックエンド機能開発

## 📊 フェーズ名の例

### 1. 数値ベース
```bash
# Phase 1
feature/phase1-integration
feature/phase1-frontend
feature/phase1-backend-infrastructure
feature/phase1-backend-features
```

### 2. バージョンベース
```bash
# MVP
feature/mvp-integration
feature/mvp-frontend
feature/mvp-backend-infrastructure
feature/mvp-backend-features
```

### 3. 機能ベース
```bash
# 音声機能
feature/voice-integration
feature/voice-frontend
feature/voice-backend-infrastructure
feature/voice-backend-features
```

## 🚀 使用方法

### 環境変数での管理（推奨）
```bash
# フェーズを設定
export PHASE="phase1"

# worktreeを作成
git worktree add -b feature/${PHASE}-integration worktree-integration
```

### チーム全体での初期設定
```bash
# チームメンバー全員が同じフェーズ名を使用
PHASE="mvp"  # チームで合意したフェーズ名

# 各自のworktreeを作成
git worktree add -b feature/${PHASE}-integration worktree-integration          # PdM
git worktree add -b feature/${PHASE}-frontend worktree-frontend               # Member 1
git worktree add -b feature/${PHASE}-backend-infrastructure worktree-backend-infrastructure  # Member 2
git worktree add -b feature/${PHASE}-backend-features worktree-backend-features  # Member 3
```

## 📝 フェーズ移行時の手順

### 1. 現在のフェーズを完了
```bash
# 統合ブランチをマージ
gh pr merge feature/${PHASE}-integration

# worktreeをクリーンアップ
git worktree remove worktree-integration
git worktree remove worktree-frontend
git worktree remove worktree-backend-infrastructure
git worktree remove worktree-backend-features
```

### 2. 新しいフェーズを開始
```bash
# 新しいフェーズ名を設定
export PHASE="phase2"

# 新しいworktreeを作成
git worktree add -b feature/${PHASE}-integration worktree-integration
# ... 以下同様
```

## ⚠️ 注意事項

1. **チーム内での統一** - 開発開始前に全員でフェーズ名を合意
2. **命名の一貫性** - ハイフンで単語を区切る（kebab-case）
3. **ドキュメント化** - 各フェーズの目的と内容を記録

## 📊 フェーズ管理表の例

| フェーズ | 期間 | 主な目標 | ステータス |
|---------|------|----------|------------|
| mvp | 2024/01/01-01/31 | 基本機能実装 | 完了 |
| phase1 | 2024/02/01-02/28 | 認証・記事管理 | 進行中 |
| phase2 | 2024/03/01-03/31 | AI機能強化 | 計画中 |
| voice | 2024/04/01-04/30 | 音声機能追加 | 未着手 |

---

この命名規則に従うことで、複数のフェーズにわたる開発を整理された形で管理できます。