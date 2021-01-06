DO $$
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM country WHERE code = 'KE') THEN
        INSERT INTO country (name, code, prefix)
        VALUES
            ('Rwanda', 'RW', 250),
            ('Kenya', 'KE', 254),	
            ('Tanzania', 'TZ', 255),
            ('Uganda', 'UG', 256)	
        ON CONFLICT (prefix) 
        DO NOTHING;
    END IF;
END $$