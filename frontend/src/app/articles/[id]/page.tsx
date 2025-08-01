'use client';

import React, { useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { formatDistanceToNow } from 'date-fns';
import { ja } from 'date-fns/locale';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { ArrowLeft, Bookmark, Clock, ExternalLink, Globe, Share2, Tag } from 'lucide-react';
import { mockArticles } from '@/lib/mock/articleData';
import type { Article } from '@/types/article';

export default function ArticleDetailPage() {
  const params = useParams();
  const router = useRouter();
  const articleId = params.id as string;
  
  const [article, setArticle] = useState<Article | null>(
    mockArticles.find(a => a.id === articleId) || null
  );

  if (!article) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center py-12">
          <h1 className="text-2xl font-bold mb-2">記事が見つかりません</h1>
          <p className="text-muted-foreground mb-4">
            指定された記事は存在しないか、削除された可能性があります。
          </p>
          <Button onClick={() => router.push('/articles')}>
            記事一覧に戻る
          </Button>
        </div>
      </div>
    );
  }

  const handleToggleFavorite = () => {
    setArticle(prev => prev ? { ...prev, isFavorite: !prev.isFavorite } : null);
  };

  const handleShare = async () => {
    if (navigator.share) {
      try {
        await navigator.share({
          title: article.title,
          text: article.summary,
          url: window.location.href,
        });
      } catch (error) {
        console.log('シェアがキャンセルされました');
      }
    } else {
      // フォールバック: URLをクリップボードにコピー
      await navigator.clipboard.writeText(window.location.href);
      alert('URLをクリップボードにコピーしました');
    }
  };

  const getReadingTime = () => {
    const minutes = Math.ceil(article.readingTimeSeconds / 60);
    return `${minutes}分`;
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="space-y-6">
        {/* ナビゲーション */}
        <div className="flex items-center gap-4">
          <Button 
            variant="ghost" 
            size="sm" 
            onClick={() => router.push('/articles')}
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            記事一覧に戻る
          </Button>

          <div className="flex items-center gap-2 ml-auto">
            <Button
              variant="ghost"
              size="sm"
              onClick={handleToggleFavorite}
            >
              <Bookmark 
                className={`h-4 w-4 ${article.isFavorite ? 'fill-current text-yellow-500' : ''}`} 
              />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={handleShare}
            >
              <Share2 className="h-4 w-4" />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => window.open(article.url, '_blank')}
            >
              <ExternalLink className="h-4 w-4" />
              元記事を開く
            </Button>
          </div>
        </div>

        {/* 記事ヘッダー */}
        <div className="space-y-4">
          <h1 className="text-3xl font-bold leading-tight">
            {article.title}
          </h1>

          <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
            <div className="flex items-center gap-1">
              <Globe className="h-4 w-4" />
              <span>{article.siteName}</span>
            </div>
            {article.author && (
              <span>{article.author}</span>
            )}
            <div className="flex items-center gap-1">
              <Clock className="h-4 w-4" />
              <span>{getReadingTime()}</span>
            </div>
            <span>
              {formatDistanceToNow(new Date(article.savedAt), { 
                addSuffix: true, 
                locale: ja 
              })}に保存
            </span>
          </div>

          {/* カテゴリとタグ */}
          <div className="flex flex-wrap items-center gap-2">
            {article.category && (
              <Badge 
                variant="secondary"
                style={{ 
                  backgroundColor: `${article.category.color}20`,
                  color: article.category.color 
                }}
              >
                {article.category.name}
              </Badge>
            )}
            {article.tags.map(tag => (
              <Badge key={tag.id} variant="outline" className="flex items-center gap-1">
                <Tag className="h-3 w-3" />
                {tag.name}
              </Badge>
            ))}
          </div>

          {/* 読み進み状況 */}
          {article.readingProgress > 0 && (
            <div className="space-y-2">
              <div className="flex items-center justify-between text-sm">
                <span>読み進み状況</span>
                <span>{Math.round(article.readingProgress * 100)}%</span>
              </div>
              <div className="w-full bg-muted rounded-full h-2">
                <div 
                  className="bg-primary h-2 rounded-full transition-all duration-300" 
                  style={{ width: `${article.readingProgress * 100}%` }}
                />
              </div>
            </div>
          )}
        </div>

        {/* サムネイル */}
        {article.thumbnailUrl && (
          <div className="relative w-full h-64 sm:h-80 rounded-lg overflow-hidden bg-muted">
            <img 
              src={article.thumbnailUrl} 
              alt={article.title}
              className="w-full h-full object-cover"
            />
          </div>
        )}

        {/* 要約 */}
        {article.summary && (
          <div className="bg-muted/50 p-6 rounded-lg">
            <h2 className="text-lg font-semibold mb-3">要約</h2>
            <p className="text-muted-foreground leading-relaxed">
              {article.summaryLong || article.summary}
            </p>
          </div>
        )}

        {/* 記事内容 */}
        <div className="prose prose-gray max-w-none">
          <h2 className="text-xl font-semibold mb-4">記事内容</h2>
          <div className="bg-muted/30 p-6 rounded-lg">
            <p className="text-muted-foreground mb-4">
              {article.content || 'この記事の本文は現在読み込み中です。実際のAPIでは、Webスクレイピングによって記事の全文が取得されます。'}
            </p>
            <p className="text-sm text-muted-foreground italic">
              注意: これはモックデータです。実際のアプリケーションでは、ここに記事の本文が表示されます。
            </p>
          </div>
        </div>

        {/* 記事メタデータ */}
        <div className="border-t pt-6">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div>
              <div className="font-medium">文字数</div>
              <div className="text-muted-foreground">
                {article.wordCount?.toLocaleString() || '不明'}文字
              </div>
            </div>
            <div>
              <div className="font-medium">読了時間</div>
              <div className="text-muted-foreground">
                {getReadingTime()}
              </div>
            </div>
            <div>
              <div className="font-medium">言語</div>
              <div className="text-muted-foreground">
                {article.language === 'ja' ? '日本語' : article.language}
              </div>
            </div>
            <div>
              <div className="font-medium">公開日</div>
              <div className="text-muted-foreground">
                {article.publishedAt 
                  ? new Date(article.publishedAt).toLocaleDateString('ja-JP')
                  : '不明'
                }
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}