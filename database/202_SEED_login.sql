DO $$
DECLARE
	usrId INTEGER := (SELECT id FROM users WHERE identifier='muthurimixphone@gmail.com');
BEGIN 
    -- Seed System User	
    IF NOT EXISTS (SELECT 1 FROM logins WHERE userId=usrId) THEN
        INSERT INTO logins (userId, lastAttempt, failedAttempts, successLogins)
        VALUES(usrId, CURRENT_TIMESTAMP, 0, 0);
    END IF;
END $$