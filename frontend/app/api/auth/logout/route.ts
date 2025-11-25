import { NextRequest, NextResponse } from 'next/server';

export async function POST(request: NextRequest) {
  try {
    const authHeader = request.headers.get('Authorization');
    const backendUrl = process.env.BACKEND_URL || 'http://localhost:8080';
    const backendSecret = process.env.BACKEND_SECRET;

    if (!backendSecret) {
      return NextResponse.json(
        { error: 'backend_secret_not_configured', message: 'Backend secret is not configured' },
        { status: 500 }
      );
    }

    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      'X-Backend-Secret': backendSecret,
    };

    // Add authorization header if provided
    if (authHeader) {
      headers['Authorization'] = authHeader;
    }

    const response = await fetch(`${backendUrl}/api/v1/auth/logout`, {
      method: 'POST',
      headers,
    });

    const data = await response.json();
    return NextResponse.json(data, { status: response.status });
  } catch (error) {
    return NextResponse.json(
      { error: 'internal_error', message: 'Failed to process request' },
      { status: 500 }
    );
  }
}

