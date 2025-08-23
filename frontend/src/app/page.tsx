'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { AuthGuard } from '@/components/auth/AuthGuard';
import { useAuth } from '@/hooks/useAuth';
import { LoadingSpinner } from '@/components/ui/loading-spinner';
import { 
  BookOpen, 
  Plus, 
  Eye, 
  Star, 
  TrendingUp,
  Calendar,
  Search,
  Settings
} from 'lucide-react';

interface DashboardStats {
  total_articles: number;
  unread_articles: number;
  favorite_articles: number;
  recent_articles: Array<{
    id: string;
    title: string;
    url: string;
    saved_at: string;
    status: 'read' | 'unread';
  }>;
}

export default function HomePage() {
  const { user } = useAuth();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        // 実際のAPIが実装されるまではモックデータを使用
        await new Promise(resolve => setTimeout(resolve, 1000)); // 1秒待機
        
        const mockStats: DashboardStats = {
          total_articles: 24,
          unread_articles: 8,
          favorite_articles: 5,
          recent_articles: [
            {
              id: '1',
              title: 'Next.js 14の新機能について',
              url: 'https://example.com/nextjs-14',
              saved_at: '2025-08-11T10:30:00Z',
              status: 'unread'
            },
            {
              id: '2',
              title: 'React Serverコンポーネントの活用法',
              url: 'https://example.com/react-server-components',
              saved_at: '2025-08-10T15:20:00Z',
              status: 'read'
            },
            {
              id: '3',
              title: 'TypeScriptでの型安全なAPI設計',
              url: 'https://example.com/typescript-api',
              saved_at: '2025-08-09T09:15:00Z',
              status: 'unread'
            }
          ]
        };
        
        setStats(mockStats);
      } catch (error) {
        console.error('Failed to fetch dashboard data:', error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchDashboardData();
  }, []);

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('ja-JP', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (isLoading) {
    return (
      <AuthGuard>
        <div className="container mx-auto px-4 py-8">
          <div className="flex justify-center">
            <LoadingSpinner />
          </div>
        </div>
      </AuthGuard>
    );
  }

  return (
    <AuthGuard>
      <div className="container mx-auto px-4 py-8">
        {/* ウェルカムセクション */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-2">
            おかえりなさい、{user?.displayName || 'ユーザー'}さん！
          </h1>
          <p className="text-gray-600">
            今日も新しい知識を蓄えましょう📚
          </p>
        </div>

        {/* 統計カード */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">保存済み記事</CardTitle>
              <BookOpen className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats?.total_articles || 0}</div>
              <p className="text-xs text-muted-foreground">
                あなたの知識ライブラリ
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">未読記事</CardTitle>
              <Eye className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-blue-600">
                {stats?.unread_articles || 0}
              </div>
              <p className="text-xs text-muted-foreground">
                読み待ちの記事
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">お気に入り</CardTitle>
              <Star className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-yellow-600">
                {stats?.favorite_articles || 0}
              </div>
              <p className="text-xs text-muted-foreground">
                特別な記事たち
              </p>
            </CardContent>
          </Card>
        </div>

        {/* クイックアクション */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
          <div className="lg:col-span-2">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Calendar className="h-5 w-5" />
                  最近保存した記事
                </CardTitle>
              </CardHeader>
              <CardContent>
                {stats?.recent_articles && stats.recent_articles.length > 0 ? (
                  <div className="space-y-3">
                    {stats.recent_articles.map((article) => (
                      <div
                        key={article.id}
                        className="flex items-center justify-between p-3 border rounded-lg hover:bg-gray-50 transition-colors"
                      >
                        <div className="flex-1 min-w-0">
                          <h4 className="text-sm font-medium truncate">
                            {article.title}
                          </h4>
                          <p className="text-xs text-gray-500">
                            {formatDate(article.saved_at)}
                          </p>
                        </div>
                        <div className="flex items-center gap-2 ml-4">
                          {article.status === 'unread' && (
                            <span className="inline-block w-2 h-2 bg-blue-500 rounded-full"></span>
                          )}
                          <Link href={`/articles/${article.id}`}>
                            <Button variant="ghost" size="sm">
                              詳細
                            </Button>
                          </Link>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <p className="text-gray-500 text-center py-8">
                    まだ記事が保存されていません
                  </p>
                )}
                <div className="mt-4 flex justify-center">
                  <Link href="/articles">
                    <Button variant="outline">
                      すべての記事を見る
                    </Button>
                  </Link>
                </div>
              </CardContent>
            </Card>
          </div>

          <div>
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <TrendingUp className="h-5 w-5" />
                  クイックアクション
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <Link href="/articles" className="block">
                  <Button className="w-full justify-start" size="sm">
                    <Plus className="mr-2 h-4 w-4" />
                    新しい記事を保存
                  </Button>
                </Link>
                
                <Link href="/articles" className="block">
                  <Button variant="outline" className="w-full justify-start" size="sm">
                    <Search className="mr-2 h-4 w-4" />
                    記事を検索
                  </Button>
                </Link>

                <Link href="/articles?status=unread" className="block">
                  <Button variant="outline" className="w-full justify-start" size="sm">
                    <Eye className="mr-2 h-4 w-4" />
                    未読記事を読む
                  </Button>
                </Link>

                <Link href="/settings" className="block">
                  <Button variant="ghost" className="w-full justify-start" size="sm">
                    <Settings className="mr-2 h-4 w-4" />
                    設定
                  </Button>
                </Link>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* ヒント */}
        <Card className="bg-blue-50 border-blue-200">
          <CardContent className="pt-6">
            <div className="flex items-start gap-3">
              <div className="text-blue-600">💡</div>
              <div>
                <h3 className="font-semibold text-blue-900 mb-1">
                  Stockleを最大限活用しよう！
                </h3>
                <p className="text-sm text-blue-800">
                  記事を保存すると、AIが自動で要約を生成します。忙しい時でも大切な情報を素早く把握できます。
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </AuthGuard>
  );
}