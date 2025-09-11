import Link from 'next/link';
import { Button } from '@/components/ui/button';

export default function Home() {
  return (
    <section className="container grid items-center gap-6 pb-8 pt-6 md:py-10">
      <div className="flex max-w-[980px] flex-col items-start gap-2">
        <h1 className="text-3xl font-extrabold leading-tight tracking-tighter md:text-4xl">
          あなたの知識を、もっと身近に。
          <br className="hidden sm:inline" />
          Stockleで、未来の自分へ最高の情報を届けよう。
        </h1>
        <p className="max-w-[700px] text-lg text-muted-foreground">
          Stockleは、気になった記事や情報を簡単に保存し、いつでもどこでもアクセスできるアプリケーションです。
          情報の海から、あなただけの知識の宝箱を作りましょう。
        </p>
      </div>
      <div className="flex gap-4">
        <Link href="/articles">
          <Button variant="default">記事を管理する</Button>
        </Link>
        <Link href="/about">
          <Button variant="outline">Stockleについて</Button>
        </Link>
      </div>
    </section>
  );
}