CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       name VARCHAR(100) NOT NULL,
                       last_name VARCHAR(100),
                       second_name VARCHAR(100),
                       birth_date DATE,
                       sex VARCHAR(10)
);
