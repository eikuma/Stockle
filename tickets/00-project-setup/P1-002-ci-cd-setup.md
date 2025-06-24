# P1-002: CI/CDパイプライン構築

## 概要
GitHub Actionsを使用したCI/CDパイプラインの構築

## 担当者
**PdM (Claude Code)**

## 優先度
**高** - 品質保証とデプロイ自動化のため

## 前提条件
- P1-001: プロジェクト基盤セットアップが完了済み
- プロジェクト構造が確定済み

## 作業内容

### 1. GitHub Actionsワークフローの作成
- [ ] `.github/workflows/` フォルダを作成
- [ ] フロントエンド用CI設定を作成
- [ ] バックエンド用CI設定を作成
- [ ] 統合テスト用CI設定を作成
- [ ] セキュリティスキャン設定を作成

### 2. フロントエンドCI設定
- [ ] Node.js 18のセットアップ
- [ ] npm installの実行
- [ ] ESLintによる静的解析
- [ ] Prettierによるフォーマットチェック
- [ ] TypeScript型チェック
- [ ] Jest/Vitestユニットテスト実行
- [ ] PlaywrightによるE2Eテスト実行
- [ ] ビルドの実行

### 3. バックエンドCI設定
- [ ] Go 1.21のセットアップ
- [ ] go mod downloadの実行
- [ ] go vetによる静的解析
- [ ] golangci-lintによる詳細チェック
- [ ] go testによるユニットテスト実行
- [ ] レースコンディションテスト
- [ ] ビルドの実行

### 4. データベースCI設定
- [ ] MySQL 8.0テストコンテナのセットアップ
- [ ] マイグレーションテストの実行
- [ ] テストデータの投入
- [ ] トリガー・ストアドプロシージャのテスト

### 5. セキュリティスキャン設定
- [ ] Dependabotの設定
- [ ] CodeQLによる脆弱性スキャン
- [ ] Secretsスキャンの設定
- [ ] OWASP ZAP セキュリティテスト

### 6. デプロイメント設定
- [ ] Staging環境へのデプロイ設定
- [ ] Production環境へのデプロイ設定
- [ ] 環境変数の管理
- [ ] ロールバック機能

## 実装詳細

### .github/workflows/frontend-ci.yml
```yaml
name: Frontend CI

on:
  push:
    branches: [main, develop]
    paths: ['frontend/**']
  pull_request:
    branches: [main]
    paths: ['frontend/**']

jobs:
  test:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: ./frontend

    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: '18'
        cache: 'npm'
        cache-dependency-path: './frontend/package-lock.json'
    
    - name: Install dependencies
      run: npm ci
    
    - name: Run ESLint
      run: npm run lint
    
    - name: Run Prettier check
      run: npm run format:check
    
    - name: Run TypeScript check
      run: npm run type-check
    
    - name: Run unit tests
      run: npm test -- --coverage
    
    - name: Build application
      run: npm run build
    
    - name: Run E2E tests
      run: npm run test:e2e

  security:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Run security audit
      run: npm audit --prefix ./frontend
```

### .github/workflows/backend-ci.yml
```yaml
name: Backend CI

on:
  push:
    branches: [main, develop]
    paths: ['backend/**']
  pull_request:
    branches: [main]
    paths: ['backend/**']

jobs:
  test:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: ./backend

    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: testpassword
          MYSQL_DATABASE: readlater_test
        ports:
          - 3306:3306
        options: --health-cmd="mysqladmin ping" --health-interval=10s --health-timeout=5s --health-retries=3

    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run go vet
      run: go vet ./...
    
    - name: Run golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        working-directory: ./backend
    
    - name: Run tests
      run: |
        go test -v -race -coverprofile=coverage.out ./...
        go tool cover -html=coverage.out -o coverage.html
      env:
        DB_HOST: localhost
        DB_PORT: 3306
        DB_USER: root
        DB_PASSWORD: testpassword
        DB_NAME: readlater_test
    
    - name: Build application
      run: go build -o api cmd/api/main.go
```

## 受入条件

### 必須条件
- [ ] PRの作成時に自動でCIが実行される
- [ ] すべてのテストが通らない場合はマージがブロックされる
- [ ] フロントエンドとバックエンドの変更が独立してテストされる
- [ ] セキュリティスキャンが実行される

### 品質条件
- [ ] CIの実行時間が10分以内
- [ ] テストカバレッジが80%以上
- [ ] 静的解析でエラーが0件
- [ ] セキュリティスキャンでHigh/Critical脆弱性が0件

## 推定時間
**12時間** (2-3日)

## 依存関係
- P1-001: プロジェクト基盤セットアップ
- 各メンバーの基盤コードが少なくとも一部完成している必要

## 完了後の次ステップ
1. 各メンバーのPR作成時のCI動作確認
2. 統合テストの実行
3. デプロイメントテストの実行

## 備考
- CIが失敗した場合は優先的に修正すること
- セキュリティアラートは即座に対応すること
- パフォーマンスの問題があれば並列実行を検討すること