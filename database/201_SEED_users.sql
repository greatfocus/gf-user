DO $$
BEGIN 
    -- Seed System User	
    IF NOT EXISTS (SELECT 1 FROM users WHERE identifier = 'muthurimixphone@gmail.com') THEN
        INSERT INTO users (type, status, identifier, password, deleted, enabled, system)
        VALUES('password', 'approved', 'muthurimixphone@gmail.com', '$2a$04$cZT44rp5yKqGox31VZpxieNq/XfoSJAMqoodhI/gUBvNvcn.kUUWe', false, true, true);
    END IF;
END $$