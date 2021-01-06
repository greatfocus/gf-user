DO $$
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM role WHERE name = 'Admin') THEN
        INSERT INTO role (name, description)
        VALUES
            ('Admin', 'Role for admin')
        ON CONFLICT (name) 
        DO NOTHING;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM role WHERE name = 'Staff') THEN
        INSERT INTO role (name, description)
        VALUES
            ('Staff', 'Role for staff')
        ON CONFLICT (name) 
        DO NOTHING;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM role WHERE name = 'Partner') THEN
        INSERT INTO role (name, description)
        VALUES
            ('Partner', 'Role for partner')
        ON CONFLICT (name) 
        DO NOTHING;
    END IF;

    IF NOT EXISTS (SELECT 1 FROM role WHERE name = 'Customer') THEN
        INSERT INTO role (name, description)
        VALUES	
            ('Customer', 'Role for customer')
        ON CONFLICT (name) 
        DO NOTHING;
    END IF;
END $$