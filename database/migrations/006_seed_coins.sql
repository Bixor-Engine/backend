-- Seed initial coins
INSERT INTO coins (name, ticker, decimal, price_decimal, logo, price, status, created_at, updated_at) VALUES 
('Bitcoin', 'BTC', 8, 2, 'https://cryptologos.cc/logos/bitcoin-btc-logo.png', 95000, 1, NOW(), NOW()),
('Ethereum', 'ETH', 18, 2, 'https://cryptologos.cc/logos/ethereum-eth-logo.png', 3200, 1, NOW(), NOW()),
('Tether', 'USDT', 6, 2, 'https://cryptologos.cc/logos/tether-usdt-logo.png', 1, 1, NOW(), NOW()),
('Binance Coin', 'BNB', 18, 2, 'https://cryptologos.cc/logos/binance-coin-bnb-logo.png', 600, 1, NOW(), NOW()),
('Solana', 'SOL', 9, 2, 'https://cryptologos.cc/logos/solana-sol-logo.png', 140, 1, NOW(), NOW());

-- Record this migration
INSERT INTO schema_migrations (version, description, checksum) 
VALUES ('006', 'Seed initial coins data', 'migration_006_seed_coins')
ON CONFLICT (version) DO NOTHING;
