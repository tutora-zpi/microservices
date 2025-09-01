--liquibase formatted sql
--changeset igor:2-extend-users

ALTER TABLE IF EXISTS users
    ADD COLUMN IF NOT EXISTS name VARCHAR(50);

ALTER TABLE IF EXISTS users
    ADD COLUMN IF NOT EXISTS surname VARCHAR(50);

ALTER TABLE IF EXISTS users
    ADD COLUMN IF NOT EXISTS avatar_key VARCHAR(500);