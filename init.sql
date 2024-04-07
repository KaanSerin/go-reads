CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    email VARCHAR(320),
    password TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)