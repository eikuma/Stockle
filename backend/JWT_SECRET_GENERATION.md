# JWT秘密鍵の生成方法

## 概要
JWT（JSON Web Token）の署名に使用する秘密鍵の安全な生成方法について説明します。

## 生成方法

### 1. OpenSSLを使用した生成（推奨）
```bash
openssl rand -base64 32
```

### 2. Node.jsを使用した生成
```bash
node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"
```

### 3. Pythonを使用した生成
```python
import secrets
import base64
print(base64.b64encode(secrets.token_bytes(32)).decode())
```

## セキュリティ要件
- **最小長**: 32バイト（256ビット）以上
- **エントロピー**: 暗号学的に安全な乱数生成器を使用
- **保管**: 環境変数として設定し、ソースコードには含めない
- **ローテーション**: 定期的な秘密鍵の更新を推奨

## 設定例
.env ファイルに設定:
```env
JWT_SECRET=your-generated-secret-key-here
JWT_EXPIRY=7d  # 7日間有効
```

## 注意事項
- 本番環境では必ず新しい秘密鍵を生成してください
- 開発環境とは異なる秘密鍵を使用してください  
- 秘密鍵は絶対にリポジトリにコミットしないでください