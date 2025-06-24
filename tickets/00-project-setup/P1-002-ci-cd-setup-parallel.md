# P1-002: CI/CDパイプライン構築（並列実行版）

## 概要
GitHub Actionsを使用したCI/CDパイプラインの構築（作業分担版）

## 担当者と作業分担
**PdM主導 + 各メンバー協力体制**

### 🎯 PdM (Claude Code) 担当作業
**推定時間**: 8時間

#### 1. GitHub Actions基盤構築
- [ ] `.github/workflows/` フォルダを作成
- [ ] ワークフロー実行権限の設定
- [ ] Secretsの設定（DB_PASSWORD, JWT_SECRET等）
- [ ] 統合テスト用ワークフローの作成

#### 2. セキュリティスキャン設定
- [ ] Dependabotの設定（`.github/dependabot.yml`）
- [ ] CodeQLによる脆弱性スキャン設定
- [ ] Secretsスキャンの設定
- [ ] OWASP ZAP セキュリティテスト設定

#### 3. デプロイメント設定
- [ ] Staging環境へのデプロイ設定
- [ ] Production環境へのデプロイ設定（準備）
- [ ] 環境変数の管理方針策定
- [ ] ロールバック機能の設計

#### 4. 統合CI/CDワークフロー
- [ ] Pull Request時の統合テスト
- [ ] マージ時の自動デプロイ設定
- [ ] 通知設定（Slack/Discord連携）

### 🎯 Member 1 (Frontend) 担当作業
**推定時間**: 3時間

#### 1. フロントエンドCI設定
- [ ] `.github/workflows/frontend-ci.yml` の作成
- [ ] Node.js 18のセットアップ設定
- [ ] npm cache設定
- [ ] ESLintによる静的解析設定

#### 2. フロントエンドテスト設定
- [ ] TypeScript型チェック設定
- [ ] Jest/Vitestユニットテスト実行設定
- [ ] ビルドテスト設定
- [ ] E2Eテスト準備（Playwright）

#### 3. フロントエンド品質チェック
- [ ] Prettierによるフォーマットチェック
- [ ] パフォーマンステスト設定
- [ ] アクセシビリティテスト設定

### 🎯 Member 2 (Backend Infrastructure) 担当作業
**推定時間**: 4時間

#### 1. バックエンドCI設定
- [ ] `.github/workflows/backend-ci.yml` の作成
- [ ] Go 1.21のセットアップ設定
- [ ] Go modules cache設定
- [ ] go vetによる静的解析設定

#### 2. バックエンドテスト設定
- [ ] golangci-lintによる詳細チェック設定
- [ ] go testによるユニットテスト実行設定
- [ ] レースコンディションテスト設定
- [ ] ビルドテスト設定

#### 3. データベースCI設定
- [ ] MySQL 8.0テストコンテナのセットアップ
- [ ] マイグレーションテストの設定
- [ ] テストデータ投入設定
- [ ] データベース接続テスト設定

### 🎯 Member 3 (Backend Features) 担当作業
**推定時間**: 2時間

#### 1. AI統合テスト設定
- [ ] AI APIモックテストの設定
- [ ] API レート制限テスト設定
- [ ] フォールバック機能テスト設定

#### 2. 非同期処理テスト設定
- [ ] Redisコンテナテスト設定
- [ ] ジョブキューテスト設定
- [ ] 非同期処理タイムアウトテスト設定

## 📋 各ワークフローファイルの詳細

### Member 1作成: .github/workflows/frontend-ci.yml
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

    strategy:
      matrix:
        node-version: [18, 20]

    steps:
    - name: Checkout code
      uses: actions/checkout@v4
    
    - name: Setup Node.js ${{ matrix.node-version }}
      uses: actions/setup-node@v4
      with:
        node-version: ${{ matrix.node-version }}
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
      run: npm test -- --coverage --watchAll=false
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./frontend/coverage/lcov.info
        flags: frontend
    
    - name: Build application
      run: npm run build
    
    - name: Run E2E tests
      run: npm run test:e2e
      if: matrix.node-version == 18  # E2Eテストは1つのバージョンでのみ実行

  lighthouse:
    runs-on: ubuntu-latest
    needs: test
    if: github.event_name == 'pull_request'
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v4
    
    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: 18
        cache: 'npm'
        cache-dependency-path: './frontend/package-lock.json'
    
    - name: Install dependencies
      working-directory: ./frontend
      run: npm ci
    
    - name: Build application
      working-directory: ./frontend
      run: npm run build
    
    - name: Run Lighthouse CI
      working-directory: ./frontend
      run: npx lhci autorun
```

### Member 2作成: .github/workflows/backend-ci.yml
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

    strategy:
      matrix:
        go-version: [1.21, 1.22]

    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: testpassword
          MYSQL_DATABASE: readlater_test
          MYSQL_USER: readlater_test
          MYSQL_PASSWORD: testpassword
        ports:
          - 3306:3306
        options: >-
          --health-cmd="mysqladmin ping --silent"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=3

      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379
        options: >-
          --health-cmd="redis-cli ping"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=3

    steps:
    - name: Checkout code
      uses: actions/checkout@v4
    
    - name: Setup Go ${{ matrix.go-version }}
      uses: actions/setup-go@v4
      with:
        go-version: ${{ matrix.go-version }}
        cache-dependency-path: './backend/go.sum'
    
    - name: Install dependencies
      run: go mod download
    
    - name: Verify dependencies
      run: go mod verify
    
    - name: Run go vet
      run: go vet ./...
    
    - name: Run golangci-lint
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
        working-directory: ./backend
        args: --timeout=5m
    
    - name: Run tests
      run: |
        go test -v -race -coverprofile=coverage.out ./...
        go tool cover -html=coverage.out -o coverage.html
      env:
        DB_HOST: localhost
        DB_PORT: 3306
        DB_USER: readlater_test
        DB_PASSWORD: testpassword
        DB_NAME: readlater_test
        REDIS_URL: redis://localhost:6379
        JWT_SECRET: test-jwt-secret
        GROQ_API_KEY: test-groq-key
        ANTHROPIC_API_KEY: test-anthropic-key
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./backend/coverage.out
        flags: backend
    
    - name: Build application
      run: go build -v -o stockle-api ./cmd/api
    
    - name: Test binary execution
      run: |
        chmod +x stockle-api
        timeout 10s ./stockle-api || if [ $? -eq 124 ]; then echo "Binary started successfully"; else exit $?; fi

  security:
    runs-on: ubuntu-latest
    steps:
    - name: Checkout code
      uses: actions/checkout@v4
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: 1.21
    
    - name: Run Gosec Security Scanner
      uses: securecodewarrior/github-action-gosec@master
      with:
        args: '-fmt sarif -out gosec.sarif ./...'
        working-directory: ./backend
    
    - name: Upload SARIF file
      uses: github/codeql-action/upload-sarif@v2
      with:
        sarif_file: gosec.sarif
```

### PdM作成: .github/workflows/integration-test.yml
```yaml
name: Integration Test

on:
  pull_request:
    branches: [main]
  push:
    branches: [main]

jobs:
  integration-test:
    runs-on: ubuntu-latest
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v4
    
    - name: Setup Node.js
      uses: actions/setup-node@v4
      with:
        node-version: 18
        cache: 'npm'
        cache-dependency-path: './frontend/package-lock.json'
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: 1.21
        cache-dependency-path: './backend/go.sum'
    
    - name: Start services
      run: |
        docker-compose up -d mysql redis
        sleep 30  # Wait for services to be ready
    
    - name: Install frontend dependencies
      working-directory: ./frontend
      run: npm ci
    
    - name: Install backend dependencies
      working-directory: ./backend
      run: go mod download
    
    - name: Build backend
      working-directory: ./backend
      run: go build -o stockle-api ./cmd/api
    
    - name: Start backend server
      working-directory: ./backend
      run: |
        export DB_HOST=localhost
        export DB_PORT=3306
        export DB_USER=readlater_app
        export DB_PASSWORD=secure_password
        export DB_NAME=readlater_db
        export JWT_SECRET=integration-test-secret
        export GROQ_API_KEY=${{ secrets.GROQ_API_KEY_TEST }}
        export ANTHROPIC_API_KEY=${{ secrets.ANTHROPIC_API_KEY_TEST }}
        ./stockle-api &
        sleep 10
    
    - name: Build frontend
      working-directory: ./frontend
      run: |
        export NEXT_PUBLIC_API_URL=http://localhost:8080
        npm run build
    
    - name: Run integration tests
      run: |
        # API ヘルスチェック
        curl -f http://localhost:8080/api/v1/health
        
        # 統合テストスクリプト実行
        chmod +x ./scripts/integration-test.sh
        ./scripts/integration-test.sh
    
    - name: Stop services
      if: always()
      run: docker-compose down
```

## 🔄 並列実行スケジュール

### Day 1: CI設定作成（各自2-3時間）
**午前**:
- PdM: GitHub Actions基盤 + セキュリティスキャン設定
- Member 1: フロントエンドCI設定作成
- Member 2: バックエンドCI設定作成
- Member 3: AI統合テスト設定作成

**午後**:
- 全員: 各自の設定をプッシュ
- PdM: 統合テスト設定作成
- 全員: CI動作確認とデバッグ

### Day 2: 品質向上・完成（各自1-2時間）
**午前**:
- PdM: デプロイメント設定 + 通知設定
- Member 1: E2E/パフォーマンステスト設定
- Member 2: セキュリティテスト設定
- Member 3: 非同期処理テスト設定

**午後**:
- 全員: CI/CDの最終確認
- 全員: パフォーマンス・品質チェック

## 受入条件

### 必須条件（全員達成必要）
- [ ] PRの作成時に自動でCIが実行される
- [ ] すべてのテストが通らない場合はマージがブロックされる
- [ ] フロントエンドとバックエンドの変更が独立してテストされる
- [ ] セキュリティスキャンが実行される
- [ ] 統合テストが正常に動作する

### 品質条件
- [ ] CIの実行時間が10分以内
- [ ] テストカバレッジが80%以上
- [ ] 静的解析でエラーが0件
- [ ] セキュリティスキャンでHigh/Critical脆弱性が0件

## 推定時間
**合計17時間** (並列実行により2日で完了)
- PdM: 8時間
- Member 1: 3時間
- Member 2: 4時間
- Member 3: 2時間

## 依存関係
- P1-001: プロジェクト基盤セットアップの完了
- 各メンバーの基盤コードが少なくとも一部完成している必要

## 完了後の次ステップ
1. 各メンバーのPR作成時のCI動作確認
2. 統合テストの実行確認
3. デプロイメントテストの実行

## 備考
- CIが失敗した場合は優先的に修正すること
- セキュリティアラートは即座に対応すること
- パフォーマンスの問題があれば並列実行を検討すること
- 各メンバーは自分の担当領域のCI設定に責任を持つ