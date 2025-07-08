-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email_status BOOLEAN DEFAULT FALSE NOT NULL,
    phone_number TEXT,
    phone_status BOOLEAN DEFAULT FALSE NOT NULL,
    referred_by UUID REFERENCES users(id) ON DELETE SET NULL,
    address TEXT,
    city TEXT,
    country TEXT,
    role TEXT DEFAULT 'user' NOT NULL,
    status TEXT DEFAULT 'pending' NOT NULL,
    kyc_status TEXT DEFAULT 'not_submitted' NOT NULL,
    twofa_enabled BOOLEAN DEFAULT FALSE NOT NULL,
    last_login_at TIMESTAMP WITH TIME ZONE,
    last_login_ip TEXT,
    device_info JSONB,
    language TEXT DEFAULT 'en' NOT NULL,
    timezone TEXT DEFAULT 'UTC',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_referred_by ON users(referred_by);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_kyc_status ON users(kyc_status);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Add constraints for enumerated values
ALTER TABLE users ADD CONSTRAINT chk_users_role 
CHECK (role IN ('user', 'superadmin', 'admin', 'technical', 'accountant', 'compliance', 'support'));

ALTER TABLE users ADD CONSTRAINT chk_users_status 
CHECK (status IN ('active', 'suspended', 'banned', 'pending', 'inactive', 'locked', 'frozen'));

ALTER TABLE users ADD CONSTRAINT chk_users_kyc_status 
CHECK (kyc_status IN ('not_submitted', 'pending', 'verified', 'rejected'));

-- Add email format validation
ALTER TABLE users ADD CONSTRAINT chk_users_email_format 
CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$');

-- Create function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Create schema_migrations table if it doesn't exist
CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY,
    description TEXT NOT NULL,
    applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    applied_by TEXT DEFAULT CURRENT_USER NOT NULL,
    checksum TEXT
);

-- Record this migration
INSERT INTO schema_migrations (version, description, checksum) 
VALUES ('001', 'Create users table with all required columns and constraints', 'migration_001_users_table')
ON CONFLICT (version) DO NOTHING; 