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

export class AuthService {
  private static readonly TOKEN_KEY = 'auth_token';
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
    const token = localStorage.getItem(this.TOKEN_KEY);
    if (!token) {
      return false;
    }
    
    // Check if token is expired
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      const currentTime = Date.now() / 1000;
      const isValid = payload.exp > currentTime;
      return isValid;
    } catch (error) {
      return false;
    }
  }

  static getToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(this.TOKEN_KEY);
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
      } catch (error) {
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
    return localStorage.getItem(this.REFRESH_TOKEN_KEY);
  }

  static setTokens(accessToken: string, refreshToken: string): void {
    if (typeof window === 'undefined') return;
    localStorage.setItem(this.TOKEN_KEY, accessToken);
    localStorage.setItem(this.REFRESH_TOKEN_KEY, refreshToken);
  }

  static clearAuth(): void {
    if (typeof window === 'undefined') return;
    localStorage.removeItem(this.TOKEN_KEY);
    localStorage.removeItem(this.REFRESH_TOKEN_KEY);
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
    } catch (error) {
      this.clearAuth();
      return false;
    }
  }

  static logout(): void {
    this.clearAuth();
  }

  static isTokenExpiringSoon(): boolean {
    if (typeof window === 'undefined') return false;
    const token = localStorage.getItem(this.TOKEN_KEY);
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
}
