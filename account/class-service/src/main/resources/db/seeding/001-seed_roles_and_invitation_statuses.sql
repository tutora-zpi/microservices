--liquibase formatted sql
--changeset igor:3-seed-invitation-statuses-and-roles

-- Invitation statuses
INSERT INTO invitation_statuses (status_name) VALUES ('ACCEPTED')
ON CONFLICT (status_name) DO NOTHING;

INSERT INTO invitation_statuses (status_name) VALUES ('DECLINED')
ON CONFLICT (status_name) DO NOTHING;

INSERT INTO invitation_statuses (status_name) VALUES ('INVITED')
ON CONFLICT (status_name) DO NOTHING;

-- User class roles
INSERT INTO member_roles (name) VALUES ('HOST')
ON CONFLICT (name) DO NOTHING;

INSERT INTO member_roles (name) VALUES ('GUEST')
ON CONFLICT (name) DO NOTHING;