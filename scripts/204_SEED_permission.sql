DO $$
DECLARE
	adminroleid INTEGER := (select id from role where name='Admin');
	staffroleid INTEGER := (select id from role where name='Staff');

	user_create INTEGER := (select id from action where name='user_create');
	user_activate INTEGER := (select id from action where name='user_activate');
	user_deactivate INTEGER := (select id from action where name='user_deactivate');
	user_delete INTEGER := (select id from action where name='user_delete');
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM permission WHERE roleid=adminroleid) THEN
       INSERT INTO permission (roleid, actionid)
        VALUES
            -- admin		
            (adminroleid, user_create),
            (adminroleid, user_activate),
            (adminroleid, user_deactivate),
            (adminroleid, user_delete)
        ON CONFLICT
        DO NOTHING;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM permission WHERE roleid=staffroleid) THEN
       INSERT INTO permission (roleid, actionid)
        VALUES
            -- staff
            (staffroleid, user_activate),
            (staffroleid, user_deactivate)
        ON CONFLICT
        DO NOTHING;
    END IF;
END $$