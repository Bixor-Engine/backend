import { NextRequest, NextResponse } from 'next/server';

export async function GET(request: NextRequest) {
  try {
    const authHeader = request.headers.get('Authorization');
    if (!authHeader) {
      return NextResponse.json(
        { error: 'unauthorized', message: 'Authorization header required' },
        { status: 401 }
      );
    }

    const backendUrl = process.env.BACKEND_URL || 'http://localhost:8080';
    const backendSecret = process.env.BACKEND_SECRET;

    if (!backendSecret) {
      return NextResponse.json(
        { error: 'backend_secret_not_configured', message: 'Backend secret is not configured' },
        { status: 500 }
      );
    }

    const response = await fetch(`${backendUrl}/api/v1/auth/me`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': authHeader,
        'X-Backend-Secret': backendSecret,
      },
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

