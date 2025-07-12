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
  token: string;
  user: User;
  message?: string;
}

export class AuthService {
  private static readonly TOKEN_KEY = 'auth_token';
  private static readonly USER_KEY = 'user_data';

  static isAuthenticated(): boolean {
    if (typeof window === 'undefined') return false;
    return !!localStorage.getItem(this.TOKEN_KEY);
  }

  static getToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(this.TOKEN_KEY);
  }

  static getUser(): User | null {
    if (typeof window === 'undefined') return null;
    const userData = localStorage.getItem(this.USER_KEY);
    return userData ? JSON.parse(userData) : null;
  }

  static setAuth(token: string, user: User): void {
    if (typeof window === 'undefined') return;
    localStorage.setItem(this.TOKEN_KEY, token);
    localStorage.setItem(this.USER_KEY, JSON.stringify(user));
  }

  static clearAuth(): void {
    if (typeof window === 'undefined') return;
    localStorage.removeItem(this.TOKEN_KEY);
    localStorage.removeItem(this.USER_KEY);
  }

  static async login(email: string, password: string): Promise<AuthResponse> {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    const response = await fetch(`${apiUrl}/api/v1/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    const data = await response.json();

    if (!response.ok) {
      throw new Error(data.message || 'Login failed');
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

    return data;
  }

  static async refreshToken(): Promise<string | null> {
    const token = this.getToken();
    if (!token) return null;

    try {
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/api/v1/auth/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        this.clearAuth();
        return null;
      }

      const data = await response.json();
      if (data.token) {
        localStorage.setItem(this.TOKEN_KEY, data.token);
        return data.token;
      }

      return null;
    } catch {
      this.clearAuth();
      return null;
    }
  }

  static logout(): void {
    this.clearAuth();
  }
}
