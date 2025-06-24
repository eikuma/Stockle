# Webスクレイピング用Goライブラリ比較・推奨事項

## 概要

Stockleアプリケーションで記事の本文抽出に使用するWebスクレイピングライブラリの比較・推奨事項をまとめます。

## 主要ライブラリ比較（2024年）

### 1. Colly（推奨：第1選択）

#### 特徴
- **パフォーマンス**: 2000リクエスト/秒の処理能力
- **成功率**: 99.5%の高い成功率
- **同期処理**: 組み込み並行処理で効率的
- **メモリ効率**: 最小限のメモリフットプリント
- **使いやすさ**: 直感的で清潔なAPI

#### 適用場面
- 静的コンテンツの記事サイト
- サーバーサイドレンダリングされたページ
- 大量のURLを効率的に処理

#### 制限事項
- JavaScript実行不可（Ajax、動的コンテンツ未対応）

#### インストール
```bash
go get github.com/gocolly/colly/v2
```

#### 使用例
```go
import "github.com/gocolly/colly/v2"

c := colly.NewCollector()
c.OnHTML("article", func(e *colly.HTMLElement) {
    title := e.ChildText("h1")
    content := e.ChildText("p")
})
c.Visit("https://example.com/article")
```

### 2. GoQuery（推奨：第2選択）

#### 特徴
- **jQuery互換**: jQueryライクな構文
- **軽量**: HTMLパースに特化
- **柔軟性**: 他のHTTPクライアントと組み合わせ可能
- **メモリ効率**: シングルスレッド操作で効率的

#### 適用場面
- HTMLの詳細な解析・操作
- net/httpと組み合わせたカスタム実装
- 軽量なスクレイピング処理

#### インストール
```bash
go get github.com/PuerkitoBio/goquery
```

#### 使用例
```go
import (
    "github.com/PuerkitoBio/goquery"
    "net/http"
)

resp, _ := http.Get("https://example.com")
defer resp.Body.Close()
doc, _ := goquery.NewDocumentFromReader(resp.Body)
title := doc.Find("h1").Text()
```

### 3. Rod（JavaScript必要時）

#### 特徴
- **ブラウザ自動化**: DevTools Protocol使用
- **JavaScript対応**: 完全なJS実行環境
- **柔軟なAPI**: Puppeteer類似の使いやすさ
- **動的コンテンツ**: SPAや非同期ロードサイト対応

#### 適用場面
- JavaScript重要なサイト
- React/Vue.js等のSPAサイト
- 動的にコンテンツが読み込まれるサイト

#### 制限事項
- パフォーマンス: Collyより低速
- リソース消費: より多くのメモリ・CPU使用

#### インストール
```bash
go get github.com/go-rod/rod
```

### 4. Chromedp（高度なブラウザ操作）

#### 特徴
- **Chrome DevTools**: Chrome DevTools Protocol使用
- **フルブラウザ**: 完全なブラウザ環境
- **スクリーンショット**: 画面キャプチャ機能
- **複雑な操作**: フォーム入力、クリック等

#### 適用場面
- 認証が必要なサイト
- 複雑なユーザー操作が必要
- スクリーンショット取得

#### 制限事項
- 重い処理: 最もリソース消費量が多い
- 複雑性: セットアップが複雑

## Stockleでの推奨構成

### Phase 1: 基本実装

#### 主要ライブラリ: Colly
```go
// go.mod
require (
    github.com/gocolly/colly/v2 v2.1.0
    github.com/PuerkitoBio/goquery v1.8.1 // バックアップ・解析用
)
```

#### 理由
1. **パフォーマンス**: 大量の記事処理に最適
2. **信頼性**: 高い成功率
3. **効率性**: 並行処理による高速化
4. **シンプル**: 実装・メンテナンスが容易

### Phase 2: 動的サイト対応

#### 追加ライブラリ: Rod
```go
// 動的コンテンツが必要な場合のみ
require github.com/go-rod/rod v0.114.0
```

## 実装戦略

### 1. フォールバック戦略

```go
type ContentExtractor struct {
    collyExtractor   *CollyExtractor
    rodExtractor     *RodExtractor    // 必要時のみ
}

func (ce *ContentExtractor) ExtractContent(url string) (*Article, error) {
    // 1. Collyで試行
    article, err := ce.collyExtractor.Extract(url)
    if err == nil && article.IsValid() {
        return article, nil
    }
    
    // 2. JavaScript必要と判断された場合のみRod使用
    if isJavaScriptRequired(url) {
        return ce.rodExtractor.Extract(url)
    }
    
    return nil, err
}
```

### 2. サイト別の最適化

```go
// サイト固有の抽出ルール
var siteRules = map[string]SiteRule{
    "medium.com": {
        titleSelector:   "h1",
        contentSelector: "article p",
        useColly:        true,
    },
    "qiita.com": {
        titleSelector:   ".it-MdContent h1",
        contentSelector: ".it-MdContent p",
        useColly:        true,
    },
    "spa-heavy-site.com": {
        useRod: true, // JavaScript必須サイト
    },
}
```

### 3. パフォーマンス最適化

```go
// Colly設定例
c := colly.NewCollector(
    colly.Async(true),
)
c.Limit(&colly.LimitRule{
    DomainGlob:  "*",
    Parallelism: 5,           // 同時接続数
    Delay:       1 * time.Second, // リクエスト間隔
})
c.SetRequestTimeout(30 * time.Second)
```

## 実装チェックリスト

### Phase 1実装項目
- [ ] Collyベースの記事抽出機能
- [ ] GoQueryを使ったHTMLパース処理
- [ ] エラーハンドリング・リトライ機構
- [ ] ユーザーエージェント・ヘッダー管理
- [ ] レート制限機能

### Phase 2実装項目
- [ ] Rod統合（JavaScript必要サイト用）
- [ ] サイト別最適化ルール
- [ ] パフォーマンス監視
- [ ] キャッシュ機能

### セキュリティ・配慮事項
- [ ] robots.txt遵守
- [ ] 適切なUser-Agent設定
- [ ] レート制限・アクセス頻度制御
- [ ] エラー処理とログ記録

## まとめ

**Stockleでの推奨構成:**
1. **メイン**: Colly（高性能・安定性）
2. **サブ**: GoQuery（軽量・柔軟性）
3. **将来対応**: Rod（JavaScript必要時）

この構成により、効率的で信頼性の高い記事抽出機能を実現できます。