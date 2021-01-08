CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	type VARCHAR(20) NOT NULL,
	email VARCHAR(100) NOT NULL,	
	password VARCHAR(100) NOT NULL,
	failedAttempts INTEGER default (0),
	lastAttempt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	successLogins INTEGER default (0),
	expiredDate TIMESTAMP NOT NULL,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status VARCHAR(20) NOT NULL,
	deleted BOOLEAN NOT NULL default(false),
	enabled BOOLEAN NOT NULL default(false),
	UNIQUE(email),
	UNIQUE(email, password)
);