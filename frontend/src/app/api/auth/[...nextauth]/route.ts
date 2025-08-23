import NextAuth from 'next-auth';
import GoogleProvider from 'next-auth/providers/google';

const handler = NextAuth({
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID!,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
    }),
  ],
  pages: {
    signIn: '/auth/signin',
    error: '/auth/error',
  },
  callbacks: {
    async jwt({ token, user, account }) {
      // Google OAuth 専用処理
      if (account?.provider === "google" && user) {
        console.log('🔑 Google OAuth detected, fetching JWT from backend');
        console.log('🔑 User data:', { email: user?.email, name: user?.name, image: user?.image });
        console.log('🔑 Account data:', { provider: account.provider, providerAccountId: account.providerAccountId });
        
        // サーバーサイドでは内部API URLを使用
        const apiUrl = process.env.INTERNAL_API_URL || process.env.NEXT_PUBLIC_API_URL;
        console.log('🔑 API URL (client):', process.env.NEXT_PUBLIC_API_URL);
        console.log('🔑 API URL (internal):', process.env.INTERNAL_API_URL);
        console.log('🔑 API URL (using):', apiUrl);
        
        try {
          const requestBody = {
            email: user.email || '',
            name: user.name || '',
            google_id: account.providerAccountId,
            image_url: user.image || ''
          };
          console.log('🔑 Request body:', requestBody);

          const response = await fetch(`${apiUrl}/auth/google`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(requestBody)
          });

          console.log('🔑 Response status:', response.status);
          console.log('🔑 Response headers:', Object.fromEntries(response.headers.entries()));

          if (response.ok) {
            const data = await response.json();
            console.log('🔑 Google Auth backend response:', data);
            token.accessToken = data.tokens.access_token;
            token.refreshToken = data.tokens.refresh_token;
            token.id = data.user.id;
            console.log('🔑 JWT tokens set successfully:', { 
              hasAccessToken: !!token.accessToken, 
              hasRefreshToken: !!token.refreshToken,
              userId: token.id
            });
          } else {
            const errorText = await response.text();
            console.error('🔑 Failed to get JWT from backend for Google auth');
            console.error('🔑 Response status:', response.status);
            console.error('🔑 Response text:', errorText);
          }
        } catch (error) {
          console.error('🔑 Error fetching JWT for Google auth:', error);
          console.error('🔑 Error details:', {
            name: error.name,
            message: error.message,
            stack: error.stack
          });
        }
      }
      
      console.log('🔑 Final token state:', {
        hasAccessToken: !!token.accessToken,
        hasRefreshToken: !!token.refreshToken,
        userId: token.id
      });
      return token;
    },
    async session({ session, token }) {
      if (session.user) {
        session.accessToken = token.accessToken;
        session.user.id = token.id as string;
      }
      return session;
    },
  },
  session: {
    strategy: 'jwt',
  },
  secret: process.env.NEXTAUTH_SECRET,
});

export { handler as GET, handler as POST };