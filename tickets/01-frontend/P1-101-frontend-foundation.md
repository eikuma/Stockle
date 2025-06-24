# P1-101: フロントエンド基盤構築

## 概要
Next.js 14 + TypeScript + Tailwind CSSによるフロントエンド基盤の構築

## 担当者
**Member 1 (Frontend Developer)**

## 優先度
**最高** - 他のフロントエンド作業の基盤

## 前提条件
- P1-001: プロジェクト基盤セットアップが完了済み
- Node.js 18以上がインストール済み

## 作業内容

### 1. Next.js プロジェクトの初期化
- [ ] `frontend/` ディレクトリに移動
- [ ] Next.js 14をApp Routerで初期化
- [ ] TypeScript設定を有効化
- [ ] `next.config.js` の設定
- [ ] フォルダ構造の作成

### 2. 必要パッケージのインストール
- [ ] UI フレームワーク: shadcn/ui, Radix UI
- [ ] スタイリング: Tailwind CSS, class-variance-authority
- [ ] 状態管理: Zustand
- [ ] フォーム: React Hook Form, Zod
- [ ] データフェッチング: TanStack Query
- [ ] 認証: NextAuth.js
- [ ] アニメーション: Framer Motion
- [ ] アイコン: Lucide React
- [ ] 開発ツール: ESLint, Prettier, Husky

### 3. TypeScript設定
- [ ] `tsconfig.json` の詳細設定
- [ ] 厳密型チェックの有効化
- [ ] パスエイリアスの設定
- [ ] 型定義ファイルの作成

### 4. Tailwind CSS設定
- [ ] `tailwind.config.js` の設定
- [ ] カスタムカラーパレットの追加
- [ ] レスポンシブブレークポイントの設定
- [ ] カスタムコンポーネントクラスの定義

### 5. ESLint + Prettier設定
- [ ] `.eslintrc.json` の設定
- [ ] `.prettierrc` の設定
- [ ] `package.json` スクリプトの追加
- [ ] Husky + lint-staged の設定

### 6. フォルダ構造の作成
- [ ] `src/app/` - App Router ページ
- [ ] `src/components/` - 再利用可能コンポーネント
- [ ] `src/features/` - 機能別モジュール
- [ ] `src/hooks/` - カスタムフック
- [ ] `src/lib/` - ユーティリティ関数
- [ ] `src/services/` - API クライアント
- [ ] `src/stores/` - Zustand ストア
- [ ] `src/types/` - TypeScript 型定義

### 7. shadcn/ui コンポーネントのセットアップ
- [ ] shadcn/ui の初期化
- [ ] 基本コンポーネントのインストール
  - [ ] Button
  - [ ] Input
  - [ ] Card
  - [ ] Dialog
  - [ ] Dropdown Menu
  - [ ] Form
  - [ ] Label
  - [ ] Select
  - [ ] Table
  - [ ] Textarea
  - [ ] Toast

## 実装詳細

### package.json 依存関係
```json
{
  "dependencies": {
    "next": "^14.0.0",
    "react": "^18.0.0",
    "react-dom": "^18.0.0",
    "@radix-ui/react-dropdown-menu": "^2.0.6",
    "@radix-ui/react-dialog": "^1.0.5",
    "@radix-ui/react-select": "^2.0.0",
    "class-variance-authority": "^0.7.0",
    "clsx": "^2.0.0",
    "tailwind-merge": "^2.0.0",
    "zustand": "^4.4.0",
    "react-hook-form": "^7.47.0",
    "zod": "^3.22.0",
    "@hookform/resolvers": "^3.3.0",
    "@tanstack/react-query": "^5.0.0",
    "next-auth": "^4.24.0",
    "framer-motion": "^10.16.0",
    "lucide-react": "^0.294.0"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "@types/react": "^18.0.0",
    "@types/react-dom": "^18.0.0",
    "typescript": "^5.0.0",
    "tailwindcss": "^3.3.0",
    "eslint": "^8.0.0",
    "eslint-config-next": "^14.0.0",
    "prettier": "^3.0.0",
    "husky": "^8.0.0",
    "lint-staged": "^15.0.0"
  }
}
```

### tailwind.config.js
```javascript
/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: ["class"],
  content: [
    './pages/**/*.{ts,tsx}',
    './components/**/*.{ts,tsx}',
    './app/**/*.{ts,tsx}',
    './src/**/*.{ts,tsx}',
  ],
  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        "2xl": "1400px",
      },
    },
    extend: {
      colors: {
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
      keyframes: {
        "accordion-down": {
          from: { height: 0 },
          to: { height: "var(--radix-accordion-content-height)" },
        },
        "accordion-up": {
          from: { height: "var(--radix-accordion-content-height)" },
          to: { height: 0 },
        },
      },
      animation: {
        "accordion-down": "accordion-down 0.2s ease-out",
        "accordion-up": "accordion-up 0.2s ease-out",
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
}
```

### src/types/index.ts
```typescript
// API Response Types
export interface User {
  id: string;
  email: string;
  displayName: string;
  profileImageUrl?: string;
  authProvider: 'email' | 'google';
  createdAt: string;
  updatedAt: string;
}

export interface Category {
  id: string;
  userId: string;
  name: string;
  color: string;
  displayOrder: number;
  isDefault: boolean;
  articleCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface Article {
  id: string;
  userId: string;
  categoryId?: string;
  url: string;
  title: string;
  summary?: string;
  thumbnailUrl?: string;
  author?: string;
  siteName?: string;
  publishedAt?: string;
  savedAt: string;
  lastAccessedAt?: string;
  status: 'unread' | 'read' | 'archived';
  isFavorite: boolean;
  readingProgress: number;
  wordCount?: number;
  category?: Category;
}

// Form Types
export interface SaveArticleForm {
  url: string;
  categoryId?: string;
  tags?: string[];
}

export interface CategoryForm {
  name: string;
  color: string;
}

// Store Types
export interface AuthStore {
  user: User | null;
  isAuthenticated: boolean;
  login: (user: User) => void;
  logout: () => void;
}

export interface ArticleStore {
  articles: Article[];
  loading: boolean;
  error: string | null;
  filters: {
    status?: 'unread' | 'read' | 'archived';
    categoryId?: string;
    search?: string;
  };
  setArticles: (articles: Article[]) => void;
  addArticle: (article: Article) => void;
  updateArticle: (id: string, article: Partial<Article>) => void;
  deleteArticle: (id: string) => void;
  setFilters: (filters: Partial<ArticleStore['filters']>) => void;
}
```

## 受入条件

### 必須条件
- [ ] `npm run dev` でNext.js開発サーバーが起動する
- [ ] `npm run build` でエラーなくビルドが完了する
- [ ] ESLint + Prettierが正常に動作する
- [ ] TypeScript型チェックがエラーなく通る
- [ ] shadcn/ui コンポーネントが正常に表示される

### 品質条件
- [ ] ページの初期表示が2秒以内
- [ ] レスポンシブデザインが正常に動作する
- [ ] アクセシビリティスコアが90以上
- [ ] TypeScript strict モードでエラーが0件

## 推定時間
**16時間** (3-4日)

## 依存関係
- P1-001: プロジェクト基盤セットアップ

## 完了後の次ステップ
1. P1-102: 認証システムの実装
2. P1-103: レイアウト・ナビゲーションの実装

## 備考
- shadcn/ui のカスタマイズは最小限に留める
- コンポーネントの再利用性を重視する
- パフォーマンスを意識した実装を行う
- アクセシビリティを考慮したコンポーネント設計