CREATE TABLE IF NOT EXISTS rights (
	id SERIAL PRIMARY KEY,
	roleId INTEGER NOT NULL REFERENCES role(id) ON DELETE CASCADE,
	userId INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status VARCHAR(20) NOT NULL,
	deleted BOOLEAN NOT NULL default(false),
	enabled BOOLEAN NOT NULL default(false),
	UNIQUE(roleId, userId)
);

DO $$ 
DECLARE
	roleId INTEGER := (select id from role where name='admin');
	userId INTEGER := (select id from users where email='mucunga90@gmail.com');
BEGIN 
	INSERT INTO rights (roleId, userId, status, deleted, enabled)
	VALUES
		(roleId, userId, 'RIGHT.APPROVED', false, true)
	ON CONFLICT
	DO NOTHING;
END $$;