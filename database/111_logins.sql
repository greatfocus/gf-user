CREATE TABLE IF NOT EXISTS logins (
	id BIGSERIAL PRIMARY KEY,
	userId INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	type VARCHAR(20) NOT NULL,
	sessionId INTEGER NULL,
	lastAttempt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	loginAttempts INTEGER default (0),
	failedLogins INTEGER default (0),
	successLogins INTEGER default (0),
	UNIQUE(userid)
);
