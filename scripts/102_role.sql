CREATE TABLE IF NOT EXISTS role (
	id SERIAL PRIMARY KEY,
	name VARCHAR(20) NOT NULL,
	description VARCHAR(100) NOT NULL,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(name)
);

INSERT INTO role (name, description)
VALUES
	('admin', 'Role for admin'),
	('staff', 'Role for staff'),
	('partner', 'Role for partner'),	
	('customer', 'Role for customer')
ON CONFLICT (name) 
DO NOTHING;
