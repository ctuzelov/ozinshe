CREATE TABLE users (
    id             BIGSERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    email          VARCHAR(255) NOT NULL UNIQUE,
    number         VARCHAR(255) NOT NULL,
    date_of_birth  DATE NOT NULL,
    user_type      VARCHAR(255) NOT NULL,
    password       VARCHAR(255) NOT NULL,
    token          TEXT NOT NULL,
    refresh_token  TEXT NOT NULL
);
