CREATE TABLE IF NOT EXISTS user_account (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    email VARCHAR(128) NOT NULL,
    profile_picture VARCHAR(256),
    provider VARCHAR(16),
    provider_id VARCHAR(64),
    created_at TIMESTAMPTZ DEFAULT now()
);
