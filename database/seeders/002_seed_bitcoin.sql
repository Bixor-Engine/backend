-- Insert Bitcoin coin with bitcoind gateway for both deposit and withdraw
INSERT INTO coins (
    name,
    ticker,
    decimal,
    price_decimal,
    logo,
    price,
    deposit_gateway,
    withdraw_gateway,
    deposit_fee,
    withdraw_fee,
    deposit_fee_type,
    withdraw_fee_type,
    confirmation,
    status,
    withdraw_status,
    deposit_status,
    website,
    explorer,
    explorer_tx,
    explorer_address,
    created_at,
    updated_at
) VALUES (
    'Bitcoin',
    'BTC',
    8,  -- decimal places
    8,  -- price decimal places
    '/assets/coins/btc.png',  -- logo path
    50000.00000000,  -- price (example, update with real price)
    ARRAY['bitcoind'],  -- deposit gateways
    ARRAY['bitcoind'],  -- withdraw gateways
    0.00010000,  -- deposit_fee (example: 0.0001 BTC)
    0.00050000,  -- withdraw_fee (example: 0.0005 BTC)
    0,  -- deposit_fee_type: 0 = fixed, 1 = percentage
    0,  -- withdraw_fee_type: 0 = fixed, 1 = percentage
    6,  -- confirmation: 6 blocks
    1,  -- status: 1 = active
    1,  -- withdraw_status: 1 = enabled
    1,  -- deposit_status: 1 = enabled
    'https://bitcoin.org',  -- website
    'https://blockstream.info',  -- explorer
    'https://blockstream.info/tx/',  -- explorer_tx
    'https://blockstream.info/address/',  -- explorer_address
    NOW(),
    NOW()
)
ON CONFLICT (ticker) DO NOTHING;  -- Prevent duplicate if seeder runs multiple times

