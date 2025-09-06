--liquibase formatted sql
--changeset igor:1-create-schema

CREATE TABLE classes (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP
);

CREATE TABLE member_roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    CONSTRAINT uq_member_roles_name UNIQUE (name)
);

CREATE TABLE members (
    id UUID PRIMARY KEY,
    class_id UUID REFERENCES classes(id) ON DELETE CASCADE NOT NULL,
    user_id UUID NOT NULL,
    user_role SERIAL REFERENCES member_roles(id) NOT NULL
);

CREATE UNIQUE INDEX idx_members_user_id_class_id ON members(user_id, class_id);
