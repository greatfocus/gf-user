CREATE TABLE IF NOT EXISTS rights (
	id SERIAL PRIMARY KEY,
	roleId INTEGER NOT NULL REFERENCES role(id) ON DELETE CASCADE,
	userId INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	deleted BOOLEAN NOT NULL default(false),
	enabled BOOLEAN NOT NULL default(false),
	UNIQUE(roleId, userId)
);

DO $$ 
DECLARE
	roleId INTEGER := (select id from role where name='admin');
	userId INTEGER := (select id from users where mobilenumber='0780904371');
BEGIN 
	INSERT INTO rights (roleId, userId, deleted, enabled)
	VALUES
		(roleId, userId, false, true)
	ON CONFLICT
	DO NOTHING;
END $$;