-- Create OTP (One-Time Password) table for email verification and other verification purposes
CREATE TABLE IF NOT EXISTS otps (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL, -- e.g., 'email-verification', 'password-reset', '2fa'
    code TEXT NOT NULL, -- 6-digit code
    used BOOLEAN DEFAULT FALSE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_otps_user_id ON otps(user_id);
CREATE INDEX IF NOT EXISTS idx_otps_type ON otps(type);
CREATE INDEX IF NOT EXISTS idx_otps_code ON otps(code);
CREATE INDEX IF NOT EXISTS idx_otps_used ON otps(used);
CREATE INDEX IF NOT EXISTS idx_otps_expires_at ON otps(expires_at);
CREATE INDEX IF NOT EXISTS idx_otps_user_type ON otps(user_id, type);
CREATE INDEX IF NOT EXISTS idx_otps_user_type_used ON otps(user_id, type, used);

-- Add constraint for OTP type values
ALTER TABLE otps ADD CONSTRAINT chk_otps_type 
CHECK (type IN ('email-verification', 'password-reset', '2fa', 'phone-verification'));

-- Add constraint to ensure code is 6 digits
ALTER TABLE otps ADD CONSTRAINT chk_otps_code_format 
CHECK (code ~ '^[0-9]{6}$');

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_otps_updated_at 
    BEFORE UPDATE ON otps 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Record this migration
INSERT INTO schema_migrations (version, description, checksum) 
VALUES ('002', 'Create OTP table for email verification and other verification purposes', 'migration_002_otp_table')
ON CONFLICT (version) DO NOTHING;

