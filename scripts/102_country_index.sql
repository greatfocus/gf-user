CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS idx_country_code ON country USING BTREE(code DESC);