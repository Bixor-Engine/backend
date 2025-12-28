-- Add global_balance column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS global_balance DECIMAL(20, 2) NOT NULL DEFAULT 0.00;

-- Record this migration
INSERT INTO schema_migrations (version, description, checksum) 
VALUES ('004', 'Add global_balance to users', 'migration_004_global_balance')
ON CONFLICT (version) DO NOTHING;
