DO $$
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM country WHERE code = 'RW') THEN
        INSERT INTO country (name, code, prefix)
        VALUES
            ('Rwanda', 'RW', 250)
        ON CONFLICT (prefix) 
        DO NOTHING;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM country WHERE code = 'KE') THEN
        INSERT INTO country (name, code, prefix)
        VALUES
            ('Kenya', 'KE', 254)	
        ON CONFLICT (prefix) 
        DO NOTHING;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM country WHERE code = 'TZ') THEN
        INSERT INTO country (name, code, prefix)
        VALUES	
            ('Tanzania', 'TZ', 255)	
        ON CONFLICT (prefix) 
        DO NOTHING;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM country WHERE code = 'UG') THEN
        INSERT INTO country (name, code, prefix)
        VALUES
            ('Uganda', 'UG', 256)	
        ON CONFLICT (prefix) 
        DO NOTHING;
    END IF;
END $$