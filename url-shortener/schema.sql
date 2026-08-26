CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(255) UNIQUE,
    original_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE clicks (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(255) REFERENCES urls(short_code),
    clicked_at TIMESTAMP DEFAULT NOW(),
    device VARCHAR(255),
    browser VARCHAR(255),
    os VARCHAR(255),
    referrer VARCHAR(255)
);