DO $$
DECLARE
	rId INTEGER := (select id from role where name='Admin');
	usId INTEGER := (select id from users where email='muthurimixphone@gmail.com');
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM rights r WHERE r.roleId=rId AND r.userId=usId) THEN
        INSERT INTO rights (roleId, userId, deleted, enabled)
        VALUES
            (roleId, userId, false, true)
        ON CONFLICT
        DO NOTHING;
    END IF;
END $$