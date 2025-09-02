--liquibase formatted sql
--changeset igor:1-create-schema

CREATE TABLE classes (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP
);

CREATE TABLE user_class_roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE users_classes (
    id BIGINT PRIMARY KEY,
    class_id UUID REFERENCES classes(id) ON DELETE CASCADE NOT NULL,
    user_id UUID NOT NULL,
    user_role SERIAL REFERENCES user_class_roles(id) NOT NULL
);

CREATE UNIQUE INDEX idx_users_classes_user_id_class_id ON users_classes(user_id, class_id);
