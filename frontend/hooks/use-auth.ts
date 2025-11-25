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
        if (!AuthService.isAuthenticated()) {
          setIsAuthenticated(false);
          setUser(null);
          setLoading(false);
          return;
        }
        
        // Try to get current user (this will auto-refresh if needed)
        const userData = await AuthService.getCurrentUser();
        
        if (!userData) {
          // Complete auth failure - clear everything
          AuthService.clearAuth();
          setIsAuthenticated(false);
          setUser(null);
          setLoading(false);
          return;
        }

        setUser(userData);
        setIsAuthenticated(true);
      } catch (error) {
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

  const requireEmailVerification = () => {
    // Only check if user is loaded and authenticated
    if (loading) return true; // Still loading, wait
    if (!isAuthenticated || !user) return true; // Not authenticated, let requireAuth handle it
    if (!user.email_status) {
      router.push('/verify-email');
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
    requireEmailVerification,
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
          AuthService.clearAuth();
          setIsAuthenticated(false);
          setUser(null);
          setLoading(false);
          return;
        }

        setUser(userData);
        setIsAuthenticated(true);
      } catch (error) {
        AuthService.clearAuth();
        setIsAuthenticated(false);
        setUser(null);
      } finally {
        setLoading(false);
      }
    }
  };
}
