# Bixor Engine Frontend

A modern Next.js frontend application for the Bixor Engine platform, featuring secure authentication, responsive design, and seamless integration with the Go backend.

## Features

- 🔐 **Secure Authentication** - JWT-based auth with proper token management
- 🎨 **Modern UI** - Beautiful, responsive design with Tailwind CSS
- ⚡ **High Performance** - Built with Next.js 15 and React 19
- 🛡️ **Type Safety** - Full TypeScript support
- 🔄 **Real-time Updates** - Automatic auth state synchronization
- 📱 **Mobile Responsive** - Works perfectly on all devices

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn
- Bixor Engine backend running on port 8080

### Installation

1. Install dependencies:
```bash
npm install
```

2. Configure environment variables:
```bash
cp .env.local.example .env.local
```

Edit `.env.local` with your configuration:
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_NAME=Bixor Engine
NEXT_PUBLIC_APP_VERSION=1.0.0
```

3. Start the development server:
```bash
npm run dev
```

4. Open [http://localhost:3000](http://localhost:3000) in your browser.

## Project Structure

```
frontend/
├── app/
│   ├── components/          # Reusable UI components
│   │   └── Navbar.tsx      # Navigation component
│   ├── lib/                # Utility libraries
│   │   └── auth.ts         # Authentication service
│   ├── auth/               # Authentication pages
│   │   ├── signin/         # Sign in page
│   │   │   └── page.tsx
│   │   └── signup/         # Sign up page
│   │       └── page.tsx
│   ├── globals.css         # Global styles
│   ├── layout.tsx          # Root layout
│   └── page.tsx            # Home page
├── public/                 # Static assets
├── .env.local              # Environment variables
└── package.json
```

## Authentication Flow

### Sign Up
1. User fills registration form (`/auth/signup`)
2. Frontend calls `POST /api/v1/auth/register`
3. On success, redirects to sign in page
4. On error, displays error message

### Sign In
1. User fills login form (`/auth/signin`)
2. Frontend calls `POST /api/v1/auth/login`
3. On success:
   - Stores JWT token in localStorage
   - Stores user data in localStorage
   - Redirects to home page
4. On error, displays error message

### Authentication State
- `AuthService` class manages all auth operations
- Automatic token storage and retrieval
- Cross-tab synchronization via localStorage events
- Automatic cleanup on logout

## API Integration

The frontend integrates with the Bixor Engine backend through RESTful APIs:

### Authentication Endpoints
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Token refresh (if implemented)

### Expected Response Format

**Registration/Login Success:**
```json
{
  "token": "jwt_token_here",
  "user": {
    "id": "user_id",
    "username": "username",
    "email": "user@example.com",
    "createdAt": "2025-01-01T00:00:00Z",
    "updatedAt": "2025-01-01T00:00:00Z"
  },
  "message": "Success message"
}
```

**Error Response:**
```json
{
  "error": "error_code",
  "message": "Human readable error message",
  "details": "Additional error details"
}
```

## Components

### Navbar Component
- Displays app branding
- Shows authentication status
- Login/logout functionality
- Responsive navigation menu

### AuthService
- Centralized authentication management
- Token storage and retrieval
- API communication
- Type-safe user data handling

## Styling

The application uses Tailwind CSS for styling with:
- Responsive design principles
- Modern color palette (blue/purple gradients)
- Consistent spacing and typography
- Hover and focus states
- Loading animations

## Development

### Available Scripts

```bash
# Development server
npm run dev

# Production build
npm run build

# Start production server
npm start

# Lint code
npm run lint
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API URL | `http://localhost:8080` |
| `NEXT_PUBLIC_APP_NAME` | Application name | `Bixor Engine` |
| `NEXT_PUBLIC_APP_VERSION` | App version | `1.0.0` |

## Deployment

### Production Build

1. Build the application:
```bash
npm run build
```

2. Start the production server:
```bash
npm start
```

### Docker Deployment

```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build
EXPOSE 3000
CMD ["npm", "start"]
```

## Security Considerations

- JWT tokens stored in localStorage (consider httpOnly cookies for production)
- CSRF protection through SameSite cookies
- XSS prevention through React's built-in escaping
- Input validation on both client and server
- Secure password requirements

## Browser Support

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Contributing

1. Follow the existing code style
2. Use TypeScript for all new components
3. Add proper error handling
4. Test authentication flows
5. Ensure responsive design

## License

This project is part of the Bixor Engine platform.
