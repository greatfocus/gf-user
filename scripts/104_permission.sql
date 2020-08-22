CREATE TABLE IF NOT EXISTS permission (
	id SERIAL PRIMARY KEY,
	roleId INTEGER NOT NULL REFERENCES role(id) ON DELETE CASCADE,
	actionId INTEGER NOT NULL REFERENCES action(id) ON DELETE CASCADE,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(roleId, actionId)
);

DO $$ 
DECLARE
	adminroleid INTEGER := (select id from role where name='admin');
	staffroleid INTEGER := (select id from role where name='staff');

	user_create INTEGER := (select id from action where name='user_create');
	user_activate INTEGER := (select id from action where name='user_activate');
	user_deactivate INTEGER := (select id from action where name='user_deactivate');
	user_delete INTEGER := (select id from action where name='user_delete');
	
BEGIN 
	INSERT INTO permission (roleid, actionid)
	VALUES
		-- admin		
		(adminroleid, user_create),
		(adminroleid, user_activate),
		(adminroleid, user_deactivate),
		(adminroleid, user_delete),

		-- staff
		(staffroleid, user_activate),
		(staffroleid, user_deactivate)

		-- agent
		-- customer
	ON CONFLICT
	DO NOTHING;
END $$;