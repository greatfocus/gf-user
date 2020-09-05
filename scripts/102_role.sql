CREATE TABLE IF NOT EXISTS role (
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(20) NOT NULL,
	description VARCHAR(100) NOT NULL,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(name)
);

INSERT INTO role (name, description)
VALUES
	('Admin', 'Role for admin'),
	('Staff', 'Role for staff'),
	('Partner', 'Role for partner'),	
	('Customer', 'Role for customer')
ON CONFLICT (name) 
DO NOTHING;
