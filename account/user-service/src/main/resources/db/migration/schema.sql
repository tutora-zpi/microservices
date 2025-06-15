--liquibase formatted sql
--changeset igor:1-create-schema

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    provider VARCHAR(20) NOT NULL,
    provider_id VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX idx_users_provider_providerid ON users (provider, provider_id);
CREATE UNIQUE INDEX idx_users_email ON users (email);

CREATE TABLE roles (
   id BIGSERIAL PRIMARY KEY,
   name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE users_roles (
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    role_id BIGINT REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id)
);
