import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { AuthService, User } from '@/lib/auth';

export function useAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const router = useRouter();

  useEffect(() => {
    const checkAuthAndFetchUser = async () => {
      try {
        console.log('useAuth: Starting auth check');
        if (!AuthService.isAuthenticated()) {
          console.log('useAuth: Not authenticated');
          setIsAuthenticated(false);
          setUser(null);
          setLoading(false);
          return;
        }

        console.log('useAuth: Token is valid, fetching user data');
        
        // Try to get current user (this will auto-refresh if needed)
        const userData = await AuthService.getCurrentUser();
        
        if (!userData) {
          console.log('useAuth: Failed to get user data, clearing auth');
          // Complete auth failure - clear everything
          AuthService.clearAuth();
          setIsAuthenticated(false);
          setUser(null);
          setLoading(false);
          return;
        }

        console.log('useAuth: User data received:', userData);
        setUser(userData);
        setIsAuthenticated(true);
      } catch (error) {
        console.error('useAuth: Auth check failed:', error);
        AuthService.clearAuth();
        setIsAuthenticated(false);
        setUser(null);
      } finally {
        setLoading(false);
      }
    };

    checkAuthAndFetchUser();
  }, []);

  const logout = () => {
    AuthService.logout();
    setUser(null);
    setIsAuthenticated(false);
    router.push('/');
  };

  const requireAuth = () => {
    if (!loading && !isAuthenticated) {
      router.push('/auth/signin');
      return false;
    }
    return true;
  };

  return {
    user,
    loading,
    isAuthenticated,
    logout,
    requireAuth,
    refetch: async () => {
      setLoading(true);
      try {
        if (!AuthService.isAuthenticated()) {
          setIsAuthenticated(false);
          setUser(null);
          setLoading(false);
          return;
        }

        // Try to get current user (this will auto-refresh if needed)
        const userData = await AuthService.getCurrentUser();
        
        if (!userData) {
          console.log('useAuth refetch: Failed to get user data, clearing auth');
          AuthService.clearAuth();
          setIsAuthenticated(false);
          setUser(null);
          setLoading(false);
          return;
        }

        setUser(userData);
        setIsAuthenticated(true);
      } catch (error) {
        console.error('Auth refetch failed:', error);
        AuthService.clearAuth();
        setIsAuthenticated(false);
        setUser(null);
      } finally {
        setLoading(false);
      }
    }
  };
}
