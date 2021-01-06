DO $$
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM role WHERE name = 'Admin') THEN
        INSERT INTO role (name, description)
        VALUES
            ('Admin', 'Role for admin'),
            ('Staff', 'Role for staff'),
            ('Partner', 'Role for partner'),	
            ('Customer', 'Role for customer')
        ON CONFLICT (name) 
        DO NOTHING;
    END IF;
END $$