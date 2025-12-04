--liquibase formatted sql
--changeset igor:3-classroom_name_fixes

ALTER TABLE classes
    ALTER COLUMN name TYPE VARCHAR(100) USING substring(name, 1, 100);

ALTER TABLE classes
    ADD CONSTRAINT chk_classroom_name_not_empty CHECK (length(trim(name)) > 0);