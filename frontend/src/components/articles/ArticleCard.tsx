'use client';

import React from 'react';
import Link from 'next/link';
import { formatDistanceToNow } from 'date-fns';
import { ja } from 'date-fns/locale';
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { Bookmark, BookmarkX, Clock, Eye, EyeOff, Globe, MoreHorizontal, Tag, Trash2 } from 'lucide-react';
import type { Article } from '@/types/article';

interface ArticleCardProps {
  article: Article;
  onToggleFavorite?: (articleId: string) => void;
  onUpdateStatus?: (articleId: string, status: Article['status']) => void;
  onDelete?: (articleId: string) => void;
}

export const ArticleCard: React.FC<ArticleCardProps> = ({
  article,
  onToggleFavorite,
  onUpdateStatus,
  onDelete,
}) => {
  const handleStatusChange = (status: Article['status']) => {
    if (onUpdateStatus) {
      onUpdateStatus(article.id, status);
    }
  };

  const getStatusIcon = () => {
    switch (article.status) {
      case 'unread':
        return <Eye className="h-4 w-4 text-muted-foreground" />;
      case 'read':
        return <EyeOff className="h-4 w-4 text-green-600" />;
      case 'archived':
        return <BookmarkX className="h-4 w-4 text-gray-600" />;
    }
  };

  const getReadingTime = () => {
    const minutes = Math.ceil(article.readingTimeSeconds / 60);
    return `${minutes}分`;
  };

  return (
    <Card className="hover:shadow-lg transition-shadow duration-200">
      <CardHeader className="pb-3">
        <div className="flex justify-between items-start gap-2">
          <Link href={`/articles/${article.id}`} className="flex-1">
            <h3 className="font-semibold text-lg line-clamp-2 hover:text-primary transition-colors">
              {article.title}
            </h3>
          </Link>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => handleStatusChange('unread')}>
                <Eye className="mr-2 h-4 w-4" />
                未読にする
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handleStatusChange('read')}>
                <EyeOff className="mr-2 h-4 w-4" />
                既読にする
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handleStatusChange('archived')}>
                <BookmarkX className="mr-2 h-4 w-4" />
                アーカイブする
              </DropdownMenuItem>
              <DropdownMenuItem 
                onClick={() => onDelete?.(article.id)} 
                className="text-destructive"
              >
                <Trash2 className="mr-2 h-4 w-4" />
                削除する
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </CardHeader>

      <CardContent className="space-y-3">
        {article.thumbnailUrl && (
          <Link href={`/articles/${article.id}`}>
            <div className="relative w-full h-48 rounded-md overflow-hidden bg-muted">
              <img 
                src={article.thumbnailUrl} 
                alt={article.title}
                className="w-full h-full object-cover hover:scale-105 transition-transform duration-300"
              />
            </div>
          </Link>
        )}

        {article.summary && (
          <p className="text-sm text-muted-foreground line-clamp-3">
            {article.summary}
          </p>
        )}

        <div className="flex items-center gap-4 text-xs text-muted-foreground">
          <div className="flex items-center gap-1">
            <Globe className="h-3 w-3" />
            <span>{article.siteName || new URL(article.url).hostname}</span>
          </div>
          {article.author && (
            <span>{article.author}</span>
          )}
          <div className="flex items-center gap-1">
            <Clock className="h-3 w-3" />
            <span>{getReadingTime()}</span>
          </div>
        </div>

        {article.tags.length > 0 && (
          <div className="flex flex-wrap gap-1">
            {article.tags.map(tag => (
              <span 
                key={tag.id}
                className="inline-flex items-center gap-1 px-2 py-1 rounded-full bg-secondary text-secondary-foreground text-xs"
              >
                <Tag className="h-3 w-3" />
                {tag.name}
              </span>
            ))}
          </div>
        )}
      </CardContent>

      <CardFooter className="pt-3 border-t">
        <div className="flex items-center justify-between w-full">
          <div className="flex items-center gap-2">
            {getStatusIcon()}
            {article.category && (
              <span 
                className="px-2 py-1 rounded-full text-xs font-medium"
                style={{ 
                  backgroundColor: `${article.category.color}20`,
                  color: article.category.color 
                }}
              >
                {article.category.name}
              </span>
            )}
          </div>
          <div className="flex items-center gap-2 text-xs text-muted-foreground">
            <span>
              {formatDistanceToNow(new Date(article.savedAt), { 
                addSuffix: true, 
                locale: ja 
              })}
            </span>
            <Button
              variant="ghost"
              size="sm"
              className="h-8 w-8 p-0"
              onClick={() => onToggleFavorite?.(article.id)}
            >
              <Bookmark 
                className={`h-4 w-4 ${article.isFavorite ? 'fill-current text-yellow-500' : ''}`} 
              />
            </Button>
          </div>
        </div>
      </CardFooter>
    </Card>
  );
};