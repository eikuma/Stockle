'use client';

import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Dialog, DialogTrigger } from '@/components/ui/dialog';
import { Plus } from 'lucide-react';
import { 
  ArticleCard, 
  SaveArticleDialog, 
  SearchBar, 
  ArticleFilters 
} from '@/components/articles';
import { mockArticles, mockCategories, mockTags } from '@/lib/mock/articleData';
import type { Article, ArticleFilters as IArticleFilters, SaveArticleForm } from '@/types/article';

export default function ArticlesPage() {
  const [articles, setArticles] = useState<Article[]>(mockArticles);
  const [filters, setFilters] = useState<IArticleFilters>({
    page: 1,
    limit: 20,
  });
  const [searchQuery, setSearchQuery] = useState('');
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  // フィルタリング処理
  const filteredArticles = articles.filter(article => {
    // 検索クエリフィルター
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      const matchesSearch = 
        article.title.toLowerCase().includes(query) ||
        article.summary?.toLowerCase().includes(query) ||
        article.author?.toLowerCase().includes(query) ||
        article.tags.some(tag => tag.name.toLowerCase().includes(query));
      
      if (!matchesSearch) return false;
    }

    // ステータスフィルター
    if (filters.status && article.status !== filters.status) {
      return false;
    }

    // カテゴリフィルター
    if (filters.categoryId && article.categoryId !== filters.categoryId) {
      return false;
    }

    // お気に入りフィルター
    if (filters.favorite && !article.isFavorite) {
      return false;
    }

    return true;
  });

  const handleSaveArticle = async (data: SaveArticleForm) => {
    // 実際のAPIコールの代わりにモックデータを作成
    const newArticle: Article = {
      id: `article-${Date.now()}`,
      userId: 'user-1',
      categoryId: data.categoryId,
      url: data.url,
      title: `保存された記事 - ${new URL(data.url).hostname}`,
      summary: 'この記事は保存されました。実際のAPIではWebスクレイピングによって記事の内容が取得されます。',
      thumbnailUrl: 'https://images.unsplash.com/photo-1432888498266-38ffec3eaf0a?w=400',
      author: '未知の著者',
      siteName: new URL(data.url).hostname,
      savedAt: new Date().toISOString(),
      status: 'unread',
      isFavorite: false,
      readingProgress: 0,
      readingTimeSeconds: 300,
      wordCount: 1500,
      language: 'ja',
      category: data.categoryId ? mockCategories.find(cat => cat.id === data.categoryId) : undefined,
      tags: data.tags ? data.tags.map(tagName => ({
        id: `tag-${Date.now()}-${tagName}`,
        userId: 'user-1',
        name: tagName,
        usageCount: 1,
      })) : [],
    };

    setArticles(prev => [newArticle, ...prev]);
  };

  const handleToggleFavorite = (articleId: string) => {
    setArticles(prev => prev.map(article => 
      article.id === articleId 
        ? { ...article, isFavorite: !article.isFavorite }
        : article
    ));
  };

  const handleUpdateStatus = (articleId: string, status: Article['status']) => {
    setArticles(prev => prev.map(article => 
      article.id === articleId 
        ? { ...article, status }
        : article
    ));
  };

  const handleDeleteArticle = (articleId: string) => {
    setArticles(prev => prev.filter(article => article.id !== articleId));
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex flex-col gap-6">
        {/* ヘッダー */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div>
            <h1 className="text-3xl font-bold">保存した記事</h1>
            <p className="text-muted-foreground mt-1">
              {filteredArticles.length}件の記事が見つかりました
            </p>
          </div>
          <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="mr-2 h-4 w-4" />
                記事を保存
              </Button>
            </DialogTrigger>
            <SaveArticleDialog
              open={isDialogOpen}
              onOpenChange={setIsDialogOpen}
              onSave={handleSaveArticle}
              categories={mockCategories}
              existingTags={mockTags}
            />
          </Dialog>
        </div>

        {/* 検索バー */}
        <SearchBar
          value={searchQuery}
          onChange={setSearchQuery}
          placeholder="タイトル、著者、タグで検索..."
        />

        {/* フィルター */}
        <ArticleFilters
          filters={filters}
          categories={mockCategories}
          onFiltersChange={setFilters}
        />

        {/* 記事一覧 */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredArticles.map(article => (
            <ArticleCard
              key={article.id}
              article={article}
              onToggleFavorite={handleToggleFavorite}
              onUpdateStatus={handleUpdateStatus}
              onDelete={handleDeleteArticle}
            />
          ))}
        </div>

        {/* 記事が見つからない場合 */}
        {filteredArticles.length === 0 && (
          <div className="text-center py-12">
            <h2 className="text-xl font-semibold mb-2">記事が見つかりません</h2>
            <p className="text-muted-foreground mb-4">
              検索条件を変更するか、新しい記事を保存してください。
            </p>
            <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
              <DialogTrigger asChild>
                <Button variant="outline">
                  <Plus className="mr-2 h-4 w-4" />
                  記事を保存
                </Button>
              </DialogTrigger>
            </Dialog>
          </div>
        )}
      </div>
    </div>
  );
}