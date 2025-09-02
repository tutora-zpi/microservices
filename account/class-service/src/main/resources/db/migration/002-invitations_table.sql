--liquibase formatted sql
--changeset igor:2-invitations_table

CREATE TABLE invitation_statuses (
    id SERIAL PRIMARY KEY,
    status_name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE invitations (
    id SERIAL PRIMARY KEY,
    class_id UUID NOT NULL REFERENCES classes(id),
    user_id UUID NOT NULL,
    status SERIAL REFERENCES invitation_statuses(id) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(class_id, user_id)
);
