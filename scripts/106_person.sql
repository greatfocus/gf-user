CREATE TABLE IF NOT EXISTS person (
	id BIGSERIAL PRIMARY KEY,
	userId INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	countryId INTEGER REFERENCES country(id),
	firstName VARCHAR(30) NOT NULL,
	middleName VARCHAR(30) NOT NULL,
	lastName VARCHAR(30) NOT NULL,
	mobileNumber VARCHAR(20) NOT NULL,
	idNumber VARCHAR(50) NOT NULL,
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	createdBy INTEGER NOT NULL,
	updatedOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updatedBy INTEGER NOT NULL,
	UNIQUE(userId),
	UNIQUE(mobilenumber),
	UNIQUE(idNumber),
	UNIQUE(mobilenumber, idNumber)	
);

DO $$ 
DECLARE
	country INTEGER := (select id from country where code='KE');
	userId INTEGER := (select id from users where email='muthurimixphone@gmail.com');

BEGIN
	INSERT INTO person (userId, countryId, firstName, middleName, lastName, mobileNumber, idNumber, createdBy, updatedBy)
	VALUES
		(userId, country, 'John', 'Peter', 'Mucunga', '0780904371', '27496388', userId, userId) 
	ON CONFLICT
	DO NOTHING;
END $$;