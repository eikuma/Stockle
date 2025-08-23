'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { BookOpen, Home, Settings, User, LogOut } from 'lucide-react';
import { useAuth } from '@/hooks/useAuth';

export const Navigation: React.FC = () => {
  const pathname = usePathname();
  const { isAuthenticated, logout } = useAuth();

  const navItems = [
    { href: '/', label: 'ホーム', icon: Home },
    { href: '/articles', label: '記事管理', icon: BookOpen },
    { href: '/settings', label: '設定', icon: Settings },
    { href: '/profile', label: 'プロフィール', icon: User },
  ];

  return (
    <nav className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto px-4">
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center gap-6">
            <Link href="/" className="text-2xl font-bold text-primary">
              Stockle
            </Link>
            <div className="hidden md:flex items-center gap-1">
              {navItems.map(item => {
                const Icon = item.icon;
                const isActive = pathname === item.href;
                return (
                  <Link key={item.href} href={item.href}>
                    <Button 
                      variant={isActive ? 'default' : 'ghost'} 
                      size="sm"
                      className="flex items-center gap-2"
                    >
                      <Icon className="h-4 w-4" />
                      {item.label}
                    </Button>
                  </Link>
                );
              })}
            </div>
          </div>
          
          <div className="flex items-center gap-4">
            {isAuthenticated ? (
              <Button 
                variant="outline" 
                size="sm"
                onClick={logout}
                className="flex items-center gap-2"
              >
                <LogOut className="h-4 w-4" />
                ログアウト
              </Button>
            ) : (
              <Link href="/auth/signin">
                <Button variant="default" size="sm">
                  ログイン
                </Button>
              </Link>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
};