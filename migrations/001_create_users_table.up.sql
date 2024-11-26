CREATE TABLE users (
                       id BIGINT PRIMARY KEY NOT NULL ,
                       phone VARCHAR UNIQUE,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       name VARCHAR(100),
                       last_name VARCHAR(100),
                       second_name VARCHAR(100),
                       birth_date DATE,
                       sex VARCHAR(10)
);
