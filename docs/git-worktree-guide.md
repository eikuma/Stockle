# Git Worktree 運用ガイド

## 🌳 Git Worktreeとは

Git Worktreeは、同一リポジトリから複数の作業ディレクトリを作成できる機能です。各ディレクトリで異なるブランチを同時に扱えるため、並列開発に最適です。

## 💡 なぜWorktreeを使うのか

### 従来の問題
- ブランチ切り替え時にnode_modulesの再インストールが必要
- 実行中のサーバーを停止する必要がある
- 未コミットの変更をstashする手間

### Worktreeのメリット
- ✅ 各メンバーが独立した環境で開発
- ✅ ブランチ切り替えなしで並列作業
- ✅ 依存関係やビルド成果物が混在しない
- ✅ 統合作業が容易

## 📁 ディレクトリ構造

```
Stockle/                          # メインリポジトリ（main branch）
├── worktree-integration/         # PdM用（統合ブランチ）
├── worktree-frontend/           # Frontend開発用
├── worktree-backend-infrastructure/      # Backend基盤開発用
└── worktree-backend-features/   # Backend機能開発用
```

## 🚀 基本的な使い方

### 1. Worktreeの作成

```bash
# 新しいブランチと同時にworktreeを作成
git worktree add -b <branch-name> <directory-name>

# 例：フロントエンド開発用
git worktree add -b feature/phase1-frontend worktree-frontend
```

### 2. 既存ブランチからWorktreeを作成

```bash
# 既存のブランチからworktreeを作成
git worktree add <directory-name> <branch-name>

# 例：
git worktree add worktree-hotfix hotfix/critical-bug
```

### 3. Worktreeの一覧表示

```bash
git worktree list

# 出力例：
# /path/to/Stockle                          80f773f [main]
# /path/to/Stockle/worktree-frontend        2026472 [feature/phase1-frontend]
# /path/to/Stockle/worktree-backend-infrastructure   ef2c05f [feature/phase1-backend-infrastructure]
```

### 4. Worktreeの削除

```bash
# worktreeを削除
git worktree remove <directory-name>

# 例：
git worktree remove worktree-frontend

# 強制削除（未コミットの変更がある場合）
git worktree remove --force worktree-frontend
```

### 5. 不要なWorktreeの整理

```bash
# 削除されたworktreeの参照をクリーンアップ
git worktree prune
```

## 📋 チーム開発での実践例

### Phase 1: セットアップ（全員）

```bash
# 1. 最新のmainを取得
git checkout main
git pull origin main

# 2. 各自のworktreeを作成
# PdM
git worktree add -b feature/phase1-integration worktree-integration

# Member 1
git worktree add -b feature/phase1-frontend worktree-frontend

# Member 2
git worktree add -b feature/phase1-backend-infrastructure worktree-backend-infrastructure

# Member 3
git worktree add -b feature/phase1-backend-features worktree-backend-features
```

### Phase 2: 並列開発

```bash
# Member 1: フロントエンド開発
cd worktree-frontend
npm install
npm run dev
# http://localhost:3000 で開発

# Member 2: バックエンド開発
cd ../worktree-backend-infrastructure
cd backend
go mod download
air
# http://localhost:8080 で開発

# 各自が独立してコミット・プッシュ
git add .
git commit -m "feat: 機能実装"
git push origin feature/phase1-frontend
```

### Phase 3: 統合作業（PdM）

```bash
# PdMのworktreeで統合
cd worktree-integration

# 各ブランチの変更を取り込む
git fetch origin
git merge origin/feature/phase1-frontend
git merge origin/feature/phase1-backend-infrastructure
git merge origin/feature/phase1-backend-features

# 統合テスト実施
docker-compose up -d
npm run test:e2e
```

## ⚠️ 注意事項

### 1. 同一ブランチの重複防止
```bash
# ❌ エラーになる例
git worktree add worktree-main main
# fatal: 'main' is already checked out at '/path/to/Stockle'
```

### 2. 作業前の確認事項
```bash
# worktreeに移動したら必ず現在のブランチを確認
git branch --show-current
```

### 3. 削除時の注意
```bash
# 未コミットの変更がある場合は警告が出る
# 必要なら先にコミットまたはstash
git stash save "一時保存"
git worktree remove worktree-frontend
```

## 🔍 トラブルシューティング

### Q1: worktreeが削除できない
```bash
# エラー: fatal: 'worktree-frontend' contains modified or untracked files
# 解決法1: 変更をコミット
cd worktree-frontend
git add . && git commit -m "WIP: 作業中"

# 解決法2: 強制削除
git worktree remove --force worktree-frontend
```

### Q2: worktreeのパスがわからなくなった
```bash
# 全worktreeのパスを表示
git worktree list --porcelain
```

### Q3: リモートブランチが見つからない
```bash
# リモートの情報を更新
git fetch origin

# リモートブランチ一覧を確認
git branch -r
```

## 🎯 ベストプラクティス

1. **命名規則の統一**
   - worktreeディレクトリ名: `worktree-<purpose>`
   - ブランチ名: `feature/<phase>-<component>`

2. **定期的な同期**
   ```bash
   # 各worktreeで定期的に実行
   git fetch origin
   git rebase origin/main
   ```

3. **クリーンな状態の維持**
   ```bash
   # 不要になったworktreeは即座に削除
   git worktree remove worktree-old-feature
   git worktree prune
   ```

4. **統合前のチェック**
   ```bash
   # 統合前に各worktreeでテスト実行
   npm test
   go test ./...
   ```

## 📚 参考リンク

- [Git公式ドキュメント - git-worktree](https://git-scm.com/docs/git-worktree)
- [GitHub: Working with Git worktree](https://github.blog/2015-07-29-git-worktree/)

---

このガイドに従うことで、チーム全員が効率的にGit Worktreeを活用した並列開発を実現できます。