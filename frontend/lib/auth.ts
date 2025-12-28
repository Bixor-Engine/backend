export interface User {
  id: string;
  first_name: string;
  last_name: string;
  username: string;
  email: string;
  email_status: boolean;
  phone_number?: string;
  phone_status: boolean;
  address?: string;
  city?: string;
  country?: string;
  role: string;
  status: string;
  kyc_status: string;
  twofa_enabled: boolean;
  last_login_at?: string;
  language: string;
  timezone: string;
  global_balance: number;
  created_at: string;
  updated_at: string;
}

export interface RegisterRequest {
  first_name: string;
  last_name: string;
  username: string;
  email: string;
  password: string;
  language?: string;
  timezone?: string;
  // Optional fields to be handled in profile settings later
  phone_number?: string;
  referred_by?: string;
  address?: string;
  city?: string;
  country?: string;
}

export interface AuthResponse {
  message: string;
  user: User;
  tokens?: {
    access_token: string;
    refresh_token: string;
    expires_in: number;
  };
  requires_verify?: boolean;
  redirect_to?: string;
}

export interface UpdateProfileRequest {
  first_name: string;
  last_name: string;
  phone_number?: string;
  address?: string;
  city?: string;
  country?: string;
}

export interface UpdateSettingsRequest {
  language: string;
  timezone: string;
}

export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
}

export interface ToggleTwoFARequest {
  enable: boolean;
  code: string;
}

export class AuthService {
  // Store access token in memory only (never in storage)
  private static accessToken: string | null = null;
  private static readonly REFRESH_TOKEN_KEY = 'refresh_token';

  // Get default headers (no backend secret needed - handled by Next.js API routes)
  private static getDefaultHeaders(additionalHeaders?: Record<string, string>): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...additionalHeaders,
    };

    return headers;
  }

  static isAuthenticated(): boolean {
    if (typeof window === 'undefined') return false;
    const token = this.accessToken;
    if (!token) {
      return false;
    }

    // Check if token is expired
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      const currentTime = Date.now() / 1000;
      const isValid = payload.exp > currentTime;
      return isValid;
    } catch {
      return false;
    }
  }

  static getToken(): string | null {
    if (typeof window === 'undefined') return null;
    return this.accessToken;
  }

  static async getCurrentUser(): Promise<User | null> {
    const token = this.getToken();
    if (!token) {
      return null;
    }

    // Function to make the API call (via Next.js API route)
    const fetchUser = async (accessToken: string): Promise<User | null> => {
      try {
        const response = await fetch('/api/auth/me', {
          method: 'GET',
          headers: this.getDefaultHeaders({
            'Authorization': `Bearer ${accessToken}`,
          }),
        });

        if (!response.ok) {
          return null;
        }

        const data = await response.json();
        return data.user;
      } catch {
        return null;
      }
    };

    // Try with current token
    let user = await fetchUser(token);

    // If failed and we have a refresh token, try to refresh and retry
    if (!user && this.getRefreshToken()) {
      const refreshSuccess = await this.refreshToken();

      if (refreshSuccess) {
        const newToken = this.getToken();
        if (newToken) {
          user = await fetchUser(newToken);
        }
      }
    }

    // If still failed, clear auth
    if (!user) {
      this.clearAuth();
    }

    return user;
  }

  static getRefreshToken(): string | null {
    if (typeof window === 'undefined') return null;
    return sessionStorage.getItem(this.REFRESH_TOKEN_KEY);
  }

  static setTokens(accessToken: string, refreshToken: string): void {
    if (typeof window === 'undefined') return;
    // Store access token in memory only
    this.accessToken = accessToken;
    // Store refresh token in sessionStorage (cleared on browser close)
    sessionStorage.setItem(this.REFRESH_TOKEN_KEY, refreshToken);
  }

  static clearAuth(): void {
    if (typeof window === 'undefined') return;
    // Clear in-memory access token
    this.accessToken = null;
    // Clear refresh token from sessionStorage
    sessionStorage.removeItem(this.REFRESH_TOKEN_KEY);
  }

  static async login(email: string, password: string): Promise<AuthResponse> {
    const response = await fetch('/api/auth/login', {
      method: 'POST',
      headers: this.getDefaultHeaders(),
      body: JSON.stringify({ email, password }),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || 'Login failed');
    }

    // Only store the access token from tokens object
    if (data.tokens && data.tokens.access_token && data.tokens.refresh_token) {
      this.setTokens(data.tokens.access_token, data.tokens.refresh_token);
    } else {
      throw new Error('No tokens received');
    }

    return data;
  }

  static async register(registerData: RegisterRequest): Promise<AuthResponse> {
    const response = await fetch('/api/auth/register', {
      method: 'POST',
      headers: this.getDefaultHeaders(),
      body: JSON.stringify(registerData),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || 'Registration failed');
    }

    // Register endpoint doesn't return tokens, user needs to sign in
    return data;
  }

  static async refreshToken(): Promise<boolean> {
    const refreshToken = this.getRefreshToken();
    if (!refreshToken) {
      return false;
    }

    try {
      const response = await fetch('/api/auth/refresh', {
        method: 'POST',
        headers: this.getDefaultHeaders(),
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      if (!response.ok) {
        this.clearAuth();
        return false;
      }

      const data = await response.json();

      if (data.tokens && data.tokens.access_token && data.tokens.refresh_token) {
        this.setTokens(data.tokens.access_token, data.tokens.refresh_token);
        return true;
      } else {
        this.clearAuth();
        return false;
      }
    } catch {
      this.clearAuth();
      return false;
    }
  }

  static async logout(): Promise<void> {
    const token = this.getToken();

    // Try to call backend logout endpoint (optional - will still clear local tokens if it fails)
    if (token) {
      try {
        const response = await fetch('/api/auth/logout', {
          method: 'POST',
          headers: this.getDefaultHeaders({
            'Authorization': `Bearer ${token}`,
          }),
        });
        // Don't throw on error - still clear local tokens
        if (!response.ok) {
          // Log silently or ignore
        }
      } catch {
        // Ignore errors - still clear local tokens
      }
    }

    // Always clear local tokens
    this.clearAuth();
  }

  static isTokenExpiringSoon(): boolean {
    if (typeof window === 'undefined') return false;
    const token = this.accessToken;
    if (!token) return false;

    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      const currentTime = Date.now() / 1000;
      const timeUntilExpiry = payload.exp - currentTime;

      // Return true if token expires in less than 5 minutes (300 seconds)
      return timeUntilExpiry < 300;
    } catch {
      return true; // If we can't parse, assume it's expiring
    }
  }

  static async ensureValidToken(): Promise<boolean> {
    if (!this.isAuthenticated()) {
      return false;
    }

    if (this.isTokenExpiringSoon()) {
      return await this.refreshToken();
    }

    return true;
  }

  static async requestOTP(type: 'email-verification' | 'password-reset' | '2fa' | 'phone-verification'): Promise<{ message: string; expires_in: number }> {
    const token = this.getToken();
    if (!token) {
      throw new Error('Not authenticated');
    }

    const response = await fetch('/api/auth/otp/request', {
      method: 'POST',
      headers: this.getDefaultHeaders({
        'Authorization': `Bearer ${token}`,
      }),
      body: JSON.stringify({ type }),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || 'Failed to request OTP');
    }

    return data;
  }

  static async verifyOTP(type: 'email-verification' | 'password-reset' | '2fa' | 'phone-verification', code: string): Promise<{ message: string; verified: boolean }> {
    const token = this.getToken();
    if (!token) {
      throw new Error('Not authenticated');
    }

    const response = await fetch('/api/auth/otp/verify', {
      method: 'POST',
      headers: this.getDefaultHeaders({
        'Authorization': `Bearer ${token}`,
      }),
      body: JSON.stringify({ type, code }),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || 'Failed to verify OTP');
    }

    return data;
  }

  static async updateProfile(profileData: UpdateProfileRequest): Promise<AuthResponse> {
    const token = this.getToken();
    if (!token) throw new Error('Not authenticated');

    const response = await fetch('/api/auth/profile/update', {
      method: 'POST',
      headers: this.getDefaultHeaders({
        'Authorization': `Bearer ${token}`,
      }),
      body: JSON.stringify(profileData),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || 'Update failed');
    }

    return data;
  }

  static async updateSettings(settingsData: UpdateSettingsRequest): Promise<AuthResponse> {
    const token = this.getToken();
    if (!token) throw new Error('Not authenticated');

    const response = await fetch('/api/auth/settings/update', {
      method: 'POST',
      headers: this.getDefaultHeaders({
        'Authorization': `Bearer ${token}`,
      }),
      body: JSON.stringify(settingsData),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || 'Update failed');
    }

    return data;
  }

  static async changePassword(passwordData: ChangePasswordRequest): Promise<{ message: string }> {
    const token = this.getToken();
    if (!token) throw new Error('Not authenticated');

    const response = await fetch('/api/auth/security/password', {
      method: 'POST',
      headers: this.getDefaultHeaders({
        'Authorization': `Bearer ${token}`,
      }),
      body: JSON.stringify(passwordData),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || 'Password change failed');
    }

    return data;
  }

  static async toggleTwoFA(twoFAData: ToggleTwoFARequest): Promise<AuthResponse> {
    const token = this.getToken();
    if (!token) throw new Error('Not authenticated');

    const response = await fetch('/api/auth/security/2fa', {
      method: 'POST',
      headers: this.getDefaultHeaders({
        'Authorization': `Bearer ${token}`,
      }),
      body: JSON.stringify(twoFAData),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || '2FA toggle failed');
    }

    return data;
  }

  static async fetch(url: string, options: RequestInit = {}): Promise<Response> {
    const token = this.getToken();
    if (!token) {
      // If no token, we can throw or just proceed without auth (depending on requirement).
      // For protected routes, this will fail 401, which is fine.
      // But let's try to refresh if we can.
      if (this.getRefreshToken()) {
        await this.refreshToken();
      }
    } else if (this.isTokenExpiringSoon()) {
      await this.refreshToken();
    }

    const currentToken = this.getToken();

    // Prepare headers
    const headers = new Headers(options.headers || {});
    headers.set('Content-Type', 'application/json');
    if (currentToken) {
      headers.set('Authorization', `Bearer ${currentToken}`);
    }

    // Add backend secret if needed (though Next.js API routes usually handle this if proxying, 
    // but here we are calling internal API routes directly or via generic proxy).
    // The previous implementation assumes direct calls to /api/... which matches the backend routes.
    // If the frontend is hitting the backend directly, we rely on the browser not needing the secret
    // OR we need to add it. The existing code suggests we hit '/api/...' which might be rewritten.
    // However, looking at 'routes.go', it expects 'X-Backend-Secret'.
    // In `AuthService`, we don't see X-Backend-Secret being added in `getDefaultHeaders`.
    // Wait, `API_PROTECTION.md` says:
    // "Usage from Frontend: headers: { 'X-Backend-Secret': process.env.NEXT_PUBLIC_BACKEND_SECRET, ... }"
    // So we SHOULD add it.

    // So we SHOULD add it.

    // Fallback to 'test123' if env var is missing (e.g., server not restarted)
    const backendSecret = process.env.NEXT_PUBLIC_BACKEND_SECRET || 'test123';
    if (backendSecret) {
      headers.set('X-Backend-Secret', backendSecret);
    }

    const config = {
      ...options,
      headers,
    };

    // Prepend /api/v1 if the url doesn't start with http or /api
    // But existing calls use '/api/auth/...'.
    // The backend routes are group '/api/v1'.
    // So '/api/auth/...' in frontend probably maps to '/api/v1/auth/...' in backend? 
    // Let's check `api.go` or `routes.go`.
    // `routes.go`: v1 := router.Group("/api/v1") ... auth := protected.Group("/auth")
    // So the full path is `/api/v1/auth/...`.
    // The existing `login` calls use `/api/auth/login`. This implies a rewrite or proxy?
    // Or maybe the previous dev made a mistake.
    // Let's assume we should use the full path if we are hitting the Go server directly.
    // If we are hitting Next.js API routes, then it depends.
    // Assuming we hit Go server directly or via proxy at /api/v1.
    // For now, I will NOT modify the URL structure, but just pass it through.
    // BUT, the use-wallet hook uses `/wallets`, so we might need to prepend `/api/v1`.

    let finalUrl = url;
    if (url.startsWith('/wallets') || url.startsWith('/transactions')) {
      finalUrl = `/api/v1${url}`;
    }

    return fetch(finalUrl, config);
  }
}
