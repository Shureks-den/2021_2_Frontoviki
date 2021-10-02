/*
DROP TABLE price_history;
DROP TABLE favorite;
DROP TABLE advert_image;
DROP TABLE advert;
DROP TABLE category;
DROP TABLE users;
*/

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email text UNIQUE NOT NULL,
    password text NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name text NOT NULL DEFAULT '',
    surname text NOT NULL DEFAULT '',
    image text NOT NULL DEFAULT '',
	rating DECIMAL(4, 2) DEFAULT 0.0
);


CREATE TABLE IF NOT EXISTS category (
	id SERIAL PRIMARY KEY,
	name text UNIQUE NOT NULL
);


CREATE TABLE IF NOT EXISTS advert (
    id SERIAL PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL DEFAULT '',
	price int NOT NULL DEFAULT 0,
    city text NOT NULL DEFAULT 'Moscow',
	latitude float NOT NULL DEFAULT 55.751244,
	longitude float NOT NULL DEFAULT 37.618423,
    published_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	date_close TIMESTAMP DEFAULT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
	views int NOT NULL DEFAULT 0,
	
    publisher_id INT NOT NULL,
	category_id INT NOT NULL,
	
    FOREIGN KEY (publisher_id) REFERENCES users (id) ON DELETE CASCADE,
	FOREIGN KEY (category_id) REFERENCES category (id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS advert_image (
	advert_id int NOT NULL,
	img_path text UNIQUE NOT NULL,
	
	FOREIGN KEY (advert_id) REFERENCES advert (id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS favorite (
	user_id int NOT NULL,
	advert_id int NOT NULL,
	
	FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
	FOREIGN KEY (advert_id) REFERENCES advert (id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS price_history (
	advert_id int NOT NULL,
	price int NOT NULL,
	change_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	
	FOREIGN KEY (advert_id) REFERENCES advert (id) ON DELETE CASCADE
);