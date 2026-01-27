-- 0003_create_nodes.sql
CREATE TABLE IF NOT EXISTS nodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    location TEXT NOT NULL,
    address TEXT NOT NULL,
    port INTEGER NOT NULL DEFAULT 443,
    server_name TEXT,
    protocols TEXT,
    available BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Index for available nodes
CREATE INDEX IF NOT EXISTS idx_nodes_available ON nodes (available);
