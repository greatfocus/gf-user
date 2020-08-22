CREATE TABLE IF NOT EXISTS userType (
	id VARCHAR(20) NOT NULL,
	name VARCHAR(20) NOT NULL,
	description VARCHAR(200) NULL,
	UNIQUE(id),
	UNIQUE(name)
);


INSERT INTO userType (id, name)
VALUES
	('admin', 'Administrator', 'Super user in the system'),
	('staff', 'Staff', 'Company Staff'),
	('partner', 'Partner', 'Our partners to grow together'),
	('customer', 'Customer', 'Direct customer')
ON CONFLICT (id) 
DO NOTHING; 