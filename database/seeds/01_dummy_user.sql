-- Seed data: Insert dummy user for testing/development
-- From migration 000011_create_s3_resources_mapping_table.up.sql

INSERT INTO users (username, password_hash, status)
VALUES ('dummy_user', 'dummy_hash', 'active');
