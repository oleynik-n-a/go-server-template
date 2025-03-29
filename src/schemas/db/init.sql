CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    token TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE
);

-- Test user (in the table password is encrypted by bcrypt)
-- {
--     "email": "admin@example.com",
--     "password": "password"
-- }

INSERT INTO users (email, password)
VALUES ('admin@example.com', '$2a$12$njzaAnQEULN1kWAd4CyG.eyK/rf4DBXLCK66VL/NTcBK0HeBQXkHe')
ON CONFLICT (email) DO NOTHING;
