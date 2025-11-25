-- Create coins table for cryptocurrency management
CREATE TABLE IF NOT EXISTS coins (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    ticker VARCHAR(10) NOT NULL,
    decimal INTEGER NOT NULL,
    price_decimal INTEGER NOT NULL,
    logo VARCHAR(500),
    price NUMERIC(20, 8) NOT NULL DEFAULT 0,
    deposit_gateway TEXT[],
    withdraw_gateway TEXT[],
    deposit_fee NUMERIC(20, 8),
    withdraw_fee NUMERIC(20, 8),
    deposit_fee_type INTEGER,
    withdraw_fee_type INTEGER,
    confirmation INTEGER,
    status INTEGER NOT NULL DEFAULT 1,
    withdraw_status INTEGER,
    deposit_status INTEGER,
    website VARCHAR(500),
    explorer VARCHAR(500),
    explorer_tx VARCHAR(500),
    explorer_address VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_coins_ticker ON coins(ticker);
CREATE INDEX IF NOT EXISTS idx_coins_status ON coins(status);
CREATE INDEX IF NOT EXISTS idx_coins_deposit_status ON coins(deposit_status);
CREATE INDEX IF NOT EXISTS idx_coins_withdraw_status ON coins(withdraw_status);
CREATE INDEX IF NOT EXISTS idx_coins_name ON coins(name);

-- Add unique constraint for ticker (each coin ticker should be unique)
ALTER TABLE coins ADD CONSTRAINT chk_coins_ticker_unique UNIQUE (ticker);

-- Add check constraints for status values (assuming 0 = inactive, 1 = active)
ALTER TABLE coins ADD CONSTRAINT chk_coins_status 
CHECK (status IN (0, 1));

-- Add check constraints for deposit/withdraw status (0 = disabled, 1 = enabled)
ALTER TABLE coins ADD CONSTRAINT chk_coins_deposit_status 
CHECK (deposit_status IS NULL OR deposit_status IN (0, 1));

ALTER TABLE coins ADD CONSTRAINT chk_coins_withdraw_status 
CHECK (withdraw_status IS NULL OR withdraw_status IN (0, 1));

-- Add check constraint for decimal values (must be positive)
ALTER TABLE coins ADD CONSTRAINT chk_coins_decimal 
CHECK (decimal >= 0 AND decimal <= 18);

ALTER TABLE coins ADD CONSTRAINT chk_coins_price_decimal 
CHECK (price_decimal >= 0 AND price_decimal <= 18);

-- Add check constraint for fee types (assuming 0 = fixed, 1 = percentage)
ALTER TABLE coins ADD CONSTRAINT chk_coins_deposit_fee_type 
CHECK (deposit_fee_type IS NULL OR deposit_fee_type IN (0, 1));

ALTER TABLE coins ADD CONSTRAINT chk_coins_withdraw_fee_type 
CHECK (withdraw_fee_type IS NULL OR withdraw_fee_type IN (0, 1));

-- Add check constraint for confirmation (must be non-negative)
ALTER TABLE coins ADD CONSTRAINT chk_coins_confirmation 
CHECK (confirmation IS NULL OR confirmation >= 0);

-- Add check constraint for price (must be non-negative)
ALTER TABLE coins ADD CONSTRAINT chk_coins_price 
CHECK (price >= 0);

-- Add check constraint for fees (must be non-negative if provided)
ALTER TABLE coins ADD CONSTRAINT chk_coins_deposit_fee 
CHECK (deposit_fee IS NULL OR deposit_fee >= 0);

ALTER TABLE coins ADD CONSTRAINT chk_coins_withdraw_fee 
CHECK (withdraw_fee IS NULL OR withdraw_fee >= 0);

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_coins_updated_at 
    BEFORE UPDATE ON coins 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Record this migration
INSERT INTO schema_migrations (version, description, checksum) 
VALUES ('003', 'Create coins table for cryptocurrency management', 'migration_003_coins_table')
ON CONFLICT (version) DO NOTHING;

