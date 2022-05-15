CREATE TABLE IF NOT EXISTS token (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES user_account(id),
    refresh_token VARCHAR(512) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
