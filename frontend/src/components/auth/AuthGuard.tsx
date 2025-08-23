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

  console.log('🔒 AuthGuard State:', { isAuthenticated, isLoading });

  useEffect(() => {
    console.log('🔒 AuthGuard useEffect:', { isAuthenticated, isLoading });
    if (!isLoading && !isAuthenticated) {
      console.log('🔒 Redirecting to /auth/signin...');
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