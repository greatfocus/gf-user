CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	type VARCHAR(10) REFERENCES userType(id),
	countryId INTEGER REFERENCES country(id),
	firstName VARCHAR(30) NOT NULL,
	middleName VARCHAR(30) NOT NULL,
	lastName VARCHAR(30) NOT NULL,
	mobileNumber VARCHAR(20) NOT NULL,
	email VARCHAR(50) NOT NULL,
	idNumber VARCHAR(50) NOT NULL,
	password VARCHAR(100) NOT NULL,
	failedAttempts INTEGER default (0),
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
	UNIQUE(idNumber),
	UNIQUE(mobilenumber),
	UNIQUE(email)	
);

DO $$ 
DECLARE
	country INTEGER := (select id from country where code='KE');

BEGIN
	INSERT INTO users (type, countryId, firstName, middleName, lastName, mobileNumber, email, idNumber, password, expiredDate, status, enabled)
	VALUES
		('admin', country, 'John', 'Peter', 'Mucunga', '0780904371', 'muthurimixphone@gmail.com', '27496388', '$2a$04$cZT44rp5yKqGox31VZpxieNq/XfoSJAMqoodhI/gUBvNvcn.kUUWe', '2025-06-01 08:22:17.460493', 'USER.APPROVED', true) 
	ON CONFLICT (idNumber, mobilenumber, email) 
	DO NOTHING;
END $$;