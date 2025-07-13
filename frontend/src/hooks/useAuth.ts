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
        profileImageUrl: session.user.image || undefined,
        authProvider: 'email', // TODO: 実際のプロバイダーを判定
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