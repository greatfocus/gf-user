DO $$
BEGIN 
    -- Seed System User	
    IF NOT EXISTS (SELECT 1 FROM users WHERE identifier = 'muthurimixphone@gmail.com') THEN
        INSERT INTO users (identifier, password, status, deleted, enabled)
        VALUES('muthurimixphone@gmail.com', PGP_SYM_ENCRYPT('$2a$04$cZT44rp5yKqGox31VZpxieNq/XfoSJAMqoodhI/gUBvNvcn.kUUWe','28922DHSderer3244D7832$wCJSH]]DS2'), 'approved', false, true);
    END IF;
END $$