CREATE TABLE IF NOT EXISTS users (
    id             BIGSERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    email          VARCHAR(255) NOT NULL UNIQUE,
    number         VARCHAR(255) NOT NULL,
    date_of_birth  DATE NOT NULL,
    user_type      VARCHAR(255) NOT NULL,
    password       VARCHAR(255) NOT NULL,
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    token          TEXT,
    refresh_token  TEXT
);

CREATE TABLE IF NOT EXISTS genres (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS age_categories (
    id SERIAL PRIMARY KEY,
    min_age INT NOT NULL,
    max_age INT NOT NULL
);

CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    release_year INT NOT NULL,
    description TEXT NOT NULL,
    popularity INT,
    youtube_id VARCHAR(11) UNIQUE,
    duration INT NOT NULL,
    director VARCHAR(100) NOT NULL,
    producer VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS series (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    release_year INT NOT NULL,
    description TEXT NOT NULL,
    popularity INT,
    duration INT NOT NULL,
    director VARCHAR(100) NOT NULL,
    producer VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS seasons (
    id SERIAL PRIMARY KEY,
    series_id INT NOT NULL,
    season_number INT NOT NULL,
    FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS episodes (
    id SERIAL PRIMARY KEY,
    season_id INT NOT NULL,
    episode_number INT NOT NULL,
    youtube_id VARCHAR(11),
    FOREIGN KEY (season_id) REFERENCES seasons(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS movie_covers(
    id SERIAL PRIMARY KEY,
    movie_id INTEGER NOT NULL,
    filename VARCHAR(255) NOT NULL,
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS series_covers(
    id SERIAL PRIMARY KEY,
    series_id INTEGER NOT NULL,
    filename VARCHAR(255) NOT NULL,
    FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS movie_screenshots(
    id SERIAL PRIMARY KEY,
    movie_id INTEGER NOT NULL,
    filename VARCHAR(255) NOT NULL,
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS series_screenshots(
    id SERIAL PRIMARY KEY,
    series_id INTEGER NOT NULL,
    filename VARCHAR(255) NOT NULL,
    FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS movie_genres (
    movie_id INT NOT NULL,
    genre_id INT NOT NULL,
    FOREIGN KEY (movie_id) REFERENCES movies(id),
    FOREIGN KEY (genre_id) REFERENCES genres(id)
);

CREATE TABLE IF NOT EXISTS series_genres (
    series_id INT NOT NULL,
    genre_id INT NOT NULL,
    FOREIGN KEY (series_id) REFERENCES series(id),
    FOREIGN KEY (genre_id) REFERENCES genres(id)
);

CREATE TABLE IF NOT EXISTS movie_age_categories (
    movie_id INT NOT NULL,
    age_category_id INT NOT NULL,
    FOREIGN KEY (movie_id) REFERENCES movies(id),
    FOREIGN KEY (age_category_id) REFERENCES age_categories(id)
);

CREATE TABLE IF NOT EXISTS series_age_categories (
    series_id INT NOT NULL,
    age_category_id INT NOT NULL,
    FOREIGN KEY (series_id) REFERENCES series(id),
    FOREIGN KEY (age_category_id) REFERENCES age_categories(id)
);

CREATE TABLE IF NOT EXISTS favorite_movies (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    movie_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS favorite_series (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    series_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS key_words (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS movie_key_words (
    movie_id INTEGER NOT NULL,
    key_word_id INTEGER NOT NULL,
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE,
    FOREIGN KEY (key_word_id) REFERENCES key_words(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS series_key_words (
    series_id INTEGER NOT NULL,
    key_word_id INTEGER NOT NULL,
    FOREIGN KEY (series_id) REFERENCES series(id) ON DELETE CASCADE,
    FOREIGN KEY (key_word_id) REFERENCES key_words(id) ON DELETE CASCADE
);

CREATE INDEX idx_movie_genres_genre_id ON movie_genres (genre_id);
CREATE INDEX idx_series_genres_genre_id ON series_genres (genre_id);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_movies_title ON movies (title);
CREATE INDEX idx_movies_release_year ON movies (release_year);
CREATE INDEX idx_movies_popularity ON movies (popularity);

CREATE INDEX idx_series_title ON series (title);
CREATE INDEX idx_series_release_year ON series (release_year);
CREATE INDEX idx_series_popularity ON series (popularity);