CREATE EXTENSION IF NOT EXISTS hstore;

CREATE UNLOGGED TABLE IF NOT EXISTS cache_kv (
  key text PRIMARY KEY,
  value hstore NOT NULL,
  expires_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS cache_kv_expires_at_idx ON cache_kv (expires_at);
