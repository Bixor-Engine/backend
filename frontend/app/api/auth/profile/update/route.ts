import { NextRequest, NextResponse } from 'next/server';

export async function POST(request: NextRequest) {
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

        const body = await request.json();

        const response = await fetch(`${backendUrl}/api/v1/auth/profile/update`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': authHeader,
                'X-Backend-Secret': backendSecret,
            },
            body: JSON.stringify(body),
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
