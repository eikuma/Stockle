# P1-102: 認証システムUI実装

## 概要
NextAuth.jsを使用した認証機能のUI実装（ログイン・サインアップ・Google OAuth）

## 担当者
**Member 1 (Frontend Developer)**

## 優先度
**最高** - アプリケーションの基本機能

## 前提条件
- P1-101: フロントエンド基盤構築が完了済み
- Member 2の認証API実装と並行作業可能

## 作業内容

### 1. NextAuth.js設定
- [x] NextAuth.js の設定ファイル作成
- [x] Google OAuth プロバイダーの設定
- [x] JWT + Database セッション戦略の設定
- [x] カスタムページの設定
- [x] セッション管理の設定

### 2. 認証ページの作成
- [x] `/auth/signin` - ログインページ
- [x] `/auth/signup` - サインアップページ
- [ ] `/auth/forgot-password` - パスワードリセットページ（後続タスク）
- [ ] `/auth/reset-password` - 新パスワード設定ページ（後続タスク）
- [ ] `/auth/verify-email` - メール確認ページ（後続タスク）

### 3. 認証フォームコンポーネント
- [x] `SignInForm` - ログインフォーム
- [x] `SignUpForm` - サインアップフォーム
- [ ] `ForgotPasswordForm` - パスワードリセットフォーム（後続タスク）
- [ ] `ResetPasswordForm` - 新パスワード設定フォーム（後続タスク）
- [x] `GoogleSignInButton` - Google認証ボタン（SignInFormに統合）

### 4. 認証状態管理
- [x] `useAuth` カスタムフック
- [x] 認証状態のZustand store
- [x] セッション情報の管理
- [x] 認証エラーハンドリング

### 5. 認証ガード機能
- [x] `AuthGuard` コンポーネント
- [ ] `withAuth` HOC（必要に応じて後続タスク）
- [x] リダイレクト機能
- [x] 未認証ユーザーの処理

### 6. UI/UX の実装
- [x] レスポンシブデザイン
- [x] ローディング状態
- [x] エラー表示
- [x] 成功メッセージ
- [x] フォームバリデーション

## 実装詳細

### app/api/auth/[...nextauth]/route.ts
```typescript
import NextAuth from 'next-auth';
import GoogleProvider from 'next-auth/providers/google';
import CredentialsProvider from 'next-auth/providers/credentials';

const handler = NextAuth({
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID!,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
    }),
    CredentialsProvider({
      name: 'credentials',
      credentials: {
        email: { label: 'Email', type: 'email' },
        password: { label: 'Password', type: 'password' },
      },
      async authorize(credentials) {
        if (!credentials?.email || !credentials?.password) {
          return null;
        }

        // バックエンドAPIでの認証処理
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/login`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            email: credentials.email,
            password: credentials.password,
          }),
        });

        if (!response.ok) {
          return null;
        }

        const user = await response.json();
        return user;
      },
    }),
  ],
  pages: {
    signIn: '/auth/signin',
    signUp: '/auth/signup',
    error: '/auth/error',
  },
  callbacks: {
    async jwt({ token, user, account }) {
      if (user) {
        token.accessToken = user.accessToken;
        token.id = user.id;
      }
      return token;
    },
    async session({ session, token }) {
      session.accessToken = token.accessToken;
      session.user.id = token.id;
      return session;
    },
  },
});

export { handler as GET, handler as POST };
```

### components/auth/SignInForm.tsx
```typescript
'use client';

import { useState } from 'react';
import { signIn } from 'next-auth/react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { AlertCircle, Mail, Lock } from 'lucide-react';

const signInSchema = z.object({
  email: z.string().email('有効なメールアドレスを入力してください'),
  password: z.string().min(8, 'パスワードは8文字以上である必要があります'),
});

type SignInForm = z.infer<typeof signInSchema>;

export function SignInForm() {
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const form = useForm<SignInForm>({
    resolver: zodResolver(signInSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  });

  const onSubmit = async (data: SignInForm) => {
    setIsLoading(true);
    setError(null);

    try {
      const result = await signIn('credentials', {
        email: data.email,
        password: data.password,
        redirect: false,
      });

      if (result?.error) {
        setError('メールアドレスまたはパスワードが正しくありません');
      } else {
        router.push('/dashboard');
      }
    } catch (error) {
      setError('ログインに失敗しました。もう一度お試しください。');
    } finally {
      setIsLoading(false);
    }
  };

  const handleGoogleSignIn = () => {
    signIn('google', { callbackUrl: '/dashboard' });
  };

  return (
    <Card className="w-full max-w-md mx-auto">
      <CardHeader>
        <CardTitle className="text-2xl text-center">ログイン</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {error && (
          <div className="flex items-center space-x-2 text-red-600 text-sm">
            <AlertCircle className="w-4 h-4" />
            <span>{error}</span>
          </div>
        )}

        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="email">メールアドレス</Label>
            <div className="relative">
              <Mail className="absolute left-3 top-3 w-4 h-4 text-gray-400" />
              <Input
                id="email"
                type="email"
                placeholder="your@example.com"
                className="pl-10"
                {...form.register('email')}
              />
            </div>
            {form.formState.errors.email && (
              <p className="text-red-600 text-sm">{form.formState.errors.email.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="password">パスワード</Label>
            <div className="relative">
              <Lock className="absolute left-3 top-3 w-4 h-4 text-gray-400" />
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                className="pl-10"
                {...form.register('password')}
              />
            </div>
            {form.formState.errors.password && (
              <p className="text-red-600 text-sm">{form.formState.errors.password.message}</p>
            )}
          </div>

          <Button
            type="submit"
            className="w-full"
            disabled={isLoading}
          >
            {isLoading ? 'ログイン中...' : 'ログイン'}
          </Button>
        </form>

        <div className="relative">
          <div className="absolute inset-0 flex items-center">
            <span className="w-full border-t" />
          </div>
          <div className="relative flex justify-center text-xs uppercase">
            <span className="bg-background px-2 text-muted-foreground">または</span>
          </div>
        </div>

        <Button
          type="button"
          variant="outline"
          onClick={handleGoogleSignIn}
          className="w-full"
        >
          <svg className="w-4 h-4 mr-2" viewBox="0 0 24 24">
            {/* Google Icon SVG */}
            <path
              fill="currentColor"
              d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
            />
            <path
              fill="currentColor"
              d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
            />
            <path
              fill="currentColor"
              d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
            />
            <path
              fill="currentColor"
              d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
            />
          </svg>
          Googleでログイン
        </Button>

        <div className="text-center text-sm">
          <span className="text-muted-foreground">アカウントをお持ちでない方は </span>
          <a href="/auth/signup" className="text-primary hover:underline">
            こちら
          </a>
        </div>

        <div className="text-center">
          <a
            href="/auth/forgot-password"
            className="text-sm text-muted-foreground hover:underline"
          >
            パスワードを忘れた方はこちら
          </a>
        </div>
      </CardContent>
    </Card>
  );
}
```

### hooks/useAuth.ts
```typescript
'use client';

import { useSession, signOut } from 'next-auth/react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/stores/authStore';
import { useEffect } from 'react';

export function useAuth() {
  const { data: session, status } = useSession();
  const router = useRouter();
  const { user, setUser, clearUser } = useAuthStore();

  useEffect(() => {
    if (session?.user) {
      setUser({
        id: session.user.id,
        email: session.user.email!,
        displayName: session.user.name!,
        profileImageUrl: session.user.image,
        authProvider: 'email', // または 'google'
        createdAt: '',
        updatedAt: '',
      });
    } else {
      clearUser();
    }
  }, [session, setUser, clearUser]);

  const logout = async () => {
    await signOut({ redirect: false });
    clearUser();
    router.push('/auth/signin');
  };

  return {
    user,
    isAuthenticated: !!session,
    isLoading: status === 'loading',
    logout,
  };
}
```

### components/auth/AuthGuard.tsx
```typescript
'use client';

import { useAuth } from '@/hooks/useAuth';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';
import { LoadingSpinner } from '@/components/ui/loading-spinner';

interface AuthGuardProps {
  children: React.ReactNode;
  fallback?: React.ReactNode;
}

export function AuthGuard({ children, fallback }: AuthGuardProps) {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push('/auth/signin');
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    return fallback || <LoadingSpinner />;
  }

  if (!isAuthenticated) {
    return null;
  }

  return <>{children}</>;
}
```

## 受入条件

### 必須条件
- [x] メール/パスワードでのログインが正常に動作する
- [x] Googleアカウントでのログインが正常に動作する（設定済み）
- [x] サインアップフローが完全に機能する
- [ ] パスワードリセット機能が動作する（後続タスク）
- [x] 認証状態が正しく管理される
- [x] 未認証時のリダイレクトが正常に動作する

### 品質条件
- [x] フォームバリデーションが適切に動作する
- [x] エラーメッセージが分かりやすい
- [x] レスポンシブデザインが適用されている
- [x] アクセシビリティが考慮されている
- [x] ローディング状態が適切に表示される

## 推定時間
**24時間** (4-5日)

## 依存関係
- P1-101: フロントエンド基盤構築
- Member 2の認証API実装と連携

## 完了後の次ステップ
1. P1-103: レイアウト・ナビゲーションの実装
2. 認証APIとの統合テスト

## 備考
- NextAuth.jsの設定はセキュリティを最優先に
- Google OAuth の設定は環境変数で管理
- エラーハンドリングを丁寧に実装
- UXを重視したスムーズな認証フロー