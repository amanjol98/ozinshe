CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    phone_number VARCHAR(100),
    born_at DATE
);

CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR,
    release_year INT,
    description VARCHAR,
    duration INT,
    poster_url VARCHAR,
    director VARCHAR,
    producer VARCHAR,
    video_id VARCHAR
);

CREATE TABLE genres (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE movie_genres (
    id SERIAL PRIMARY KEY,
    movie_id INT NOT NULL,
    genre_id INT NOT NULL,
    Foreign Key (movie_id) REFERENCES movies (id) ON DELETE CASCADE,
    Foreign Key (genre_id) REFERENCES genres (id) ON DELETE CASCADE,
    UNIQUE (movie_id, genre_id)
);

CREATE TABLE movie_categories (
    id SERIAL PRIMARY KEY,
    movie_id INT NOT NULL,
    category_id INT NOT NULL,
    Foreign Key (movie_id) REFERENCES movies (id) ON DELETE CASCADE,
    Foreign Key (category_id) REFERENCES categories (id) ON DELETE CASCADE,
    UNIQUE (movie_id, category_id)
);

CREATE TABLE favorites (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    movie_id INT NOT NULL,
    Foreign Key (user_id) REFERENCES users (id) ON DELETE CASCADE,
    Foreign Key (movie_id) REFERENCES movies (id) ON DELETE CASCADE,
    UNIQUE (user_id, movie_id)
);

CREATE TABLE seasons (
    id SERIAL PRIMARY KEY,
    movie_id INT NOT NULL,
    season_number INT,
    Foreign Key (movie_id) REFERENCES movies (id) ON DELETE CASCADE,
    UNIQUE (movie_id, season_number)
);

CREATE TABLE episodes (
    id SERIAL PRIMARY KEY,
    season_id INT NOT NULL,
    video_id VARCHAR,
    episode_number INT NOT NULL,
    Foreign Key (season_id) REFERENCES seasons (id) ON DELETE CASCADE,
    UNIQUE (season_id, episode_number)
);

CREATE TABLE movie_screenshots (
    id SERIAL PRIMARY KEY,
    movie_id INT NOT NULL,
    image_url VARCHAR,
    order_number INT,
    Foreign Key (movie_id) REFERENCES movies (id) ON DELETE CASCADE
);