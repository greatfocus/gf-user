CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	type VARCHAR(10) NOT NULL,
	firstName VARCHAR(20) NOT NULL,
	middleName VARCHAR(20) NOT NULL,
	lastName VARCHAR(20) NOT NULL,
	mobileNumber VARCHAR(14) NOT NULL,
	email VARCHAR(50) NOT NULL,
	password VARCHAR(100) NOT NULL,
	failedAttempts FLOAT default (0),
	lastAttempt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	lastChange TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expiredDate TIMESTAMP NOT NULL,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	createdBy INTEGER NULL,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedBy INTEGER NULL,
	status VARCHAR(20) NOT NULL,
	deleted BOOLEAN NOT NULL default(false),
	enabled BOOLEAN NOT NULL default(false),
	UNIQUE(mobilenumber),
	UNIQUE(email),
	UNIQUE(email, mobilenumber)
);

INSERT INTO users (type, firstName, middleName, lastName, mobileNumber, email, password, expiredDate, status, enabled)
VALUES
	('password', 'John', 'Peter', 'Mucunga', '0780904371', 'muthurimixphone@gmail.com', '$2a$04$cZT44rp5yKqGox31VZpxieNq/XfoSJAMqoodhI/gUBvNvcn.kUUWe', '2025-06-01 08:22:17.460493', 'USER.APPROVED', true) 
ON CONFLICT (email) 
DO NOTHING;