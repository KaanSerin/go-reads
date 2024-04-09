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