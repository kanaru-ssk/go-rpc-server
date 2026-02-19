CREATE TABLE tasks (
    id         TEXT        PRIMARY KEY,
    title      TEXT        NOT NULL,
    status     TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);