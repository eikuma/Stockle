'use client';

import { useState, useEffect } from 'react';
import { useSession } from 'next-auth/react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { AuthGuard } from '@/components/auth/AuthGuard';
import { useAuth } from '@/hooks/useAuth';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { LoadingSpinner } from '@/components/ui/loading-spinner';
import { AlertCircle, CheckCircle, User } from 'lucide-react';

const profileSchema = z.object({
  display_name: z
    .string()
    .min(1, '表示名を入力してください')
    .max(50, '表示名は50文字以内で入力してください'),
  email: z.string().email('有効なメールアドレスを入力してください'),
});

type ProfileForm = z.infer<typeof profileSchema>;

interface UserProfile {
  id: string;
  email: string;
  display_name: string;
  profile_image_url?: string;
  auth_provider: string;
}

export default function ProfilePage() {
  const { user, isAuthenticated } = useAuth();
  const { data: session } = useSession();
  const [isLoading, setIsLoading] = useState(true);
  const [isUpdating, setIsUpdating] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
  const [userProfile, setUserProfile] = useState<UserProfile | null>(null);

  const form = useForm<ProfileForm>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      display_name: '',
      email: '',
    },
  });

  useEffect(() => {
    const fetchUserProfile = async () => {
      if (!isAuthenticated) return;

      console.log('🔍 NEXT_PUBLIC_API_URL:', process.env.NEXT_PUBLIC_API_URL);
      console.log('🔍 NextAuth Session:', session);
      console.log('🔍 Session AccessToken:', session?.accessToken);
      console.log('🔍 LocalStorage Token:', localStorage.getItem('access_token'));
      
      try {
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/users/me`, {
          headers: {
            'Authorization': `Bearer ${session?.accessToken || ''}`,
          },
        });

        if (response.ok) {
          const profile = await response.json();
          setUserProfile(profile);
          form.reset({
            display_name: profile.display_name,
            email: profile.email,
          });
        }
      } catch (error) {
        console.error('Failed to fetch user profile:', error);
        setMessage({ type: 'error', text: 'プロフィール情報の取得に失敗しました' });
      } finally {
        setIsLoading(false);
      }
    };

    fetchUserProfile();
  }, [isAuthenticated, form, session]);

  const onSubmit = async (data: ProfileForm) => {
    setIsUpdating(true);
    setMessage(null);

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/users/me`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session?.accessToken || ''}`,
        },
        body: JSON.stringify({
          display_name: data.display_name,
          email: data.email,
        }),
      });

      if (response.ok) {
        const updatedProfile = await response.json();
        setUserProfile(updatedProfile);
        setMessage({ type: 'success', text: 'プロフィールを更新しました' });
      } else {
        throw new Error('Failed to update profile');
      }
    } catch (error) {
      console.error('Failed to update profile:', error);
      setMessage({ type: 'error', text: 'プロフィールの更新に失敗しました' });
    } finally {
      setIsUpdating(false);
    }
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
        <div className="max-w-2xl mx-auto">
          <Card>
            <CardHeader>
              <CardTitle className="text-2xl">プロフィール設定</CardTitle>
            </CardHeader>
            <CardContent className="space-y-6">
              {message && (
                <div className={`flex items-center space-x-2 p-3 rounded-md text-sm ${
                  message.type === 'success'
                    ? 'bg-green-100 text-green-700'
                    : 'bg-red-100 text-red-700'
                }`}>
                  {message.type === 'success' ? (
                    <CheckCircle className="w-4 h-4" />
                  ) : (
                    <AlertCircle className="w-4 h-4" />
                  )}
                  <span>{message.text}</span>
                </div>
              )}

              {/* プロフィール画像 */}
              <div className="flex items-center space-x-4">
                <Avatar className="w-20 h-20">
                  <AvatarImage 
                    src={userProfile?.profile_image_url || user?.profileImageUrl} 
                    alt="プロフィール画像" 
                  />
                  <AvatarFallback>
                    <User className="w-8 h-8" />
                  </AvatarFallback>
                </Avatar>
                <div>
                  <h3 className="text-lg font-semibold">プロフィール画像</h3>
                  <p className="text-sm text-gray-600">
                    {userProfile?.auth_provider === 'google' 
                      ? 'Googleアカウントから取得されています'
                      : 'デフォルトアバターを使用中'
                    }
                  </p>
                </div>
              </div>

              {/* プロフィールフォーム */}
              <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="display_name">表示名</Label>
                  <Input
                    id="display_name"
                    placeholder="表示名を入力してください"
                    {...form.register('display_name')}
                  />
                  {form.formState.errors.display_name && (
                    <p className="text-red-600 text-sm">
                      {form.formState.errors.display_name.message}
                    </p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="email">メールアドレス</Label>
                  <Input
                    id="email"
                    type="email"
                    placeholder="メールアドレスを入力してください"
                    {...form.register('email')}
                  />
                  {form.formState.errors.email && (
                    <p className="text-red-600 text-sm">
                      {form.formState.errors.email.message}
                    </p>
                  )}
                  <p className="text-sm text-gray-600">
                    メールアドレスを変更した場合、再度認証が必要になります
                  </p>
                </div>

                <div className="flex justify-end">
                  <Button 
                    type="submit" 
                    disabled={isUpdating}
                    className="w-full sm:w-auto"
                  >
                    {isUpdating ? '更新中...' : 'プロフィールを更新'}
                  </Button>
                </div>
              </form>

              {/* アカウント情報 */}
              <div className="border-t pt-6">
                <h3 className="text-lg font-semibold mb-4">アカウント情報</h3>
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-600">ユーザーID:</span>
                    <span>{userProfile?.id || 'N/A'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-600">認証方法:</span>
                    <span>
                      {userProfile?.auth_provider === 'google' ? 'Google' : 'メール認証'}
                    </span>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </AuthGuard>
  );
}