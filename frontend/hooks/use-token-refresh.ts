import { useEffect } from 'react';
import { AuthService } from '@/lib/auth';

export function useTokenRefresh() {
  useEffect(() => {
    let refreshInterval: NodeJS.Timeout;

    const startTokenRefreshTimer = () => {
      // Check every 4 minutes (240 seconds) if token needs refresh
      refreshInterval = setInterval(async () => {
        try {
          if (AuthService.isAuthenticated()) {
            await AuthService.ensureValidToken();
          }
        } catch (error) {
          console.error('Token refresh failed:', error);
        }
      }, 240000); // 4 minutes
    };

    // Start the timer if user is authenticated
    if (AuthService.isAuthenticated()) {
      startTokenRefreshTimer();
    }

    // Cleanup on unmount
    return () => {
      if (refreshInterval) {
        clearInterval(refreshInterval);
      }
    };
  }, []);

  // Manual refresh function
  const refreshNow = async () => {
    try {
      return await AuthService.refreshToken();
    } catch (error) {
      console.error('Manual token refresh failed:', error);
      return false;
    }
  };

  return { refreshNow };
}
