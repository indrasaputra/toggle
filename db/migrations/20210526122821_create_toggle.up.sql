BEGIN;

CREATE TABLE IF NOT EXISTS toggles (
  id            BIGSERIAL       PRIMARY KEY,
  key           TEXT            UNIQUE NOT NULL,
  is_enabled    BOOLEAN,
  description   TEXT,
  created_at    TIMESTAMP,
  updated_at    TIMESTAMP
);

COMMIT;