CREATE TABLE IF NOT EXISTS urls (
    id          TEXT PRIMARY KEY,
    original    TEXT NOT NULL,
    short_code  TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    clicks      INTEGER NOT NULL DEFAULT 0,
    expires_at  TIMESTAMPTZ
);
