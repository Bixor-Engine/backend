# Development Tools

This directory contains interactive CLI tools for testing and development of the Bixor Trading Engine.

## Tools Overview

| Tool | Purpose | Location |
|------|---------|----------|
| [Authentication Tool](#authentication-tool) | Test complete auth flow, manage users | `tools/auth/` |
| [Password Hash Tool](#password-hash-tool) | Generate/verify Argon2i hashes | `tools/hash/` |

## Authentication Tool

Interactive CLI tool for testing the complete authentication system without needing a frontend.

### Usage
```bash
go run tools/auth/main.go
```

### Features

| Option | Description |
|--------|-------------|
| **1. Register account** | Create new user accounts with validation |
| **2. Login** | Test login with JWT token generation |
| **3. Verify user (activate account)** | Activate user accounts directly in database |
| **4. Refresh token** | Test JWT token refresh functionality |
| **5. Exit** | Close the tool |

### Example Flow

1. **Register a new user**:
   - Enter email, username, password
   - System validates input and creates account
   - User starts in "inactive" status

2. **Activate the account**:
   - Use option 3 to change status from "inactive" to "active"
   - Direct database update without manual SQL

3. **Test login**:
   - Login with activated account
   - Receive JWT access and refresh tokens
   - Tokens displayed for frontend integration testing

4. **Refresh tokens**:
   - Use refresh token to get new access token
   - Test token renewal flow

### API Endpoints Tested
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User authentication
- `POST /api/v1/auth/refresh` - Token refresh
- Direct database operations for user management

## Password Hash Tool

Interactive tool for testing Argon2i password hashing and verification.

### Usage
```bash
go run tools/hash/main.go
```

### Features

| Option | Description |
|--------|-------------|
| **1. Hash a password** | Generate Argon2i hash from plain text |
| **2. Verify a hash** | Test password against existing hash |
| **3. Exit** | Close the tool |

### Use Cases

#### Testing Password Security
```bash
# 1. Hash a password
Input: "MySecurePassword123!"
Output: $argon2i$v=19$m=65536,t=3,p=2$...

# 2. Verify the hash
Input password: "MySecurePassword123!"
Input hash: $argon2i$v=19$m=65536,t=3,p=2$...
Result: ✅ MATCH!
```

#### Debugging Authentication Issues
- Generate test hashes for database seeding
- Verify if password verification logic works correctly
- Test different password formats and special characters

### Security Parameters
- **Algorithm**: Argon2i (data-independent variant)
- **Memory**: 64MB (65536 KB)
- **Iterations**: 3
- **Parallelism**: 2 threads
- **Salt**: 16 random bytes, Base64 encoded
- **Hash**: 32 bytes, Base64 encoded

## Prerequisites

### Environment Setup
Both tools require a properly configured environment:

```bash
# 1. Database must be running and migrated
bixor database setup

# 2. Environment variables configured (.env file)
DATABASE_URL=postgres://username:password@localhost:5432/bixor?sslmode=disable
JWT_SECRET=your_secure_jwt_secret_key_here
JWT_EXPIRES_HOURS=24

# 3. API server running (for auth tool)
bixor server start
```

### Dependencies
- Go 1.21+
- PostgreSQL database
- Bixor API server running on localhost:8080

## Development Workflow

### Testing New Authentication Features
1. Use **hash tool** to understand password hashing behavior
2. Use **auth tool** to test complete registration/login flow
3. Verify JWT tokens work correctly with your frontend
4. Test edge cases (wrong passwords, inactive users, etc.)

### Debugging Issues
1. **Registration problems**: Use auth tool to test validation
2. **Login failures**: Use hash tool to verify password hashing
3. **JWT issues**: Use auth tool to inspect token content
4. **Database problems**: Use auth tool's user verification feature

## API Integration Testing

The authentication tool eliminates the need for frontend development during API testing:

```bash
# Instead of building a frontend form, just run:
go run tools/auth/main.go

# Test complete flows:
# Register → Activate → Login → Refresh
# All with real API calls and database operations
```

## Security Notes

⚠️ **Development Only**: These tools are for development and testing only.

🔒 **Password Safety**: The hash tool uses the same Argon2i implementation as the production API.

🔑 **JWT Tokens**: Generated tokens are real and can be used with the API during testing.

📊 **Database Access**: The auth tool performs real database operations - be careful in production environments.

## Troubleshooting

### Common Issues

| Problem | Solution |
|---------|----------|
| "Connection refused" | Ensure API server is running on port 8080 |
| "Database connection failed" | Check DATABASE_URL in .env file |
| "JWT secret missing" | Add JWT_SECRET to .env file |
| "User not found" | Use option 3 to verify user was created correctly |

### Getting Help

1. Check the main project [README](../README.md) for setup instructions
2. Verify your `.env` file matches `.env.example`
3. Ensure database is properly migrated with `bixor database status`
4. Check API server logs for detailed error messages 