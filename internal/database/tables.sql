-- DROP TABLE views_;
-- DROP TABLE rating_statistics;
-- DROP TABLE rating;
-- DROP TABLE cart;
-- DROP TABLE price_history;
-- DROP TABLE favorite;
-- DROP TABLE advert_image;
-- DROP TABLE advert;
-- DROP TABLE category;
-- DROP TABLE users;


-- CREATE TEXT SEARCH DICTIONARY russian_ispell (
--     TEMPLATE = ispell,
--     DictFile = russian,
--     AffFile = russian,
--     StopWords = russian
-- );

-- CREATE TEXT SEARCH CONFIGURATION ru (COPY=russian);

-- ALTER TEXT SEARCH CONFIGURATION ru
--     ALTER MAPPING FOR hword, hword_part, word
--     WITH russian_ispell, russian_stem;

-- CREATE EXTENSION postgis;
-- CREATE EXTENSION postgis_topology;


CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email text UNIQUE NOT NULL,
	phone text NOT NULL DEFAULT '',
    password text NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    name text NOT NULL DEFAULT '',
    surname text NOT NULL DEFAULT '',
    image text NOT NULL DEFAULT ''
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
    location text NOT NULL DEFAULT 'Moscow',
	latitude float NOT NULL DEFAULT 55.751244,
	longitude float NOT NULL DEFAULT 37.618423,
    published_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	date_close TIMESTAMP NOT NULL DEFAULT to_timestamp(0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
	views int NOT NULL DEFAULT 0,
	amount int NOT NULL DEFAULT 1,
	is_new BOOLEAN NOT NULL DEFAULT TRUE,

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

CREATE TABLE IF NOT EXISTS cart (
	user_id int NOT NULL,
	advert_id int NOT NULL,
	amount int NOT NULL,

	FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
	FOREIGN KEY (advert_id) REFERENCES advert (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS rating (
	user_from int NOT NULL,
	user_to int NOT NULL,
	rate int NOT NULL,

	FOREIGN KEY (user_from) REFERENCES users (id) ON DELETE CASCADE,
	FOREIGN KEY (user_to) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS rating_statistics (
	user_id int NOT NULL,
	sum int NOT NULL DEFAULT 0,
	count int NOT NULL DEFAULT 0,

	FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS views_ (
	advert_id int NOT NULL,
	count int NOT NULL DEFAULT 0,

	FOREIGN KEY (advert_id) REFERENCES advert (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS messages (
	user_from int NOT NULL,
	user_to int NOT NULL,
	msg VARCHAR(255),

	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (user_from) REFERENCES users (id) ON DELETE RESTRICT,
	FOREIGN KEY (user_to) REFERENCES users (id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS dialogs (
	user1 int NOT NULL,
	user2 int NOT NULL,

	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (user1) REFERENCES users (id) ON DELETE RESTRICT,
	FOREIGN KEY (user2) REFERENCES users (id) ON DELETE RESTRICT
);


-- INSERT INTO category (name) values ('одежда'), ('обувь'), ('животные');
-- INSERT INTO advert (name, publisher_id, category_id) values ('Худи спортивная', 2, 1), ('Манчкин', 1, 3);
-- INSERT INTO advert_image (advert_id, img_path) VALUES (2, 'hudi1'), (2, 'hudi2');
