'use client';

import { useState } from 'react';
import { AuthService } from '@/lib/auth';

export default function DebugAuth() {
  const [result, setResult] = useState<unknown>(null);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleLogin = async () => {
    try {
      console.log('=== DEBUG LOGIN START ===');
      
      // Clear any existing auth first
      AuthService.clearAuth();
      console.log('Cleared existing auth');
      
      // Check what's in localStorage before login
      console.log('localStorage before login:', {
        auth_token: localStorage.getItem('auth_token'),
        refresh_token: localStorage.getItem('refresh_token')
      });
      
      const response = await AuthService.login(email, password);
      console.log('Login response:', response);
      
      // Check what's in localStorage after login
      console.log('localStorage after login:', {
        auth_token: localStorage.getItem('auth_token'),
        refresh_token: localStorage.getItem('refresh_token')
      });
      
      console.log('=== DEBUG LOGIN END ===');
      setResult(response);
    } catch (error) {
      console.error('Login error:', error);
      setResult({ error: error instanceof Error ? error.message : String(error) });
    }
  };

  const checkLocalStorage = () => {
    const storage = {
      auth_token: localStorage.getItem('auth_token'),
      refresh_token: localStorage.getItem('refresh_token')
    };
    console.log('Current localStorage:', storage);
    setResult(storage);
  };

  const clearStorage = () => {
    AuthService.clearAuth();
    console.log('Cleared localStorage');
    checkLocalStorage();
  };

  return (
    <div className="container mx-auto p-8">
      <h1 className="text-2xl font-bold mb-6">Debug Authentication</h1>
      
      <div className="space-y-4 max-w-md">
        <div>
          <label className="block text-sm font-medium mb-1">Email:</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md"
            placeholder="Enter email"
          />
        </div>
        
        <div>
          <label className="block text-sm font-medium mb-1">Password:</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md"
            placeholder="Enter password"
          />
        </div>
        
        <div className="space-x-2">
          <button
            onClick={handleLogin}
            className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600"
          >
            Test Login
          </button>
          <button
            onClick={checkLocalStorage}
            className="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600"
          >
            Check Storage
          </button>
          <button
            onClick={clearStorage}
            className="px-4 py-2 bg-red-500 text-white rounded-md hover:bg-red-600"
          >
            Clear Storage
          </button>
        </div>
      </div>
      
      {result && (
        <div className="mt-6">
          <h2 className="text-lg font-semibold mb-2">Result:</h2>
          <pre className="bg-gray-100 p-4 rounded-md overflow-auto text-sm">
            {JSON.stringify(result, null, 2)}
          </pre>
        </div>
      )}
      
      <div className="mt-6">
        <h2 className="text-lg font-semibold mb-2">Instructions:</h2>
        <ol className="list-decimal list-inside space-y-1 text-sm">
          <li>Enter your email and password</li>
          <li>Click "Test Login" and check the browser console</li>
          <li>Click "Check Storage" to see what's in localStorage</li>
          <li>Open browser DevTools (F12) → Application → Local Storage</li>
          <li>Look for both 'auth_token' and 'refresh_token' keys</li>
        </ol>
      </div>
    </div>
  );
}
