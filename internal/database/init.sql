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

CREATE TABLE IF NOT EXISTS advts (
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL DEFAULT '',
    price INT NOT NULL DEFAULT 0,
    location text NOT NULL DEFAULT 'Moscow',
    published_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    image text NOT NULL DEFAULT '',
    publisher_id INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    FOREIGN KEY (publisher_id) REFERENCES users (id) ON DELETE CASCADE
);
