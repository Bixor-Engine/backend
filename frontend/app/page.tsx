'use client';

import Link from "next/link";
import { useState, useEffect } from "react";
import Navbar from "./components/Navbar";
import { AuthService, User } from "./lib/auth";

export default function Home() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    const checkAuth = () => {
      setIsAuthenticated(AuthService.isAuthenticated());
      setUser(AuthService.getUser());
    };

    checkAuth();
    
    // Listen for storage changes (when user logs in/out in another tab)
    const handleStorageChange = () => {
      checkAuth();
    };

    window.addEventListener('storage', handleStorageChange);
    return () => window.removeEventListener('storage', handleStorageChange);
  }, []);

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50">
      {/* Navigation */}
      <Navbar />

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-900 sm:text-6xl">
            Welcome to{" "}
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-purple-600">
              Bixor Engine
            </span>
          </h1>
          <p className="mt-6 text-lg leading-8 text-gray-600 max-w-2xl mx-auto">
            A powerful and modern web application platform built with Go backend and Next.js frontend.
            Experience seamless authentication, robust APIs, and beautiful user interfaces.
          </p>
          
          {!isAuthenticated && (
            <div className="mt-10 flex items-center justify-center gap-x-6">
              <Link
                href="/auth/signup"
                className="rounded-md bg-blue-600 px-6 py-3 text-lg font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 transition-colors"
              >
                Get started
              </Link>
              <Link
                href="/auth/signin"
                className="text-lg font-semibold leading-6 text-gray-900 hover:text-blue-600 transition-colors"
              >
                Sign in <span aria-hidden="true">→</span>
              </Link>
            </div>
          )}

          {isAuthenticated && (
            <div className="mt-10">
              <div className="bg-white rounded-lg shadow-md p-8 max-w-md mx-auto">
                <h2 className="text-2xl font-bold text-gray-900 mb-4">Dashboard</h2>
                <p className="text-gray-600 mb-6">Welcome back! You are successfully authenticated.</p>
                <div className="space-y-3">
                  <div className="text-left">
                    <span className="text-sm font-medium text-gray-500">Name:</span>
                    <span className="ml-2 text-sm text-gray-900">{user?.first_name} {user?.last_name}</span>
                  </div>
                  <div className="text-left">
                    <span className="text-sm font-medium text-gray-500">Username:</span>
                    <span className="ml-2 text-sm text-gray-900">{user?.username || 'N/A'}</span>
                  </div>
                  <div className="text-left">
                    <span className="text-sm font-medium text-gray-500">Email:</span>
                    <span className="ml-2 text-sm text-gray-900">{user?.email || 'N/A'}</span>
                  </div>
                  <div className="text-left">
                    <span className="text-sm font-medium text-gray-500">Status:</span>
                    <span className="ml-2 text-sm text-gray-900 capitalize">{user?.status || 'N/A'}</span>
                  </div>
                  <div className="text-left">
                    <span className="text-sm font-medium text-gray-500">Role:</span>
                    <span className="ml-2 text-sm text-gray-900 capitalize">{user?.role || 'N/A'}</span>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Features Section */}
        <div className="mt-20">
          <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3">
            <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
              <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mb-4">
                <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Secure Authentication</h3>
              <p className="text-gray-600">Robust JWT-based authentication with secure password hashing using Argon2.</p>
            </div>

            <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
              <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center mb-4">
                <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">High Performance</h3>
              <p className="text-gray-600">Built with Go backend for exceptional performance and Next.js for optimal user experience.</p>
            </div>

            <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
              <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center mb-4">
                <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Modern Stack</h3>
              <p className="text-gray-600">Cutting-edge technology stack with TypeScript, Tailwind CSS, and RESTful APIs.</p>
            </div>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="bg-white border-t mt-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="text-center text-gray-500">
            <p>&copy; 2025 Bixor Engine. Built with ❤️ using Go and Next.js.</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
