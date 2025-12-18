/**
 * Migration Utility: localStorage to sessionStorage
 * 
 * This utility helps migrate users from the old localStorage token storage
 * to the new secure in-memory + sessionStorage approach.
 * 
 * Run this once when the app initializes to migrate existing users.
 */

export function migrateTokenStorage(): void {
    if (typeof window === 'undefined') return;

    try {
        // Check if old tokens exist in localStorage
        const oldAccessToken = localStorage.getItem('auth_token');
        const oldRefreshToken = localStorage.getItem('refresh_token');

        if (oldRefreshToken && !sessionStorage.getItem('refresh_token')) {
            // Migrate refresh token to sessionStorage
            sessionStorage.setItem('refresh_token', oldRefreshToken);
            console.log('[Migration] Moved refresh token from localStorage to sessionStorage');
        }

        // Clean up old tokens from localStorage
        if (oldAccessToken) {
            localStorage.removeItem('auth_token');
            console.log('[Migration] Removed access token from localStorage (now stored in memory)');
        }
        if (oldRefreshToken) {
            localStorage.removeItem('refresh_token');
            console.log('[Migration] Cleaned up old refresh token from localStorage');
        }
    } catch (error) {
        console.error('[Migration] Error migrating token storage:', error);
    }
}
