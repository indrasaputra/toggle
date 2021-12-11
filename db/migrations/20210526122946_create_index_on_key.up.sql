BEGIN;

CREATE INDEX IF NOT EXISTS index_on_key_on_toggles ON toggles USING btree (key);

COMMIT;
