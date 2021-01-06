DO $$
DECLARE
	country INTEGER := (select id from country where code='KE');
	userId INTEGER := (select id from users where email='muthurimixphone@gmail.com');
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM person) THEN
        INSERT INTO person (userId, countryId, firstName, middleName, lastName, mobileNumber, idNumber, createdBy, updatedBy)
        VALUES
            (userId, country, 'John', 'Peter', 'Mucunga', '0780904371', '27496388', userId, userId) 
        ON CONFLICT
        DO NOTHING;
    END IF;
END $$