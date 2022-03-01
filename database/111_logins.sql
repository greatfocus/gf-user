CREATE TABLE IF NOT EXISTS logins (
	id BIGSERIAL PRIMARY KEY,
	userId INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	lastAttempt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	failedAttempts INTEGER default (0),
	successLogins INTEGER default (0),
	UNIQUE(userid)
);

-- make this time series
-- columns user_id, login_timestamp, device_type in the time series database