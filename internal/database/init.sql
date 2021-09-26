CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username text NOT NULL,
    email text UNIQUE NOT NULL,
    password text NOT NULL,
    created_at TIMESTAMP NOT NULL,
    name text NOT NULL DEFAULT '',
    surname text NOT NULL DEFAULT '',
    image text NOT NULL DEFAULT ''
);