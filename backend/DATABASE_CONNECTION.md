# データベース接続設定

## 概要
MySQL 8.0データベースへの接続設定と接続文字列の構成について説明します。

## 接続文字列の形式

### GORM用の接続文字列
```go
dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
    dbUser,     // DB_USER
    dbPassword, // DB_PASSWORD  
    dbHost,     // DB_HOST
    dbPort,     // DB_PORT
    dbName,     // DB_NAME
)
```

### 環境変数から構成する例
```go
package config

import (
    "fmt"
    "os"
)

func GetDatabaseDSN() string {
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    
    return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        dbUser, dbPassword, dbHost, dbPort, dbName)
}
```

## Docker Compose設定との対応

Docker Composeで定義されているMySQL設定:
```yaml
mysql:
  image: mysql:8.0
  environment:
    MYSQL_ROOT_PASSWORD: rootpassword
    MYSQL_DATABASE: readlater_db
    MYSQL_USER: readlater_app
    MYSQL_PASSWORD: secure_password
  ports:
    - "3306:3306"
```

対応する環境変数:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=readlater_app
DB_PASSWORD=secure_password
DB_NAME=readlater_db
```

## 接続オプション
- `charset=utf8mb4`: UTF-8文字セット（絵文字対応）
- `parseTime=True`: TIME/DATE型を自動的にtime.Timeに変換
- `loc=Local`: タイムゾーンをローカルに設定

## セキュリティ考慮事項
- パスワードは環境変数で管理
- 本番環境では強力なパスワードを使用
- SSL/TLS接続の有効化（本番環境）