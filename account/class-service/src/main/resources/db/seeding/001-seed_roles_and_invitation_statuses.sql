--liquibase formatted sql
--changeset igor:4-seed-invitation-statuses-and-roles

-- Invitation statuses
INSERT INTO invitation_statuses (status_name) VALUES ('ACCEPTED')
ON CONFLICT (status_name) DO NOTHING;

INSERT INTO invitation_statuses (status_name) VALUES ('DECLINED')
ON CONFLICT (status_name) DO NOTHING;

INSERT INTO invitation_statuses (status_name) VALUES ('INVITED')
ON CONFLICT (status_name) DO NOTHING;

-- User class roles
INSERT INTO user_class_roles (name) VALUES ('HOST')
ON CONFLICT (name) DO NOTHING;

INSERT INTO user_class_roles (name) VALUES ('GUEST')
ON CONFLICT (name) DO NOTHING;