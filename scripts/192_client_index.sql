CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_client_id ON client USING BTREE(id DESC);