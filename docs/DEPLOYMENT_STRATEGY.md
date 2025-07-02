# Stockle デプロイメント戦略

## 概要
Stockleの安全で効率的なデプロイメント戦略を定義し、本番環境への継続的デリバリーを実現する。

## 🏗️ デプロイメント環境

### 1. 環境構成

| 環境 | 目的 | URL | デプロイタイミング |
|------|------|-----|-------------------|
| **Development** | 開発・機能テスト | `localhost` | 各開発者のローカル環境 |
| **Staging** | 統合テスト・受入テスト | `https://staging.stockle.app` | mainブランチへのマージ |
| **Production** | 本番サービス | `https://stockle.app` | リリースタグ作成時 |

### 2. インフラ構成

```mermaid
graph TB
    subgraph "Production Environment"
        LB[Load Balancer<br/>ALB]
        
        subgraph "Frontend"
            FE1[Next.js<br/>Vercel/Netlify]
        end
        
        subgraph "Backend"
            API1[Go API Server<br/>Container 1]
            API2[Go API Server<br/>Container 2]
        end
        
        subgraph "Data Layer"
            DB[(MySQL<br/>RDS)]
            CACHE[(Redis<br/>ElastiCache)]
        end
        
        subgraph "External Services"
            GROQ[Groq API]
            CLAUDE[Claude API]
            OAUTH[Google OAuth]
        end
    end
    
    LB --> FE1
    LB --> API1
    LB --> API2
    API1 --> DB
    API2 --> DB
    API1 --> CACHE
    API2 --> CACHE
    API1 --> GROQ
    API1 --> CLAUDE
    FE1 --> OAUTH
```

## 🚀 デプロイメント戦略

### 1. ブルーグリーンデプロイメント

#### 概要
- **Blue Environment**: 現在の本番環境
- **Green Environment**: 新バージョンの環境
- 瞬時切り替えによるダウンタイム最小化

#### 実装手順
```bash
# 1. Green環境へのデプロイ
docker build -t stockle-api:v1.2.0 .
docker tag stockle-api:v1.2.0 registry.com/stockle-api:v1.2.0
docker push registry.com/stockle-api:v1.2.0

# 2. Green環境での動作確認
kubectl apply -f k8s/green-deployment.yaml
kubectl wait --for=condition=available --timeout=300s deployment/stockle-api-green

# 3. ヘルスチェック実行
./scripts/health-check.sh https://green.stockle.app

# 4. ロードバランサー切り替え
kubectl patch service stockle-api -p '{"spec":{"selector":{"version":"green"}}}'

# 5. Blue環境のクリーンアップ（一定期間後）
kubectl delete deployment stockle-api-blue
```

### 2. カナリアデプロイメント

#### 高リスク変更の段階的展開
```yaml
# k8s/canary-deployment.yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: stockle-api-rollout
spec:
  replicas: 5
  strategy:
    canary:
      steps:
      - setWeight: 10    # 10%のトラフィック
      - pause: {duration: 5m}
      - setWeight: 25    # 25%のトラフィック  
      - pause: {duration: 10m}
      - setWeight: 50    # 50%のトラフィック
      - pause: {duration: 10m}
      - setWeight: 100   # 100%のトラフィック
```

## 📋 デプロイメントプロセス

### 1. 事前準備チェックリスト

#### 開発完了時
- [ ] 全テストが通過している
- [ ] セキュリティスキャンが完了している
- [ ] パフォーマンステストが完了している
- [ ] コードレビューが完了している
- [ ] リリースノートが作成されている

#### デプロイ前確認
- [ ] データベースマイグレーションが準備されている
- [ ] 環境変数の設定が完了している
- [ ] 外部サービスとの接続確認が完了している
- [ ] ロールバック手順が明確である
- [ ] 監視・アラート設定が完了している

### 2. デプロイメントフロー

#### Phase 1: データベースマイグレーション
```bash
#!/bin/bash
# scripts/deploy-migration.sh

echo "Starting database migration..."

# 1. バックアップ作成
mysqldump -h $DB_HOST -u $DB_USER -p$DB_PASSWORD $DB_NAME > backup_$(date +%Y%m%d_%H%M%S).sql

# 2. マイグレーション実行
migrate -path migrations -database "mysql://$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME" up

# 3. マイグレーション確認
if [ $? -eq 0 ]; then
    echo "Migration completed successfully"
else
    echo "Migration failed, rolling back..."
    migrate -path migrations -database "mysql://$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME" down 1
    exit 1
fi
```

#### Phase 2: アプリケーションデプロイ
```yaml
# .github/workflows/deploy.yml
name: Deploy to Production

on:
  push:
    tags:
      - 'v*'

jobs:
  deploy-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Build Docker Image
        run: |
          docker build -t stockle-api:${{ github.ref_name }} ./backend
          
      - name: Push to Registry
        run: |
          echo ${{ secrets.DOCKER_PASSWORD }} | docker login -u ${{ secrets.DOCKER_USERNAME }} --password-stdin
          docker push stockle-api:${{ github.ref_name }}
          
      - name: Deploy to Kubernetes
        run: |
          kubectl set image deployment/stockle-api stockle-api=stockle-api:${{ github.ref_name }}
          kubectl rollout status deployment/stockle-api
          
      - name: Health Check
        run: |
          ./scripts/health-check.sh https://api.stockle.app
          
  deploy-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          
      - name: Install Dependencies
        run: |
          cd frontend
          npm ci
          
      - name: Build
        run: |
          cd frontend
          npm run build
          
      - name: Deploy to Vercel
        run: |
          cd frontend
          vercel --prod --token=${{ secrets.VERCEL_TOKEN }}
```

### 3. デプロイ後確認

#### 自動化された確認
```bash
#!/bin/bash
# scripts/post-deploy-check.sh

echo "Starting post-deployment checks..."

# 1. ヘルスチェック
curl -f https://api.stockle.app/health || exit 1

# 2. 基本機能テスト
./scripts/api-integration-test.sh https://api.stockle.app

# 3. パフォーマンス確認
ab -n 100 -c 10 https://api.stockle.app/health

# 4. エラーログ確認
kubectl logs deployment/stockle-api --tail=50 | grep -i error

# 5. メトリクス確認
curl -s https://monitoring.stockle.app/metrics | grep -E "(response_time|error_rate)"

echo "Post-deployment checks completed successfully"
```

## 🔄 ロールバック戦略

### 1. 自動ロールバック条件

#### エラー率しきい値
```yaml
# monitoring/rollback-rules.yaml
rules:
  - alert: HighErrorRate
    expr: error_rate > 0.05  # 5%以上のエラー率
    for: 2m
    action: rollback
    
  - alert: HighResponseTime
    expr: avg_response_time > 2000  # 2秒以上の応答時間
    for: 5m
    action: rollback
    
  - alert: LowAvailability
    expr: availability < 0.95  # 95%未満の可用性
    for: 1m
    action: rollback
```

### 2. 手動ロールバック手順

#### データベースロールバック
```bash
#!/bin/bash
# scripts/rollback-database.sh

echo "Starting database rollback..."

# 1. 現在のバックアップ作成
mysqldump -h $DB_HOST -u $DB_USER -p$DB_PASSWORD $DB_NAME > pre_rollback_$(date +%Y%m%d_%H%M%S).sql

# 2. 指定バージョンへのロールバック
ROLLBACK_VERSION=${1:-1}
migrate -path migrations -database "mysql://$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME" down $ROLLBACK_VERSION

echo "Database rollback completed"
```

#### アプリケーションロールバック
```bash
#!/bin/bash
# scripts/rollback-application.sh

PREVIOUS_VERSION=${1:-$(git describe --tags --abbrev=0 HEAD~1)}

echo "Rolling back to version: $PREVIOUS_VERSION"

# 1. Kubernetesデプロイメントロールバック
kubectl rollout undo deployment/stockle-api

# 2. ロールバック状況確認
kubectl rollout status deployment/stockle-api

# 3. ヘルスチェック
./scripts/health-check.sh https://api.stockle.app

echo "Application rollback completed"
```

## 📊 監視・アラート

### 1. デプロイメントメトリクス

#### 追跡すべき指標
```yaml
# monitoring/deployment-metrics.yaml
metrics:
  deployment:
    - deployment_frequency    # デプロイ頻度
    - lead_time              # リードタイム
    - deployment_success_rate # デプロイ成功率
    - rollback_frequency     # ロールバック頻度
    
  performance:
    - response_time          # 応答時間
    - throughput            # スループット
    - error_rate            # エラー率
    - availability          # 可用性
    
  business:
    - user_satisfaction     # ユーザー満足度
    - feature_adoption      # 機能採用率
    - conversion_rate       # コンバージョン率
```

### 2. アラート設定

#### Slack通知設定
```yaml
# monitoring/alerts.yaml
channels:
  critical: "#incidents"
  warning: "#engineering"
  info: "#deployments"

alerts:
  deployment_started:
    channel: "#deployments"
    message: "🚀 Deployment started: {{ .Version }}"
    
  deployment_success:
    channel: "#deployments" 
    message: "✅ Deployment successful: {{ .Version }}"
    
  deployment_failed:
    channel: "#incidents"
    message: "❌ Deployment failed: {{ .Version }}\nError: {{ .Error }}"
    
  rollback_triggered:
    channel: "#incidents"
    message: "⚠️ Automatic rollback triggered\nReason: {{ .Reason }}"
```

## 🔐 セキュリティ考慮事項

### 1. デプロイメント時のセキュリティ

#### シークレット管理
```yaml
# k8s/secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: stockle-secrets
type: Opaque
data:
  jwt-secret: <base64-encoded-value>
  db-password: <base64-encoded-value>
  groq-api-key: <base64-encoded-value>
```

#### ネットワークセキュリティ
```yaml
# k8s/network-policy.yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: stockle-network-policy
spec:
  podSelector:
    matchLabels:
      app: stockle-api
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: nginx-ingress
    ports:
    - protocol: TCP
      port: 8080
```

### 2. 脆弱性スキャン

#### デプロイ前セキュリティチェック
```bash
#!/bin/bash
# scripts/security-scan.sh

echo "Starting security scan..."

# 1. Dockerイメージスキャン
trivy image stockle-api:latest

# 2. 依存関係脆弱性チェック
# Go
go list -json -m all | nancy sleuth

# Node.js
cd frontend && npm audit

# 3. 静的コード解析
# Go
gosec ./...

# TypeScript
cd frontend && npm run lint:security

echo "Security scan completed"
```

## 📈 継続的改善

### 1. デプロイメント品質向上

#### メトリクス分析
- デプロイ成功率の改善
- ロールバック頻度の削減
- デプロイ時間の短縮

#### プロセス改善
- 自動化範囲の拡大
- テスト品質の向上
- 監視精度の向上

### 2. DevOps文化の推進

#### チーム教育
- デプロイメントベストプラクティス共有
- インシデント対応訓練
- ツール活用スキル向上

#### 定期レビュー
- 月次デプロイメント振り返り
- インシデント分析・改善
- 新技術・ツールの評価

---

**注意**: このドキュメントは継続的に更新される生きた文書です。チームの経験と技術の進歩に合わせて改善していきます。