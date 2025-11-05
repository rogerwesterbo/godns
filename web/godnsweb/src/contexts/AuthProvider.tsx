import { useEffect, useState } from 'react';
import type { ReactNode } from 'react';
import { jwtDecode } from 'jwt-decode';
import * as authService from '../services/auth';
import { AuthContext } from './AuthContext';
import type { User, AuthContextType } from './AuthContext';

interface TokenPayload {
  realm_access?: {
    roles: string[];
  };
}

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    initializeAuth();
  }, []);

  const initializeAuth = async () => {
    try {
      const token = await authService.getValidAccessToken();

      if (token) {
        const decoded = jwtDecode<User>(token);

        // Extract roles from realm_access
        const tokenData = jwtDecode<TokenPayload>(token);
        const roles = tokenData.realm_access?.roles || [];

        setUser({
          ...decoded,
          roles,
        });
      }
    } catch (error) {
      console.error('Failed to initialize auth:', error);
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  };

  const login = async () => {
    try {
      const authUrl = await authService.buildAuthorizationUrl();
      window.location.href = authUrl;
    } catch (error) {
      console.error('Login failed:', error);
    }
  };

  const logout = async () => {
    const idToken = authService.getIdToken();
    await authService.logout(idToken || undefined);
  };

  const getAccessToken = async (): Promise<string | null> => {
    return await authService.getValidAccessToken();
  };

  const value: AuthContextType = {
    user,
    isAuthenticated: !!user,
    isLoading,
    login,
    logout,
    getAccessToken,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
