CREATE TABLE IF NOT EXISTS comment (
    id SERIAL PRIMARY KEY,
    post_slug VARCHAR(255) NOT NULL,
    user_id INT REFERENCES user_account(id),
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    parent_id INT REFERENCES comment(id)
);
