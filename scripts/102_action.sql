CREATE TABLE IF NOT EXISTS action (
	id SERIAL PRIMARY KEY,
	name VARCHAR(30) NOT NULL,
	description VARCHAR(100) NOT NULL,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(name)
);

INSERT INTO action (name, description)
VALUES
	('user_activate', 'Action to allow user activate'),
	('user_deactivate', 'Action to allow user deactivate')
ON CONFLICT (name) 
DO NOTHING;