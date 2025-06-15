--liquibase formatted sql
--changeset igor:1-seed-roles

INSERT INTO roles (name) VALUES ('USER')
ON CONFLICT (name) DO NOTHING;

INSERT INTO roles (name) VALUES ('TEACHER')
ON CONFLICT (name) DO NOTHING;