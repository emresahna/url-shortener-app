CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS urls;

CREATE TABLE users (
    id uuid PRIMARY KEY default uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    is_deleted bool default false,
    is_active bool default true,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp default null,
    deleted_At timestamp default null
);

CREATE TABLE urls (
    id uuid PRIMARY KEY default uuid_generate_v4(),
    original_url TEXT NOT NULL,
    shortened_code VARCHAR(10) UNIQUE NOT NULL,
    user_id uuid REFERENCES users(id) default null,
    is_deleted bool default false,
    is_active bool default true,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp default null,
    deleted_At timestamp default null
);

CREATE TABLE click_counts (
    url_id uuid PRIMARY KEY REFERENCES urls(id) ON DELETE CASCADE,
    total_clicks bigint NOT NULL DEFAULT 0
);

CREATE INDEX idx_shortened_code ON urls(shortened_code);
CREATE INDEX idx_user_id ON urls(user_id);