CREATE TABLE sh_task.outbox
(
    id         UUID PRIMARY KEY,
    topic      TEXT      NOT NULL,
    key        TEXT      NOT NULL,
    payload    BYTEA     NOT NULL,
    created_at TIMESTAMP NOT NULL,
    processed  BOOLEAN   NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_outbox_processed ON sh_task.outbox(processed);