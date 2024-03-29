CREATE TABLE IF NOT EXISTS permission (
	id BIGSERIAL PRIMARY KEY,
	actionId INTEGER NOT NULL REFERENCES action(id) ON DELETE CASCADE,
	userId INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted BOOLEAN NOT NULL default(false),
	enabled BOOLEAN NOT NULL default(true),
	UNIQUE(actionId, userId)
);