# 記事保存アプリ - データベース設計書

## 1. 概要

### 1.1 データベース選定
- **RDBMS**: MySQL 8.0以上
- **理由**: 
  - 高速な読み書き性能
  - 豊富な運用実績
  - JSON型のサポート
  - 全文検索機能（FULLTEXT INDEX）
  - 無料で利用可能

### 1.2 設計方針
- 正規化を基本とし、パフォーマンスが必要な箇所は非正規化
- UUIDを主キーとして使用（分散システムへの拡張性）
- created_at, updated_atを全テーブルに付与
- 論理削除（deleted_at）を採用
- 適切なインデックスの設定
- InnoDBストレージエンジンを使用（トランザクション対応）

### 1.3 命名規則
- テーブル名: 複数形のスネークケース（例: articles, users）
- カラム名: スネークケース（例: created_at, user_id）
- インデックス名: idx_テーブル名_カラム名
- 外部キー名: fk_テーブル名_参照テーブル名

### 1.4 文字コード設定
```sql
-- データベース作成時の設定
CREATE DATABASE readlater_db
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;
```

## 2. Phase 1: 基本テーブル設計

### 2.1 users（ユーザー）

```sql
CREATE TABLE users (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255), -- NULLの場合はOAuth認証のみ
    display_name VARCHAR(100) NOT NULL,
    profile_image_url TEXT,
    auth_provider VARCHAR(50) DEFAULT 'email', -- 'email' or 'google'
    google_id VARCHAR(255) UNIQUE,
    email_verified BOOLEAN DEFAULT FALSE,
    email_verification_token VARCHAR(255),
    email_verification_expires_at TIMESTAMP NULL,
    password_reset_token VARCHAR(255),
    password_reset_expires_at TIMESTAMP NULL,
    last_login_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_users_email (email),
    INDEX idx_users_google_id (google_id),
    INDEX idx_users_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 2.2 categories（カテゴリ）

```sql
CREATE TABLE categories (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL,
    color VARCHAR(7) DEFAULT '#6B7280', -- HEXカラーコード
    display_order INT NOT NULL DEFAULT 0,
    is_default BOOLEAN DEFAULT FALSE,
    article_count INT DEFAULT 0, -- 非正規化: パフォーマンスのため
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE KEY uk_user_name (user_id, name),
    INDEX idx_categories_user_id (user_id),
    INDEX idx_categories_deleted_at (deleted_at),
    CONSTRAINT fk_categories_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 2.3 articles（記事）

```sql
CREATE TABLE articles (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    category_id CHAR(36),
    url TEXT NOT NULL,
    url_hash VARCHAR(64) NOT NULL, -- SHA256ハッシュ for 重複チェック
    title TEXT NOT NULL,
    content MEDIUMTEXT, -- 記事本文（プレーンテキスト）
    content_hash VARCHAR(64), -- 内容の変更検知用
    summary TEXT,
    summary_short TEXT, -- Phase 4用: 短縮版要約
    summary_long TEXT, -- Phase 4用: 詳細版要約
    thumbnail_url TEXT,
    author VARCHAR(255),
    site_name VARCHAR(255),
    published_at TIMESTAMP NULL,
    saved_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_accessed_at TIMESTAMP NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'unread', -- 'unread', 'read', 'archived'
    is_favorite BOOLEAN DEFAULT FALSE,
    reading_progress DECIMAL(3,2) DEFAULT 0.00, -- 0.00 to 1.00
    reading_time_seconds INT DEFAULT 0,
    word_count INT,
    language VARCHAR(10) DEFAULT 'ja',
    category_confidence_score DECIMAL(3,2), -- 0.00 to 1.00
    summary_generation_status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'processing', 'completed', 'failed'
    summary_generated_at TIMESTAMP NULL,
    summary_model_version VARCHAR(50),
    metadata JSON, -- OGP情報など柔軟なメタデータ保存用
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    UNIQUE KEY uk_user_url (user_id, url_hash),
    INDEX idx_articles_user_id (user_id),
    INDEX idx_articles_category_id (category_id),
    INDEX idx_articles_status (status),
    INDEX idx_articles_saved_at (saved_at DESC),
    INDEX idx_articles_url_hash (url_hash),
    INDEX idx_articles_deleted_at (deleted_at),
    INDEX idx_articles_user_status (user_id, status),
    FULLTEXT idx_articles_title (title),
    FULLTEXT idx_articles_content (content),
    FULLTEXT idx_articles_summary (summary),
    CONSTRAINT fk_articles_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_articles_categories FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 2.4 tags（タグ）

```sql
CREATE TABLE tags (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    name VARCHAR(50) NOT NULL,
    usage_count INT DEFAULT 0, -- 非正規化: 使用頻度表示用
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_tag (user_id, name),
    INDEX idx_tags_user_id (user_id),
    INDEX idx_tags_name (name),
    CONSTRAINT fk_tags_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 2.5 article_tags（記事-タグ関連）

```sql
CREATE TABLE article_tags (
    article_id CHAR(36) NOT NULL,
    tag_id CHAR(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (article_id, tag_id),
    INDEX idx_article_tags_tag_id (tag_id),
    CONSTRAINT fk_article_tags_articles FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    CONSTRAINT fk_article_tags_tags FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 2.6 user_sessions（セッション管理）

```sql
CREATE TABLE user_sessions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    token_hash VARCHAR(64) NOT NULL UNIQUE, -- リフレッシュトークンのハッシュ
    device_info TEXT,
    ip_address VARCHAR(45), -- IPv4/IPv6対応
    expires_at TIMESTAMP NOT NULL,
    last_used_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_sessions_user_id (user_id),
    INDEX idx_user_sessions_token_hash (token_hash),
    INDEX idx_user_sessions_expires_at (expires_at),
    CONSTRAINT fk_user_sessions_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 2.7 user_preferences（ユーザー設定）

```sql
CREATE TABLE user_preferences (
    user_id CHAR(36) PRIMARY KEY,
    default_category_id CHAR(36),
    items_per_page INT DEFAULT 20,
    default_sort VARCHAR(50) DEFAULT 'saved_at',
    default_order VARCHAR(10) DEFAULT 'desc',
    default_view VARCHAR(20) DEFAULT 'list', -- 'list' or 'grid'
    timezone VARCHAR(50) DEFAULT 'Asia/Tokyo',
    language VARCHAR(10) DEFAULT 'ja',
    theme VARCHAR(20) DEFAULT 'light', -- 'light', 'dark', 'auto'
    email_notifications BOOLEAN DEFAULT TRUE,
    push_notifications BOOLEAN DEFAULT FALSE,
    digest_enabled BOOLEAN DEFAULT FALSE,
    digest_frequency VARCHAR(20), -- 'daily', 'weekly', 'monthly'
    digest_time TIME,
    digest_categories JSON, -- カテゴリIDの配列
    preferences JSON DEFAULT '{}', -- その他の設定用
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_preferences_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_preferences_categories FOREIGN KEY (default_category_id) REFERENCES categories(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## 3. Phase 2: 音声機能テーブル

### 3.1 podcasts（ポッドキャスト）

```sql
CREATE TABLE podcasts (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    article_id CHAR(36) NOT NULL UNIQUE, -- 1記事1ポッドキャスト
    audio_url TEXT NOT NULL,
    audio_storage_path TEXT, -- ストレージ内のパス
    duration_seconds INT NOT NULL,
    file_size_bytes BIGINT NOT NULL,
    voice_type VARCHAR(50) DEFAULT 'female',
    speed DECIMAL(2,1) DEFAULT 1.0,
    language VARCHAR(10) DEFAULT 'ja',
    generation_status VARCHAR(20) DEFAULT 'completed',
    generation_model VARCHAR(50),
    generated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_podcasts_article_id (article_id),
    INDEX idx_podcasts_deleted_at (deleted_at),
    CONSTRAINT fk_podcasts_articles FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 3.2 playlists（プレイリスト）

```sql
CREATE TABLE playlists (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_auto_generated BOOLEAN DEFAULT FALSE,
    auto_generation_rule JSON, -- 自動生成ルール
    total_duration_seconds INT DEFAULT 0,
    article_count INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_playlists_user_id (user_id),
    INDEX idx_playlists_deleted_at (deleted_at),
    CONSTRAINT fk_playlists_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 3.3 playlist_items（プレイリスト項目）

```sql
CREATE TABLE playlist_items (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    playlist_id CHAR(36) NOT NULL,
    article_id CHAR(36) NOT NULL,
    position INT NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_playlist_article (playlist_id, article_id),
    UNIQUE KEY uk_playlist_position (playlist_id, position),
    INDEX idx_playlist_items_article_id (article_id),
    CONSTRAINT fk_playlist_items_playlists FOREIGN KEY (playlist_id) REFERENCES playlists(id) ON DELETE CASCADE,
    CONSTRAINT fk_playlist_items_articles FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 3.4 playback_history（再生履歴）

```sql
CREATE TABLE playback_history (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    podcast_id CHAR(36) NOT NULL,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_position_seconds INT DEFAULT 0,
    completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP NULL,
    total_played_seconds INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_playback_history_user_id (user_id),
    INDEX idx_playback_history_podcast_id (podcast_id),
    INDEX idx_playback_history_started_at (started_at DESC),
    CONSTRAINT fk_playback_history_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_playback_history_podcasts FOREIGN KEY (podcast_id) REFERENCES podcasts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## 4. Phase 3: 分析・統計テーブル

### 4.1 reading_sessions（読書セッション）

```sql
CREATE TABLE reading_sessions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    article_id CHAR(36) NOT NULL,
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP NULL,
    duration_seconds INT,
    scroll_depth DECIMAL(3,2), -- 0.00 to 1.00
    interaction_count INT DEFAULT 0, -- クリック、ハイライト等
    device_type VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_reading_sessions_user_id (user_id),
    INDEX idx_reading_sessions_article_id (article_id),
    INDEX idx_reading_sessions_started_at (started_at DESC),
    CONSTRAINT fk_reading_sessions_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_reading_sessions_articles FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 4.2 article_similarities（記事類似度）

```sql
CREATE TABLE article_similarities (
    article_id_1 CHAR(36) NOT NULL,
    article_id_2 CHAR(36) NOT NULL,
    similarity_score DECIMAL(3,2) NOT NULL, -- 0.00 to 1.00
    calculation_method VARCHAR(50) DEFAULT 'cosine',
    calculated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (article_id_1, article_id_2),
    INDEX idx_article_similarities_article_id_2 (article_id_2),
    INDEX idx_article_similarities_score (similarity_score DESC),
    CONSTRAINT fk_article_similarities_articles_1 FOREIGN KEY (article_id_1) REFERENCES articles(id) ON DELETE CASCADE,
    CONSTRAINT fk_article_similarities_articles_2 FOREIGN KEY (article_id_2) REFERENCES articles(id) ON DELETE CASCADE,
    CONSTRAINT chk_article_order CHECK (article_id_1 < article_id_2) -- 重複防止
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 4.3 article_groups（記事グループ）

```sql
CREATE TABLE article_groups (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    auto_generated BOOLEAN DEFAULT TRUE,
    group_type VARCHAR(50) DEFAULT 'similarity', -- 'similarity', 'manual', 'topic'
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_article_groups_user_id (user_id),
    CONSTRAINT fk_article_groups_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 4.4 article_group_items（グループ項目）

```sql
CREATE TABLE article_group_items (
    group_id CHAR(36) NOT NULL,
    article_id CHAR(36) NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (group_id, article_id),
    INDEX idx_article_group_items_article_id (article_id),
    CONSTRAINT fk_article_group_items_groups FOREIGN KEY (group_id) REFERENCES article_groups(id) ON DELETE CASCADE,
    CONSTRAINT fk_article_group_items_articles FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 4.5 user_statistics（ユーザー統計）

```sql
CREATE TABLE user_statistics (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    date DATE NOT NULL,
    reading_time_seconds INT DEFAULT 0,
    articles_read INT DEFAULT 0,
    articles_saved INT DEFAULT 0,
    podcasts_listened INT DEFAULT 0,
    listening_time_seconds INT DEFAULT 0,
    categories_read JSON DEFAULT '{}', -- {"tech": 5, "business": 3}
    tags_used JSON DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_date (user_id, date),
    INDEX idx_user_statistics_date (date DESC),
    INDEX idx_user_statistics_user_date (user_id, date DESC),
    CONSTRAINT fk_user_statistics_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## 5. Phase 4: 自動化・レコメンデーションテーブル

### 5.1 digest_settings（ダイジェスト設定）

```sql
CREATE TABLE digest_settings (
    user_id CHAR(36) PRIMARY KEY,
    enabled BOOLEAN DEFAULT FALSE,
    frequency VARCHAR(20) NOT NULL, -- 'daily', 'weekly', 'monthly'
    send_time TIME NOT NULL DEFAULT '09:00:00',
    timezone VARCHAR(50) DEFAULT 'Asia/Tokyo',
    last_sent_at TIMESTAMP NULL,
    next_send_at TIMESTAMP NULL,
    max_articles INT DEFAULT 10,
    include_categories JSON,
    exclude_categories JSON,
    include_stats BOOLEAN DEFAULT TRUE,
    include_recommendations BOOLEAN DEFAULT TRUE,
    delivery_method VARCHAR(20) DEFAULT 'email', -- 'email', 'push', 'webhook'
    webhook_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_digest_settings_next_send_at (next_send_at),
    CONSTRAINT fk_digest_settings_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 5.2 digest_history（ダイジェスト履歴）

```sql
CREATE TABLE digest_history (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    article_count INT NOT NULL,
    article_ids JSON,
    delivery_status VARCHAR(20) DEFAULT 'sent', -- 'sent', 'failed', 'bounced'
    opened_at TIMESTAMP NULL,
    click_count INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_digest_history_user_id (user_id),
    INDEX idx_digest_history_sent_at (sent_at DESC),
    CONSTRAINT fk_digest_history_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 5.3 recommendations（レコメンデーション）

```sql
CREATE TABLE recommendations (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    article_id CHAR(36) NOT NULL,
    score DECIMAL(3,2) NOT NULL, -- 0.00 to 1.00
    reason_type VARCHAR(50) NOT NULL, -- 'content_based', 'collaborative', 'trending'
    reason_detail TEXT,
    generated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    interaction VARCHAR(20), -- 'viewed', 'saved', 'dismissed'
    interacted_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_article (user_id, article_id),
    INDEX idx_recommendations_score (score DESC),
    INDEX idx_recommendations_expires_at (expires_at),
    INDEX idx_recommendations_user_score (user_id, score DESC),
    CONSTRAINT fk_recommendations_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_recommendations_articles FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 5.4 user_interests（ユーザー興味モデル）

```sql
CREATE TABLE user_interests (
    user_id CHAR(36) PRIMARY KEY,
    category_scores JSON DEFAULT '{}', -- {"tech": 0.8, "business": 0.6}
    tag_scores JSON DEFAULT '{}',
    keyword_scores JSON DEFAULT '{}',
    author_scores JSON DEFAULT '{}',
    domain_scores JSON DEFAULT '{}',
    reading_time_preference VARCHAR(20), -- 'short', 'medium', 'long'
    active_hours JSON DEFAULT '{}', -- {"9": 0.8, "20": 0.9} (hour: activity_score)
    model_version VARCHAR(50),
    last_calculated_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_user_interests_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 5.5 automation_rules（自動化ルール）

```sql
CREATE TABLE automation_rules (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    enabled BOOLEAN DEFAULT TRUE,
    trigger_type VARCHAR(50) NOT NULL, -- 'on_save', 'scheduled', 'condition_met'
    trigger_config JSON NOT NULL, -- トリガーの詳細設定
    conditions JSON NOT NULL, -- 条件式
    actions JSON NOT NULL, -- 実行するアクション
    last_triggered_at TIMESTAMP NULL,
    trigger_count INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_automation_rules_user_id (user_id),
    INDEX idx_automation_rules_enabled (enabled),
    CONSTRAINT fk_automation_rules_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 5.6 automation_logs（自動化実行ログ）

```sql
CREATE TABLE automation_logs (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    rule_id CHAR(36) NOT NULL,
    triggered_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    trigger_data JSON,
    actions_executed JSON,
    status VARCHAR(20) DEFAULT 'success', -- 'success', 'failed', 'partial'
    error_message TEXT,
    execution_time_ms INT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_automation_logs_rule_id (rule_id),
    INDEX idx_automation_logs_triggered_at (triggered_at DESC),
    CONSTRAINT fk_automation_logs_rules FOREIGN KEY (rule_id) REFERENCES automation_rules(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

## 6. 共通機能テーブル

### 6.1 job_queue（非同期ジョブ管理）

```sql
CREATE TABLE job_queue (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36),
    job_type VARCHAR(50) NOT NULL, -- 'summarize', 'generate_podcast', 'calculate_similarity'
    priority INT DEFAULT 5, -- 1-10, 10が最高
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'processing', 'completed', 'failed'
    payload JSON NOT NULL,
    result JSON,
    error_message TEXT,
    retry_count INT DEFAULT 0,
    max_retries INT DEFAULT 3,
    started_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_job_queue_status (status),
    INDEX idx_job_queue_priority (priority DESC),
    INDEX idx_job_queue_user_id (user_id),
    CONSTRAINT fk_job_queue_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### 6.2 api_logs（APIアクセスログ）

```sql
CREATE TABLE api_logs (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    user_id CHAR(36),
    method VARCHAR(10) NOT NULL,
    path TEXT NOT NULL,
    status_code INT NOT NULL,
    response_time_ms INT,
    ip_address VARCHAR(45),
    user_agent TEXT,
    request_id CHAR(36),
    error_code VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_api_logs_user_id (user_id),
    INDEX idx_api_logs_created_at (created_at),
    INDEX idx_api_logs_path (path(255)),
    CONSTRAINT fk_api_logs_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci
PARTITION BY RANGE (YEAR(created_at) * 100 + MONTH(created_at)) (
    PARTITION p202401 VALUES LESS THAN (202402),
    PARTITION p202402 VALUES LESS THAN (202403),
    PARTITION p202403 VALUES LESS THAN (202404),
    PARTITION p202404 VALUES LESS THAN (202405),
    PARTITION p202405 VALUES LESS THAN (202406),
    PARTITION p202406 VALUES LESS THAN (202407)
);
```

## 7. パフォーマンス最適化

### 7.1 複合インデックス

```sql
-- よく使用される検索パターン用の複合インデックス
ALTER TABLE articles ADD INDEX idx_user_category_status (user_id, category_id, status);
ALTER TABLE articles ADD INDEX idx_user_saved_at (user_id, saved_at DESC);
ALTER TABLE user_statistics ADD INDEX idx_user_year_month (user_id, YEAR(date), MONTH(date));
```

### 7.2 ビューの作成

```sql
-- カテゴリ別記事数ビュー
CREATE VIEW v_category_article_counts AS
SELECT 
    c.id,
    c.user_id,
    c.name,
    COUNT(a.id) as article_count
FROM categories c
LEFT JOIN articles a ON c.id = a.category_id AND a.deleted_at IS NULL
WHERE c.deleted_at IS NULL
GROUP BY c.id, c.user_id, c.name;

-- ユーザー別統計サマリービュー
CREATE VIEW v_user_summary AS
SELECT 
    u.id as user_id,
    COUNT(DISTINCT a.id) as total_articles,
    COUNT(DISTINCT CASE WHEN a.status = 'unread' THEN a.id END) as unread_count,
    COUNT(DISTINCT CASE WHEN a.status = 'read' THEN a.id END) as read_count,
    COUNT(DISTINCT t.id) as total_tags
FROM users u
LEFT JOIN articles a ON u.id = a.user_id AND a.deleted_at IS NULL
LEFT JOIN tags t ON u.id = t.user_id
WHERE u.deleted_at IS NULL
GROUP BY u.id;
```

## 8. ストアドプロシージャ・関数

### 8.1 記事保存時のカテゴリカウント更新

```sql
DELIMITER //

CREATE TRIGGER update_category_count_after_insert
AFTER INSERT ON articles
FOR EACH ROW
BEGIN
    IF NEW.category_id IS NOT NULL THEN
        UPDATE categories 
        SET article_count = article_count + 1 
        WHERE id = NEW.category_id;
    END IF;
END//

CREATE TRIGGER update_category_count_after_update
AFTER UPDATE ON articles
FOR EACH ROW
BEGIN
    IF OLD.category_id != NEW.category_id THEN
        IF OLD.category_id IS NOT NULL THEN
            UPDATE categories 
            SET article_count = article_count - 1 
            WHERE id = OLD.category_id;
        END IF;
        IF NEW.category_id IS NOT NULL THEN
            UPDATE categories 
            SET article_count = article_count + 1 
            WHERE id = NEW.category_id;
        END IF;
    END IF;
END//

CREATE TRIGGER update_category_count_after_delete
AFTER DELETE ON articles
FOR EACH ROW
BEGIN
    IF OLD.category_id IS NOT NULL THEN
        UPDATE categories 
        SET article_count = article_count - 1 
        WHERE id = OLD.category_id;
    END IF;
END//

DELIMITER ;
```

### 8.2 タグ使用数の更新

```sql
DELIMITER //

CREATE TRIGGER update_tag_usage_after_insert
AFTER INSERT ON article_tags
FOR EACH ROW
BEGIN
    UPDATE tags 
    SET usage_count = usage_count + 1 
    WHERE id = NEW.tag_id;
END//

CREATE TRIGGER update_tag_usage_after_delete
AFTER DELETE ON article_tags
FOR EACH ROW
BEGIN
    UPDATE tags 
    SET usage_count = usage_count - 1 
    WHERE id = OLD.tag_id;
END//

DELIMITER ;
```

## 9. セキュリティ設定

### 9.1 ユーザー権限設定

```sql
-- アプリケーション用ユーザーの作成
CREATE USER 'readlater_app'@'%' IDENTIFIED BY 'secure_password';

-- 必要最小限の権限付与
GRANT SELECT, INSERT, UPDATE, DELETE ON readlater_db.* TO 'readlater_app'@'%';

-- ストアドプロシージャの実行権限
GRANT EXECUTE ON readlater_db.* TO 'readlater_app'@'%';

-- 読み取り専用ユーザー（分析用）
CREATE USER 'readlater_readonly'@'%' IDENTIFIED BY 'readonly_password';
GRANT SELECT ON readlater_db.* TO 'readlater_readonly'@'%';
```

### 9.2 暗号化設定

```sql
-- 機密データの暗号化（MySQL 8.0の透過的データ暗号化）
ALTER TABLE users ENCRYPTION='Y';
ALTER TABLE user_sessions ENCRYPTION='Y';
ALTER TABLE api_logs ENCRYPTION='Y';
```

## 10. バックアップ・メンテナンス

### 10.1 バックアップスクリプト

```bash
#!/bin/bash
# 日次バックアップスクリプト
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/mysql"

# フルバックアップ
mysqldump -u root -p \
  --single-transaction \
  --routines \
  --triggers \
  --events \
  readlater_db > $BACKUP_DIR/readlater_db_$DATE.sql

# 圧縮
gzip $BACKUP_DIR/readlater_db_$DATE.sql

# 7日以上前のバックアップを削除
find $BACKUP_DIR -name "*.sql.gz" -mtime +7 -delete
```

### 10.2 定期メンテナンス

```sql
-- テーブル最適化（月次実行）
OPTIMIZE TABLE articles;
OPTIMIZE TABLE user_statistics;
OPTIMIZE TABLE api_logs;

-- インデックス統計の更新
ANALYZE TABLE articles;
ANALYZE TABLE categories;
ANALYZE TABLE tags;
```

## 11. 監視クエリ

### 11.1 テーブルサイズ監視

```sql
SELECT 
    table_name,
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)',
    table_rows AS 'Row Count'
FROM information_schema.tables
WHERE table_schema = 'readlater_db'
ORDER BY (data_length + index_length) DESC;
```

### 11.2 スロークエリ監視

```sql
-- スロークエリログの有効化
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 2;
SET GLOBAL slow_query_log_file = '/var/log/mysql/slow.log';

-- パフォーマンススキーマを使用した分析
SELECT 
    digest_text,
    count_star,
    avg_timer_wait/1000000000 AS avg_time_ms,
    sum_timer_wait/1000000000 AS total_time_ms
FROM performance_schema.events_statements_summary_by_digest
ORDER BY sum_timer_wait DESC
LIMIT 10;
```