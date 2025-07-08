-- Insert superadmin user
-- Password is hashed using Argon2id via PostgreSQL crypt() function
INSERT INTO users (
    id,
    first_name,
    last_name,
    username,
    email,
    password,
    email_status,
    phone_number,
    phone_status,
    referred_by,
    address,
    city,
    country,
    role,
    status,
    kyc_status,
    twofa_enabled,
    last_login_at,
    last_login_ip,
    device_info,
    language,
    timezone,
    created_at,
    updated_at,
    deleted_at
) VALUES (
    gen_random_uuid(),
    'Super',
    'Admin',
    'superadmin',
    'superadmin@bixor.io',
    '$argon2i$v=19$m=16,t=2,p=1$MzFocW1NVXpzMEtoVGJ6eg$aVXX9rnNCUK8KeJZu+c7XA', -- Argon2 hash of: H5w8xLLz}0drN&ey{0?!
    TRUE,  -- email_status: verified
    NULL,  -- phone_number: nullable
    FALSE, -- phone_status: not verified since no phone
    NULL,  -- referred_by: nullable
    NULL,  -- address: nullable
    NULL,  -- city: nullable
    NULL,  -- country: nullable
    'superadmin', -- role: superadmin
    'active',     -- status: active
    'verified',   -- kyc_status: verified
    FALSE, -- twofa_enabled: not enabled by default
    NULL,  -- last_login_at: nullable
    NULL,  -- last_login_ip: nullable
    NULL,  -- device_info: nullable
    'en',  -- language: default English
    'UTC', -- timezone: default UTC
    NOW(), -- created_at: current timestamp
    NOW(), -- updated_at: current timestamp
    NULL   -- deleted_at: nullable (not deleted)
)
ON CONFLICT (username) DO NOTHING; -- Prevent duplicate if seeder runs multiple times

-- Also prevent duplicate by email
INSERT INTO users (
    id,
    first_name,
    last_name,
    username,
    email,
    password,
    email_status,
    phone_number,
    phone_status,
    referred_by,
    address,
    city,
    country,
    role,
    status,
    kyc_status,
    twofa_enabled,
    last_login_at,
    last_login_ip,
    device_info,
    language,
    timezone,
    created_at,
    updated_at,
    deleted_at
) 
SELECT 
    gen_random_uuid(),
    'Super',
    'Admin',
    'superadmin',
    'superadmin@bixor.io',
    'H5w8xLLz}0drN&ey{0?!',
    TRUE,
    NULL,
    FALSE,
    NULL,
    NULL,
    NULL,
    NULL,
    'superadmin',
    'active',
    'verified',
    FALSE,
    NULL,
    NULL,
    NULL,
    'en',
    'UTC',
    NOW(),
    NOW(),
    NULL
WHERE NOT EXISTS (
    SELECT 1 FROM users WHERE email = 'superadmin@bixor.io'
); 