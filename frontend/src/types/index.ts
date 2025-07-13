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