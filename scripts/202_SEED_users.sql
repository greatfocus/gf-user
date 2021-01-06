DO $$
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM users WHERE type = 'admin') THEN
        INSERT INTO users (type, email, password, expiredDate, status, enabled)
        VALUES
            ('admin', 'muthurimixphone@gmail.com', '$2a$04$cZT44rp5yKqGox31VZpxieNq/XfoSJAMqoodhI/gUBvNvcn.kUUWe', '2025-06-01 08:22:17.460493', 'USER.APPROVED', true) 
        ON CONFLICT
        DO NOTHING;
    END IF;
END $$