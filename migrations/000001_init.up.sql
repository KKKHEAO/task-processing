CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA sh_task;

CREATE TABLE sh_task.tasks
(
    id         UUID PRIMARY KEY,
    type       TEXT      NOT NULL,
    payload    BYTEA     NOT NULL,
    status     TEXT      NOT NULL,
    retries    INT       NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_tasks_status ON sh_task.tasks(status);