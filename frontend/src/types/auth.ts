export interface User {
  id: string;
  email: string;
  displayName: string;
  profileImageUrl?: string;
  authProvider: 'email' | 'google';
  createdAt: string;
  updatedAt: string;
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
}