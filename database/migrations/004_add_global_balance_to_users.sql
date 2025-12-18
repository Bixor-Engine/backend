-- Add global_balance column to users table
ALTER TABLE users ADD COLUMN global_balance DECIMAL(20, 2) NOT NULL DEFAULT 0.00;
