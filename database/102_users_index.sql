CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_users_id ON users USING BTREE(id DESC);