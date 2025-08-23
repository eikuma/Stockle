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
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { AuthGuard } from '@/components/auth/AuthGuard';
import { useAuth } from '@/hooks/useAuth';
import { LoadingSpinner } from '@/components/ui/loading-spinner';
import { AlertCircle, CheckCircle, Trash2, Shield, Eye } from 'lucide-react';
import { Separator } from '@/components/ui/separator';

const appSettingsSchema = z.object({
  items_per_page: z.number().min(10).max(100),
  default_sort: z.enum(['created_at', 'saved_at', 'title']),
  default_view: z.enum(['grid', 'list']),
});

const passwordSchema = z.object({
  current_password: z.string().min(1, '現在のパスワードを入力してください'),
  new_password: z.string().min(8, 'パスワードは8文字以上で入力してください'),
  confirm_password: z.string().min(1, 'パスワード確認を入力してください'),
}).refine((data) => data.new_password === data.confirm_password, {
  message: 'パスワードが一致しません',
  path: ['confirm_password'],
});

type AppSettingsForm = z.infer<typeof appSettingsSchema>;
type PasswordForm = z.infer<typeof passwordSchema>;

interface UserSettings {
  preferences: {
    items_per_page: number;
    default_sort: string;
    default_view: string;
  };
  auth_provider: string;
}

export default function SettingsPage() {
  const { user, logout } = useAuth();
  const { data: session } = useSession();
  const [isLoading, setIsLoading] = useState(true);
  const [isUpdatingSettings, setIsUpdatingSettings] = useState(false);
  const [isChangingPassword, setIsChangingPassword] = useState(false);
  const [isDeletingAccount, setIsDeletingAccount] = useState(false);
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);
  const [userSettings, setUserSettings] = useState<UserSettings | null>(null);
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [deleteConfirmText, setDeleteConfirmText] = useState('');

  const settingsForm = useForm<AppSettingsForm>({
    resolver: zodResolver(appSettingsSchema),
    defaultValues: {
      items_per_page: 20,
      default_sort: 'saved_at',
      default_view: 'grid',
    },
  });

  const passwordForm = useForm<PasswordForm>({
    resolver: zodResolver(passwordSchema),
    defaultValues: {
      current_password: '',
      new_password: '',
      confirm_password: '',
    },
  });

  useEffect(() => {
    const fetchUserSettings = async () => {
      console.log('🔍 NEXT_PUBLIC_API_URL:', process.env.NEXT_PUBLIC_API_URL);
      console.log('🔍 NextAuth Session:', session);
      console.log('🔍 Session AccessToken:', session?.accessToken);
      
      try {
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/users/me`, {
          headers: {
            'Authorization': `Bearer ${session?.accessToken || ''}`,
          },
        });

        if (response.ok) {
          const data = await response.json();
          setUserSettings(data);
          if (data.preferences) {
            settingsForm.reset({
              items_per_page: data.preferences.items_per_page || 20,
              default_sort: data.preferences.default_sort || 'saved_at',
              default_view: data.preferences.default_view || 'grid',
            });
          }
        }
      } catch (error) {
        console.error('Failed to fetch user settings:', error);
        setMessage({ type: 'error', text: '設定の取得に失敗しました' });
      } finally {
        setIsLoading(false);
      }
    };

    fetchUserSettings();
  }, [settingsForm, session]);

  const onUpdateSettings = async (data: AppSettingsForm) => {
    setIsUpdatingSettings(true);
    setMessage(null);

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/users/me`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session?.accessToken || ''}`,
        },
        body: JSON.stringify({
          preferences: {
            items_per_page: data.items_per_page,
            default_sort: data.default_sort,
            default_view: data.default_view,
          },
        }),
      });

      if (response.ok) {
        setMessage({ type: 'success', text: '設定を更新しました' });
      } else {
        throw new Error('Failed to update settings');
      }
    } catch (error) {
      console.error('Failed to update settings:', error);
      setMessage({ type: 'error', text: '設定の更新に失敗しました' });
    } finally {
      setIsUpdatingSettings(false);
    }
  };

  const onChangePassword = async (data: PasswordForm) => {
    setIsChangingPassword(true);
    setMessage(null);

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/users/me/change-password`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${session?.accessToken || ''}`,
        },
        body: JSON.stringify({
          current_password: data.current_password,
          new_password: data.new_password,
        }),
      });

      if (response.ok) {
        passwordForm.reset();
        setMessage({ type: 'success', text: 'パスワードを変更しました' });
      } else {
        throw new Error('Failed to change password');
      }
    } catch (error) {
      console.error('Failed to change password:', error);
      setMessage({ type: 'error', text: 'パスワードの変更に失敗しました' });
    } finally {
      setIsChangingPassword(false);
    }
  };

  const onDeleteAccount = async () => {
    if (deleteConfirmText !== 'DELETE') {
      setMessage({ type: 'error', text: '削除確認テキストが正しくありません' });
      return;
    }

    setIsDeletingAccount(true);

    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/users/me`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${session?.accessToken || ''}`,
        },
      });

      if (response.ok) {
        setMessage({ type: 'success', text: 'アカウントを削除しました' });
        setTimeout(() => {
          logout();
        }, 2000);
      } else {
        throw new Error('Failed to delete account');
      }
    } catch (error) {
      console.error('Failed to delete account:', error);
      setMessage({ type: 'error', text: 'アカウントの削除に失敗しました' });
    } finally {
      setIsDeletingAccount(false);
      setShowDeleteDialog(false);
      setDeleteConfirmText('');
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
        <div className="max-w-2xl mx-auto space-y-6">
          <h1 className="text-3xl font-bold">設定</h1>

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

          {/* アプリケーション設定 */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Eye className="w-5 h-5" />
                <span>表示設定</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={settingsForm.handleSubmit(onUpdateSettings)} className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="items_per_page">1ページの表示件数</Label>
                    <Select 
                      value={settingsForm.watch('items_per_page')?.toString()} 
                      onValueChange={(value) => settingsForm.setValue('items_per_page', parseInt(value))}
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="10">10件</SelectItem>
                        <SelectItem value="20">20件</SelectItem>
                        <SelectItem value="30">30件</SelectItem>
                        <SelectItem value="50">50件</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="default_sort">デフォルトの並び順</Label>
                    <Select 
                      value={settingsForm.watch('default_sort')} 
                      onValueChange={(value) => settingsForm.setValue('default_sort', value as any)}
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="saved_at">保存日時</SelectItem>
                        <SelectItem value="created_at">作成日時</SelectItem>
                        <SelectItem value="title">タイトル</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="default_view">ビューモード</Label>
                    <Select 
                      value={settingsForm.watch('default_view')} 
                      onValueChange={(value) => settingsForm.setValue('default_view', value as any)}
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="grid">グリッド表示</SelectItem>
                        <SelectItem value="list">リスト表示</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>

                <div className="flex justify-end">
                  <Button type="submit" disabled={isUpdatingSettings}>
                    {isUpdatingSettings ? '更新中...' : '設定を保存'}
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>

          {/* パスワード変更 */}
          {userSettings?.auth_provider !== 'google' && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Shield className="w-5 h-5" />
                  <span>パスワード変更</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <form onSubmit={passwordForm.handleSubmit(onChangePassword)} className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="current_password">現在のパスワード</Label>
                    <Input
                      id="current_password"
                      type="password"
                      {...passwordForm.register('current_password')}
                    />
                    {passwordForm.formState.errors.current_password && (
                      <p className="text-red-600 text-sm">
                        {passwordForm.formState.errors.current_password.message}
                      </p>
                    )}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="new_password">新しいパスワード</Label>
                    <Input
                      id="new_password"
                      type="password"
                      {...passwordForm.register('new_password')}
                    />
                    {passwordForm.formState.errors.new_password && (
                      <p className="text-red-600 text-sm">
                        {passwordForm.formState.errors.new_password.message}
                      </p>
                    )}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="confirm_password">パスワード確認</Label>
                    <Input
                      id="confirm_password"
                      type="password"
                      {...passwordForm.register('confirm_password')}
                    />
                    {passwordForm.formState.errors.confirm_password && (
                      <p className="text-red-600 text-sm">
                        {passwordForm.formState.errors.confirm_password.message}
                      </p>
                    )}
                  </div>

                  <div className="flex justify-end">
                    <Button type="submit" disabled={isChangingPassword}>
                      {isChangingPassword ? '変更中...' : 'パスワードを変更'}
                    </Button>
                  </div>
                </form>
              </CardContent>
            </Card>
          )}

          {/* アカウント削除 */}
          <Card className="border-red-200">
            <CardHeader>
              <CardTitle className="flex items-center space-x-2 text-red-600">
                <Trash2 className="w-5 h-5" />
                <span>危険な操作</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <h3 className="text-lg font-semibold text-red-600">アカウント削除</h3>
                  <p className="text-sm text-gray-600">
                    アカウントを削除すると、保存した全ての記事とデータが永久に削除されます。
                    この操作は取り消すことができません。
                  </p>
                </div>

                <Dialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
                  <DialogTrigger asChild>
                    <Button variant="destructive">
                      アカウントを削除
                    </Button>
                  </DialogTrigger>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>アカウント削除の確認</DialogTitle>
                      <DialogDescription>
                        本当にアカウントを削除しますか？この操作は取り消すことができません。
                        削除を実行するには、下のフィールドに「DELETE」と入力してください。
                      </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4">
                      <Input
                        placeholder="DELETE"
                        value={deleteConfirmText}
                        onChange={(e) => setDeleteConfirmText(e.target.value)}
                      />
                    </div>
                    <DialogFooter>
                      <Button variant="outline" onClick={() => {
                        setShowDeleteDialog(false);
                        setDeleteConfirmText('');
                      }}>
                        キャンセル
                      </Button>
                      <Button
                        variant="destructive"
                        onClick={onDeleteAccount}
                        disabled={isDeletingAccount || deleteConfirmText !== 'DELETE'}
                      >
                        {isDeletingAccount ? '削除中...' : 'アカウントを削除'}
                      </Button>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </AuthGuard>
  );
}