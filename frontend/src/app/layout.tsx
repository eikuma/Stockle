import type { Metadata } from 'next';
import './globals.css';

export const metadata: Metadata = {
  title: 'Stockle',
  description: 'Save articles for later reading',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  );
}