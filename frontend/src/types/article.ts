export interface Article {
  id: string;
  userId: string;
  categoryId?: string;
  url: string;
  title: string;
  content?: string;
  summary?: string;
  summaryShort?: string;
  summaryLong?: string;
  thumbnailUrl?: string;
  author?: string;
  siteName?: string;
  publishedAt?: string;
  savedAt: string;
  lastAccessedAt?: string;
  status: 'unread' | 'read' | 'archived';
  isFavorite: boolean;
  readingProgress: number;
  readingTimeSeconds: number;
  wordCount?: number;
  language: string;
  category?: Category;
  tags: Tag[];
}

export interface Category {
  id: string;
  userId: string;
  name: string;
  color: string;
  displayOrder: number;
  isDefault: boolean;
  articleCount: number;
  articles?: Article[];
}

export interface Tag {
  id: string;
  userId: string;
  name: string;
  usageCount: number;
}

export interface SaveArticleForm {
  url: string;
  categoryId?: string;
  tags?: string[];
}

export interface ArticleFilters {
  status?: 'unread' | 'read' | 'archived';
  categoryId?: string;
  search?: string;
  page?: number;
  limit?: number;
  favorite?: boolean;
}

export interface ArticlesResponse {
  articles: Article[];
  total: number;
  page: number;
  limit: number;
}

export interface ArticleResponse {
  article: Article;
}

export interface SaveArticleResponse {
  message: string;
  article: Article;
}

export interface UpdateArticleResponse {
  message: string;
}

export interface DeleteArticleResponse {
  message: string;
}