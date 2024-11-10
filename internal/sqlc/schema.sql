CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS urls;

CREATE TABLE users (
    id uuid PRIMARY KEY default uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_deleted bool default false,
    is_active bool default true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz default null,
    deleted_At timestamptz default null
);

CREATE TABLE urls (
    id uuid PRIMARY KEY default uuid_generate_v4(),
    original_url TEXT NOT NULL,
    shortened_code VARCHAR(10) UNIQUE NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expire_time timestamptz default null,
    is_deleted bool default false,
    is_active bool default true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz default null,
    deleted_At timestamptz default null
);

CREATE INDEX idx_shortened_code ON urls(shortened_code);
CREATE INDEX idx_user_id ON urls(user_id);