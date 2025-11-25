# API Protection System

This document describes the three-tier API protection system implemented in the Bixor Trading Engine.

## Overview

The API uses three levels of protection:

1. **Public API** - No authentication required
2. **Backend Secret Protected** - Requires backend secret (for frontend requests)
3. **Personal API** - Requires user token (for future personal API access)

## Protection Levels

### 1. Public API Routes

These routes can be accessed without any authentication:

- `GET /` - Landing page
- `GET /api/v1/health` - Health check
- `GET /api/v1/status` - Service status
- `GET /api/v1/info` - API information

**Usage:** Direct access, no headers required.

### 2. Backend Secret Protected Routes

These routes require the `X-Backend-Secret` header (or `X-API-Secret` as alternative):

- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh JWT tokens
- `GET /api/v1/auth/me` - Get current user (also requires JWT)
- `POST /api/v1/auth/otp/request` - Request OTP (also requires JWT)
- `POST /api/v1/auth/otp/verify` - Verify OTP (also requires JWT)

**Usage from Frontend:**
```typescript
headers: {
  'X-Backend-Secret': process.env.NEXT_PUBLIC_BACKEND_SECRET,
  'Authorization': 'Bearer <jwt_token>', // For authenticated endpoints
}
```

**Usage from External:**
```bash
curl -H "X-Backend-Secret: your-secret-here" \
     -H "Authorization: Bearer <jwt_token>" \
     https://api.example.com/api/v1/auth/me
```

### 3. Personal API Routes (Future)

These routes will require user-specific API tokens (not JWT):

- `GET /api/v1/personal/*` - Personal API endpoints (to be implemented)

**Usage:** Will use custom user API tokens generated in user settings.

## Configuration

### Backend (.env)

```env
BACKEND_SECRET=your-super-secret-key-here-change-in-production
```

**Important:** 
- Generate a strong random secret for production
- Never commit the actual secret to version control
- If `BACKEND_SECRET` is not set, the middleware allows all requests (development mode)

### Frontend (.env.local)

```env
NEXT_PUBLIC_BACKEND_SECRET=your-super-secret-key-here-change-in-production
```

**Important:**
- Must match the backend `BACKEND_SECRET`
- The `NEXT_PUBLIC_` prefix makes it available in the browser
- This is safe because it's only used to authenticate frontend requests

## Middleware

### BackendSecretMiddleware

Located in `internal/middleware/auth.go`:

- Checks for `X-Backend-Secret` or `X-API-Secret` header
- Compares against `BACKEND_SECRET` environment variable
- Returns 401 if secret is missing or invalid
- Allows requests if `BACKEND_SECRET` is not set (development mode)

### PublicMiddleware

- Passthrough middleware for public routes
- No validation performed

### UserTokenMiddleware

- Placeholder for future personal API token validation
- Currently a passthrough

## Security Considerations

1. **Backend Secret:**
   - Use a strong, randomly generated secret (minimum 32 characters)
   - Rotate secrets periodically
   - Use different secrets for development and production

2. **CORS:**
   - Backend secret headers are allowed in CORS configuration
   - Adjust CORS settings based on your deployment

3. **Development Mode:**
   - If `BACKEND_SECRET` is not set, all requests are allowed
   - This is convenient for development but should never be used in production

## Route Organization

Routes are organized in `internal/routes/routes.go`:

```go
// Public routes
public := v1.Group("")
public.Use(middleware.PublicMiddleware())

// Backend secret protected routes
protected := v1.Group("")
protected.Use(middleware.BackendSecretMiddleware())

// Personal API routes (future)
personal := v1.Group("/personal")
personal.Use(middleware.UserTokenMiddleware())
```

## Frontend Integration

The `AuthService` class automatically includes the backend secret in all requests:

```typescript
// Automatically adds X-Backend-Secret header
const response = await AuthService.login(email, password);
```

The secret is read from `process.env.NEXT_PUBLIC_BACKEND_SECRET` and included in all API requests.

## Testing

### Test Public Endpoint
```bash
curl http://localhost:8080/api/v1/health
```

### Test Protected Endpoint (without secret - should fail)
```bash
curl http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test"}'
```

### Test Protected Endpoint (with secret - should work)
```bash
curl http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "X-Backend-Secret: your-secret-here" \
  -d '{"email":"test@example.com","password":"test"}'
```

## Future Enhancements

1. **Personal API Tokens:**
   - Users will be able to generate personal API tokens
   - These tokens will be stored in the database
   - Tokens will have scopes/permissions
   - Tokens can be revoked by users

2. **Rate Limiting:**
   - Implement rate limiting per secret/token
   - Different limits for different protection levels

3. **IP Whitelisting:**
   - Optional IP whitelisting for backend secret
   - Additional security layer

