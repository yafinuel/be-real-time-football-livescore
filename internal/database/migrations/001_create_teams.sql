CREATE TABLE IF NOT EXISTS teams (
    id VARCHAR(64) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    short_name VARCHAR(10) NOT NULL,
    logo_url TEXT
);