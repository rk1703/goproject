CREATE DATABASE instagram;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    password VARCHAR(100)
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    caption TEXT,
    image_url TEXT,
    posted_at TIMESTAMP,
    user_id INT REFERENCES users(id)
);
