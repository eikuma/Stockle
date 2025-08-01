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

// Form Types
export interface CategoryForm {
  name: string;
  color: string;
}

// Re-export article types
export * from './article';
import type { Article } from './article';

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