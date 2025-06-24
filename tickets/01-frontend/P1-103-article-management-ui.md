# P1-103: 記事管理UI実装

## 概要
記事保存・一覧表示・詳細表示・検索機能のUI実装

## 担当者
**Member 1 (Frontend Developer)**

## 優先度
**最高** - コア機能の実装

## 前提条件
- P1-101: フロントエンド基盤構築が完了済み
- P1-102: 認証システムUIが完了済み
- Member 2の記事管理APIと並行作業可能

## 作業内容

### 1. 記事保存フォーム
- [ ] URL入力フォームの実装
- [ ] カテゴリ選択機能
- [ ] タグ入力機能
- [ ] バリデーション機能
- [ ] 保存状況の表示

### 2. 記事一覧画面
- [ ] 記事カードコンポーネント
- [ ] グリッド/リスト表示切り替え
- [ ] 無限スクロール実装
- [ ] 既読/未読フィルター
- [ ] カテゴリフィルター
- [ ] 並び順変更機能

### 3. 記事詳細画面
- [ ] 記事詳細表示
- [ ] 要約表示機能
- [ ] 元記事リンク
- [ ] カテゴリ・タグ編集
- [ ] ステータス変更機能
- [ ] 削除機能

### 4. 検索機能
- [ ] 検索バーコンポーネント
- [ ] 検索結果表示
- [ ] 高度な検索オプション
- [ ] 検索履歴機能
- [ ] 検索候補表示

### 5. カテゴリ管理
- [ ] カテゴリ一覧表示
- [ ] カテゴリ作成・編集・削除
- [ ] カテゴリ色設定
- [ ] ドラッグ&ドロップ並び替え

### 6. レスポンシブデザイン
- [ ] モバイル対応
- [ ] タブレット対応
- [ ] デスクトップ最適化
- [ ] タッチ操作対応

## 実装詳細

### components/articles/ArticleCard.tsx
```typescript
'use client';

import { useState } from 'react';
import Image from 'next/image';
import Link from 'next/link';
import { Card, CardContent, CardFooter } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { Eye, EyeOff, Star, MoreVertical, ExternalLink } from 'lucide-react';
import { Article } from '@/types';
import { cn } from '@/lib/utils';

interface ArticleCardProps {
  article: Article;
  onStatusChange: (id: string, status: 'read' | 'unread') => void;
  onFavoriteToggle: (id: string) => void;
  onDelete: (id: string) => void;
  view?: 'grid' | 'list';
}

export function ArticleCard({ 
  article, 
  onStatusChange, 
  onFavoriteToggle, 
  onDelete,
  view = 'grid' 
}: ArticleCardProps) {
  const [imageError, setImageError] = useState(false);

  const handleStatusToggle = () => {
    onStatusChange(article.id, article.status === 'read' ? 'unread' : 'read');
  };

  const handleFavoriteToggle = () => {
    onFavoriteToggle(article.id);
  };

  return (
    <Card className={cn(
      'group transition-all duration-200 hover:shadow-md',
      view === 'list' && 'flex flex-row',
      article.status === 'unread' && 'border-blue-200'
    )}>
      {article.thumbnailUrl && !imageError && (
        <div className={cn(
          'relative overflow-hidden rounded-t-lg',
          view === 'list' ? 'w-32 h-24 flex-shrink-0' : 'h-48'
        )}>
          <Image
            src={article.thumbnailUrl}
            alt={article.title}
            fill
            className="object-cover group-hover:scale-105 transition-transform duration-200"
            onError={() => setImageError(true)}
          />
          {article.status === 'unread' && (
            <div className="absolute top-2 left-2">
              <Badge variant="secondary" className="bg-blue-100 text-blue-800">
                未読
              </Badge>
            </div>
          )}
        </div>
      )}

      <CardContent className={cn(
        'p-4',
        view === 'list' && 'flex-1'
      )}>
        <div className="space-y-2">
          <div className="flex items-start justify-between">
            <Link href={`/articles/${article.id}`}>
              <h3 className="font-semibold line-clamp-2 hover:text-blue-600 transition-colors">
                {article.title}
              </h3>
            </Link>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="sm" className="opacity-0 group-hover:opacity-100">
                  <MoreVertical className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={handleStatusToggle}>
                  {article.status === 'read' ? (
                    <>
                      <EyeOff className="h-4 w-4 mr-2" />
                      未読にする
                    </>
                  ) : (
                    <>
                      <Eye className="h-4 w-4 mr-2" />
                      既読にする
                    </>
                  )}
                </DropdownMenuItem>
                <DropdownMenuItem onClick={handleFavoriteToggle}>
                  <Star className={cn(
                    "h-4 w-4 mr-2",
                    article.isFavorite && "fill-yellow-400 text-yellow-400"
                  )} />
                  {article.isFavorite ? 'お気に入り解除' : 'お気に入り追加'}
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <a href={article.url} target="_blank" rel="noopener noreferrer">
                    <ExternalLink className="h-4 w-4 mr-2" />
                    元記事を開く
                  </a>
                </DropdownMenuItem>
                <DropdownMenuItem 
                  onClick={() => onDelete(article.id)}
                  className="text-red-600"
                >
                  削除
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>

          {article.summary && (
            <p className="text-sm text-muted-foreground line-clamp-3">
              {article.summary}
            </p>
          )}

          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              {article.category && (
                <Badge variant="outline" style={{ backgroundColor: article.category.color + '20' }}>
                  {article.category.name}
                </Badge>
              )}
              {article.author && (
                <span className="text-xs text-muted-foreground">
                  by {article.author}
                </span>
              )}
            </div>
            <time className="text-xs text-muted-foreground">
              {new Date(article.savedAt).toLocaleDateString('ja-JP')}
            </time>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
```

### components/articles/SaveArticleDialog.tsx
```typescript
'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import { Plus, Loader2 } from 'lucide-react';
import { useCategories } from '@/hooks/useCategories';
import { useSaveArticle } from '@/hooks/useArticles';

const saveArticleSchema = z.object({
  url: z.string().url('有効なURLを入力してください'),
  categoryId: z.string().optional(),
  tags: z.string().optional(),
});

type SaveArticleForm = z.infer<typeof saveArticleSchema>;

export function SaveArticleDialog() {
  const [open, setOpen] = useState(false);
  const { data: categories } = useCategories();
  const saveArticleMutation = useSaveArticle();

  const form = useForm<SaveArticleForm>({
    resolver: zodResolver(saveArticleSchema),
    defaultValues: {
      url: '',
      categoryId: '',
      tags: '',
    },
  });

  const onSubmit = async (data: SaveArticleForm) => {
    try {
      const tags = data.tags
        ? data.tags.split(',').map(tag => tag.trim()).filter(Boolean)
        : [];

      await saveArticleMutation.mutateAsync({
        url: data.url,
        categoryId: data.categoryId || undefined,
        tags,
      });

      form.reset();
      setOpen(false);
    } catch (error) {
      console.error('Failed to save article:', error);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus className="h-4 w-4 mr-2" />
          記事を保存
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>記事を保存</DialogTitle>
        </DialogHeader>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="url">URL *</Label>
            <Input
              id="url"
              placeholder="https://example.com/article"
              {...form.register('url')}
            />
            {form.formState.errors.url && (
              <p className="text-sm text-red-600">{form.formState.errors.url.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="category">カテゴリ</Label>
            <Select
              value={form.watch('categoryId')}
              onValueChange={(value) => form.setValue('categoryId', value)}
            >
              <SelectTrigger>
                <SelectValue placeholder="カテゴリを選択" />
              </SelectTrigger>
              <SelectContent>
                {categories?.map((category) => (
                  <SelectItem key={category.id} value={category.id}>
                    <div className="flex items-center space-x-2">
                      <div
                        className="w-3 h-3 rounded-full"
                        style={{ backgroundColor: category.color }}
                      />
                      <span>{category.name}</span>
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="tags">タグ (カンマ区切り)</Label>
            <Textarea
              id="tags"
              placeholder="例: AI, 機械学習, 技術"
              rows={2}
              {...form.register('tags')}
            />
          </div>

          <div className="flex justify-end space-x-2">
            <Button type="button" variant="outline" onClick={() => setOpen(false)}>
              キャンセル
            </Button>
            <Button type="submit" disabled={saveArticleMutation.isPending}>
              {saveArticleMutation.isPending && (
                <Loader2 className="h-4 w-4 mr-2 animate-spin" />
              )}
              保存
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
```

### hooks/useArticles.ts
```typescript
'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { articlesService } from '@/services/articles';
import { Article, SaveArticleForm } from '@/types';

export function useArticles(filters?: {
  status?: 'unread' | 'read' | 'archived';
  categoryId?: string;
  search?: string;
  page?: number;
  limit?: number;
}) {
  return useQuery({
    queryKey: ['articles', filters],
    queryFn: () => articlesService.getArticles(filters),
  });
}

export function useArticle(id: string) {
  return useQuery({
    queryKey: ['articles', id],
    queryFn: () => articlesService.getArticle(id),
    enabled: !!id,
  });
}

export function useSaveArticle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: SaveArticleForm) => articlesService.saveArticle(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['articles'] });
    },
  });
}

export function useUpdateArticle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Partial<Article> }) =>
      articlesService.updateArticle(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['articles'] });
      queryClient.invalidateQueries({ queryKey: ['articles', variables.id] });
    },
  });
}

export function useDeleteArticle() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => articlesService.deleteArticle(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['articles'] });
    },
  });
}
```

### app/dashboard/page.tsx
```typescript
'use client';

import { useState } from 'react';
import { useArticles } from '@/hooks/useArticles';
import { useCategories } from '@/hooks/useCategories';
import { ArticleCard } from '@/components/articles/ArticleCard';
import { SaveArticleDialog } from '@/components/articles/SaveArticleDialog';
import { SearchBar } from '@/components/articles/SearchBar';
import { ArticleFilters } from '@/components/articles/ArticleFilters';
import { ViewToggle } from '@/components/articles/ViewToggle';
import { LoadingSpinner } from '@/components/ui/loading-spinner';
import { AuthGuard } from '@/components/auth/AuthGuard';

export default function DashboardPage() {
  const [filters, setFilters] = useState({
    status: undefined as 'unread' | 'read' | 'archived' | undefined,
    categoryId: undefined as string | undefined,
    search: undefined as string | undefined,
  });
  const [view, setView] = useState<'grid' | 'list'>('grid');

  const { data: articles, isLoading, error } = useArticles(filters);
  const { data: categories } = useCategories();

  if (isLoading) {
    return <LoadingSpinner />;
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-red-600">記事の読み込みに失敗しました</p>
      </div>
    );
  }

  return (
    <AuthGuard>
      <div className="container mx-auto px-4 py-8">
        <div className="flex items-center justify-between mb-8">
          <h1 className="text-3xl font-bold">保存した記事</h1>
          <SaveArticleDialog />
        </div>

        <div className="flex flex-col md:flex-row gap-4 mb-6">
          <div className="flex-1">
            <SearchBar
              value={filters.search || ''}
              onChange={(search) => setFilters(prev => ({ ...prev, search }))}
            />
          </div>
          <ArticleFilters
            categories={categories || []}
            filters={filters}
            onChange={setFilters}
          />
          <ViewToggle view={view} onChange={setView} />
        </div>

        {articles?.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-muted-foreground mb-4">まだ記事が保存されていません</p>
            <SaveArticleDialog />
          </div>
        ) : (
          <div className={cn(
            'gap-4',
            view === 'grid' 
              ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3' 
              : 'flex flex-col'
          )}>
            {articles?.map((article) => (
              <ArticleCard
                key={article.id}
                article={article}
                view={view}
                onStatusChange={(id, status) => {
                  // Update article status
                }}
                onFavoriteToggle={(id) => {
                  // Toggle favorite status
                }}
                onDelete={(id) => {
                  // Delete article
                }}
              />
            ))}
          </div>
        )}
      </div>
    </AuthGuard>
  );
}
```

## 受入条件

### 必須条件
- [ ] 記事保存フォームが正常に動作する
- [ ] 記事一覧が正常に表示される
- [ ] 検索機能が正常に動作する
- [ ] フィルタリング機能が正常に動作する
- [ ] レスポンシブデザインが正常に動作する
- [ ] 既読/未読の切り替えが正常に動作する

### 品質条件
- [ ] ページの読み込みが3秒以内
- [ ] 無限スクロールが滑らか
- [ ] アクセシビリティスコアが90以上
- [ ] モバイルでのタッチ操作が快適
- [ ] エラーハンドリングが適切

## 推定時間
**32時間** (6-8日)

## 依存関係
- P1-101: フロントエンド基盤構築
- P1-102: 認証システムUI
- Member 2の記事管理API

## 完了後の次ステップ
1. Member 2のAPIとの統合テスト
2. パフォーマンス最適化
3. ユーザビリティテスト

## 備考
- UX/UIを最優先に考慮
- パフォーマンスを意識した実装
- アクセシビリティを重視
- モバイルファーストでの実装