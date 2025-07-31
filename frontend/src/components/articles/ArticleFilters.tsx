'use client';

import React from 'react';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { DropdownMenu, DropdownMenuCheckboxItem, DropdownMenuContent, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { Archive, BookmarkCheck, Eye, EyeOff, Filter, Heart, RotateCcw, X } from 'lucide-react';
import type { Category, ArticleFilters as IArticleFilters } from '@/types/article';

interface ArticleFiltersProps {
  filters: IArticleFilters;
  categories: Category[];
  onFiltersChange: (filters: IArticleFilters) => void;
  className?: string;
}

export const ArticleFilters: React.FC<ArticleFiltersProps> = ({
  filters,
  categories,
  onFiltersChange,
  className,
}) => {
  const handleStatusChange = (status: IArticleFilters['status']) => {
    onFiltersChange({
      ...filters,
      status: filters.status === status ? undefined : status,
    });
  };

  const handleCategoryChange = (categoryId: string) => {
    onFiltersChange({
      ...filters,
      categoryId: categoryId === 'all' ? undefined : categoryId,
    });
  };

  const handleFavoriteChange = () => {
    onFiltersChange({
      ...filters,
      favorite: filters.favorite ? undefined : true,
    });
  };

  const handleClearFilters = () => {
    onFiltersChange({
      page: 1,
      limit: filters.limit,
    });
  };

  const getActiveFiltersCount = () => {
    let count = 0;
    if (filters.status) count++;
    if (filters.categoryId) count++;
    if (filters.favorite) count++;
    return count;
  };

  const selectedCategory = categories.find(cat => cat.id === filters.categoryId);

  return (
    <div className={`flex flex-wrap items-center gap-2 ${className}`}>
      {/* ステータスフィルター */}
      <div className="flex gap-1">
        <Button
          variant={filters.status === 'unread' ? 'default' : 'outline'}
          size="sm"
          onClick={() => handleStatusChange('unread')}
          className="h-8"
        >
          <Eye className="mr-1 h-3 w-3" />
          未読
        </Button>
        <Button
          variant={filters.status === 'read' ? 'default' : 'outline'}
          size="sm"
          onClick={() => handleStatusChange('read')}
          className="h-8"
        >
          <EyeOff className="mr-1 h-3 w-3" />
          既読
        </Button>
        <Button
          variant={filters.status === 'archived' ? 'default' : 'outline'}
          size="sm"
          onClick={() => handleStatusChange('archived')}
          className="h-8"
        >
          <Archive className="mr-1 h-3 w-3" />
          アーカイブ
        </Button>
      </div>

      {/* カテゴリフィルター */}
      <Select
        value={filters.categoryId || 'all'}
        onValueChange={handleCategoryChange}
      >
        <SelectTrigger className="w-48 h-8">
          <SelectValue placeholder="カテゴリで絞り込み" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">すべてのカテゴリ</SelectItem>
          {categories.map(category => (
            <SelectItem key={category.id} value={category.id}>
              <div className="flex items-center gap-2">
                <div 
                  className="w-3 h-3 rounded-full" 
                  style={{ backgroundColor: category.color }}
                />
                <span>{category.name}</span>
                <Badge variant="secondary" className="ml-auto text-xs">
                  {category.articleCount}
                </Badge>
              </div>
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      {/* お気に入りフィルター */}
      <Button
        variant={filters.favorite ? 'default' : 'outline'}
        size="sm"
        onClick={handleFavoriteChange}
        className="h-8"
      >
        <Heart className={`mr-1 h-3 w-3 ${filters.favorite ? 'fill-current' : ''}`} />
        お気に入り
      </Button>

      {/* その他のフィルター */}
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="outline" size="sm" className="h-8">
            <Filter className="mr-1 h-3 w-3" />
            フィルター
            {getActiveFiltersCount() > 0 && (
              <Badge variant="secondary" className="ml-1 text-xs min-w-[1.25rem] h-5">
                {getActiveFiltersCount()}
              </Badge>
            )}
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-56">
          <DropdownMenuLabel>表示件数</DropdownMenuLabel>
          <DropdownMenuCheckboxItem
            checked={filters.limit === 20}
            onCheckedChange={() => onFiltersChange({ ...filters, limit: 20, page: 1 })}
          >
            20件
          </DropdownMenuCheckboxItem>
          <DropdownMenuCheckboxItem
            checked={filters.limit === 50}
            onCheckedChange={() => onFiltersChange({ ...filters, limit: 50, page: 1 })}
          >
            50件
          </DropdownMenuCheckboxItem>
          <DropdownMenuCheckboxItem
            checked={filters.limit === 100}
            onCheckedChange={() => onFiltersChange({ ...filters, limit: 100, page: 1 })}
          >
            100件
          </DropdownMenuCheckboxItem>
        </DropdownMenuContent>
      </DropdownMenu>

      {/* アクティブフィルターの表示 */}
      {getActiveFiltersCount() > 0 && (
        <>
          <div className="flex flex-wrap gap-1">
            {filters.status && (
              <div className="flex items-center gap-1 px-2 py-1 bg-primary/10 text-primary rounded-md text-xs">
                <span>ステータス: {
                  filters.status === 'unread' ? '未読' :
                  filters.status === 'read' ? '既読' : 'アーカイブ'
                }</span>
                <button
                  onClick={() => handleStatusChange(filters.status!)}
                  className="hover:opacity-80"
                >
                  <X className="h-3 w-3" />
                </button>
              </div>
            )}
            {selectedCategory && (
              <div className="flex items-center gap-1 px-2 py-1 bg-primary/10 text-primary rounded-md text-xs">
                <div 
                  className="w-2 h-2 rounded-full" 
                  style={{ backgroundColor: selectedCategory.color }}
                />
                <span>{selectedCategory.name}</span>
                <button
                  onClick={() => handleCategoryChange('all')}
                  className="hover:opacity-80"
                >
                  <X className="h-3 w-3" />
                </button>
              </div>
            )}
            {filters.favorite && (
              <div className="flex items-center gap-1 px-2 py-1 bg-primary/10 text-primary rounded-md text-xs">
                <Heart className="h-3 w-3 fill-current" />
                <span>お気に入り</span>
                <button
                  onClick={handleFavoriteChange}
                  className="hover:opacity-80"
                >
                  <X className="h-3 w-3" />
                </button>
              </div>
            )}
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={handleClearFilters}
            className="h-8 text-muted-foreground"
          >
            <RotateCcw className="mr-1 h-3 w-3" />
            リセット
          </Button>
        </>
      )}
    </div>
  );
};