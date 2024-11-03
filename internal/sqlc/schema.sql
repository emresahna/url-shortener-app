DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS urls;

CREATE TABLE users (
    id uuid PRIMARY KEY default uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE urls (
    id uuid PRIMARY KEY default uuid_generate_v4(),
    original_url TEXT NOT NULL,
    shortened_code VARCHAR(10) UNIQUE NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_shortened_code ON urls(shortened_code);
CREATE INDEX idx_user_id ON urls(user_id);