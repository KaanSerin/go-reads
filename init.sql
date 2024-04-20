CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO roles (name)
VALUES ('admin'),
    ('user');
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    email VARCHAR(320),
    password TEXT,
    role_id INT DEFAULT 2 constraint users_roles_id_fk references roles,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    author VARCHAR(100) NOT NULL,
    genre VARCHAR(50) NOT NULL,
    publication_date DATE NOT NULL,
    publisher VARCHAR(100) NOT NULL,
    isbn VARCHAR(50) NOT NULL,
    page_count SMALLINT NOT NULL,
    language VARCHAR(50) NOT NULL,
    format VARCHAR(50) NOT NULL
);