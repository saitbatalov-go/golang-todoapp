CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE SCHEMA todoapp;

CREATE TABLE todoapp.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version INTEGER NOT NULL DEFAULT 1,
    full_name VARCHAR(100) NOT NULL CHECK(char_length(full_name) BETWEEN 1 AND 100),
    phone_number VARCHAR(15) CHECK(
        phone_number ~ '^\+[0-9]+$'
        AND
        char_length(phone_number) BETWEEN 10 AND 15
    )
);

CREATE TABLE todoapp.tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version INTEGER NOT NULL DEFAULT 1,
    title VARCHAR(100) NOT NULL CHECK(char_length(title) BETWEEN 1 AND 100),
    description VARCHAR(1000) CHECK(char_length(description) BETWEEN 1 AND 1000),
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    author_user_id UUID NOT NULL REFERENCES todoapp.users(id) ON DELETE CASCADE,
    CHECK (
        (completed = TRUE AND completed_at IS NOT NULL AND completed_at >= created_at)
        OR
        (completed = FALSE AND completed_at IS NULL)
    )
);

CREATE INDEX idx_tasks_author_user_id ON todoapp.tasks(author_user_id);
CREATE INDEX idx_tasks_completed ON todoapp.tasks(completed);
CREATE INDEX idx_users_phone_number ON todoapp.users(phone_number);