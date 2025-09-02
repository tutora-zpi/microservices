--liquibase formatted sql
--changeset igor:3-unique-role-name

ALTER TABLE user_class_roles
    ADD CONSTRAINT uq_user_class_roles_name UNIQUE (name);