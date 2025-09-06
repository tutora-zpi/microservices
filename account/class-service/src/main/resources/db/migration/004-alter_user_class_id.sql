--liquibase formatted sql
--changeset igor:5-alter_id_user_class

ALTER TABLE users_classes
    ADD COLUMN id_new UUID;

ALTER TABLE users_classes
    DROP COLUMN id;

ALTER TABLE users_classes
    RENAME COLUMN id_new TO id;

ALTER TABLE users_classes
    ADD PRIMARY KEY (id);