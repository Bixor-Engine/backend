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
}

export class AuthService {
  private static readonly TOKEN_KEY = 'auth_token';
  private static readonly REFRESH_TOKEN_KEY = 'refresh_token';

  static isAuthenticated(): boolean {
    if (typeof window === 'undefined') return false;
    const token = localStorage.getItem(this.TOKEN_KEY);
    if (!token) {
      console.log('AuthService: No token found');
      return false;
    }
    
    // Check if token is expired
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      const currentTime = Date.now() / 1000;
      const isValid = payload.exp > currentTime;
      console.log('AuthService: Token validation:', { 
        exp: payload.exp, 
        current: currentTime, 
        isValid,
        timeLeft: payload.exp - currentTime 
      });
      return isValid;
    } catch (error) {
      console.log('AuthService: Token parsing failed:', error);
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
      console.log('AuthService: No token for getCurrentUser');
      return null;
    }

    // Function to make the API call
    const fetchUser = async (accessToken: string): Promise<User | null> => {
      try {
        const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
        console.log('AuthService: Fetching user from', `${apiUrl}/api/v1/auth/me`);
        const response = await fetch(`${apiUrl}/api/v1/auth/me`, {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${accessToken}`,
            'Content-Type': 'application/json',
          },
        });

        if (!response.ok) {
          console.log('AuthService: API responded with error:', response.status, response.statusText);
          return null;
        }

        const data = await response.json();
        console.log('AuthService: User data received:', data);
        return data.user;
      } catch (error) {
        console.log('AuthService: getCurrentUser fetch failed:', error);
        return null;
      }
    };

    // Try with current token
    let user = await fetchUser(token);
    
    // If failed and we have a refresh token, try to refresh and retry
    if (!user && this.getRefreshToken()) {
      console.log('AuthService: Token might be expired, trying to refresh');
      const refreshSuccess = await this.refreshToken();
      
      if (refreshSuccess) {
        const newToken = this.getToken();
        if (newToken) {
          console.log('AuthService: Retrying with refreshed token');
          user = await fetchUser(newToken);
        }
      }
    }

    // If still failed, clear auth
    if (!user) {
      console.log('AuthService: Failed to get user data, clearing auth');
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
    console.log('AuthService: setTokens called with:', {
      accessToken: accessToken ? 'present' : 'missing',
      refreshToken: refreshToken ? 'present' : 'missing'
    });
    localStorage.setItem(this.TOKEN_KEY, accessToken);
    localStorage.setItem(this.REFRESH_TOKEN_KEY, refreshToken);
    console.log('AuthService: Tokens stored. Verification:', {
      auth_token: localStorage.getItem(this.TOKEN_KEY) ? 'stored' : 'not stored',
      refresh_token: localStorage.getItem(this.REFRESH_TOKEN_KEY) ? 'stored' : 'not stored'
    });
  }

  static clearAuth(): void {
    if (typeof window === 'undefined') return;
    localStorage.removeItem(this.TOKEN_KEY);
    localStorage.removeItem(this.REFRESH_TOKEN_KEY);
  }

  static async login(email: string, password: string): Promise<AuthResponse> {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    console.log('AuthService: Attempting login to', `${apiUrl}/api/v1/auth/login`);
    
    const response = await fetch(`${apiUrl}/api/v1/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    const data = await response.json();
    console.log('AuthService: Login response:', data);

    if (!response.ok) {
      throw new Error(data.message || 'Login failed');
    }

    // Only store the access token from tokens object
    if (data.tokens && data.tokens.access_token && data.tokens.refresh_token) {
      console.log('AuthService: Storing access and refresh tokens');
      this.setTokens(data.tokens.access_token, data.tokens.refresh_token);
    } else {
      console.log('AuthService: No tokens in response:', data);
      throw new Error('No tokens received');
    }

    return data;
  }

  static async register(registerData: RegisterRequest): Promise<AuthResponse> {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    const response = await fetch(`${apiUrl}/api/v1/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
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
      console.log('AuthService: No refresh token available');
      return false;
    }

    try {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
      console.log('AuthService: Attempting to refresh token');
      
      const response = await fetch(`${apiUrl}/api/v1/auth/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });

      if (!response.ok) {
        console.log('AuthService: Token refresh failed:', response.status);
        this.clearAuth();
        return false;
      }

      const data = await response.json();
      console.log('AuthService: Token refresh response:', data);

      if (data.tokens && data.tokens.access_token && data.tokens.refresh_token) {
        console.log('AuthService: Storing new tokens after refresh');
        this.setTokens(data.tokens.access_token, data.tokens.refresh_token);
        return true;
      } else {
        console.log('AuthService: No tokens in refresh response');
        this.clearAuth();
        return false;
      }
    } catch (error) {
      console.log('AuthService: Token refresh error:', error);
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
      console.log('AuthService: Token expiring soon, refreshing...');
      return await this.refreshToken();
    }

    return true;
  }
}
